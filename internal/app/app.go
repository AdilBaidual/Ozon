// nolint: all // not compatible with Uber.Fx
package app

import (
	"Service/config"
	"Service/internal/auth"
	authRepo "Service/internal/auth/repository"
	authUC "Service/internal/auth/usecase"
	"Service/internal/core/delivery"
	graph2 "Service/internal/graph"
	"Service/internal/middleware"
	"Service/pkg/httpserver"
	"Service/pkg/jaeger"
	"Service/pkg/paseto"
	postgresStorage "Service/pkg/storage/postgres"
	valkeyStorage "Service/pkg/storage/valkey"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"strconv"
)

func NewApp() fx.Option {
	return fx.Options(
		ConfigModule(),
		LoggerModule(),
		JaegerModule(),
		StorageModule(),
		PostgresModule(),
		RepositoryModule(),
		UseCaseModule(),
		GraphqlModule(),
		CheckInitializedModules(),
	)
}

func ConfigModule() fx.Option {
	return fx.Module("config",
		fx.Provide(
			config.New,
		),
	)
}

func LoggerModule() fx.Option {
	return fx.Module("logger",
		fx.Provide(
			func() *zap.Logger {
				encoderCfg := zap.NewProductionEncoderConfig()
				encoderCfg.TimeKey = "timestamp"
				encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

				cfg := zap.Config{
					Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
					Development:       false,
					DisableCaller:     false,
					DisableStacktrace: false,
					Sampling:          nil,
					Encoding:          "json",
					EncoderConfig:     encoderCfg,
					OutputPaths: []string{
						"stderr",
					},
					ErrorOutputPaths: []string{
						"stderr",
					},
				}

				return zap.Must(cfg.Build())
			},
		),
	)
}

func JaegerModule() fx.Option {
	return fx.Module("jaeger",
		fx.Provide(
			func(cfg *config.Config) config.Jaeger {
				return cfg.Jaeger
			},
			jaeger.InitJaeger,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, tracer *sdktrace.TracerProvider, cfg config.Jaeger, logger *zap.Logger, shutdowner fx.Shutdowner) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						err := tracer.Shutdown(ctx)
						if err != nil {
							logger.Error("Error shutting down tracer provider", zap.Error(err))
							return err
						}
						return nil
					},
				})
			},
		),
	)
}

func StorageModule() fx.Option {
	return fx.Module("storage",
		fx.Provide(
			func(cfg *config.Config) *valkeyStorage.Storage {
				return valkeyStorage.New(valkeyStorage.Config{
					ClientName: config.ServiceName,
					Host:       cfg.Valkey.Host,
					Port:       cfg.Valkey.Port,
					Password:   cfg.Valkey.Password,
					Database:   2,
				})
			},
		),
	)
}

func PostgresModule() fx.Option {
	return fx.Module("postgres",
		fx.Provide(
			func(cfg *config.Config) (*sqlx.DB, error) {
				return postgresStorage.New(postgresStorage.Config{
					Host:            cfg.Postgres.Host,
					Port:            cfg.Postgres.Port,
					User:            cfg.Postgres.User,
					Password:        cfg.Postgres.Password,
					DBName:          cfg.Postgres.DBName,
					SSLMode:         cfg.Postgres.SSLMode,
					ApplicationName: config.ServiceName,
					PgDriver:        cfg.Postgres.PgDriver,
				})
			},
		),
	)
}

// TODO: ADD SWITCH LOGIC
func RepositoryModule() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			fx.Annotate(
				authRepo.NewRepository,
				fx.As(new(auth.Repo)),
			),
		),
	)
}

func UseCaseModule() fx.Option {
	return fx.Module("usecase",
		fx.Provide(
			func(cfg *config.Config) (*paseto.Paseto, error) {
				return paseto.NewPaseto(cfg.Paseto.PasetoSecret)
			},
			fx.Annotate(
				authUC.NewUseCase,
				fx.As(new(auth.UseCase)),
			),
		),
	)
}

func GraphqlModule() fx.Option {
	return fx.Module("graphql",
		fx.Provide(
			func(cfg *config.Config) httpserver.Config {
				return cfg.Server
			},
			graph2.NewResolver,
			middleware.NewMiddleware,
			delivery.NewHandler,
			fx.Annotate(
				func(h *delivery.Handler) *gin.Engine {
					return h.InitRoutes()
				},
				fx.As(new(http.Handler)),
			),
			httpserver.NewServer,
		),
		fx.Invoke(
			//func(router *gin.Engine, logger *zap.Logger, cfg config.Server, resolver *graph2.Resolver) {
			//	srv := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: resolver}))
			//	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
			//	http.Handle("/query", srv)
			//
			//	logger.Info(fmt.Sprintf("starting GraphQL server {%s}", net.JoinHostPort(cfg.Host, cfg.Port)))
			//	logger.Fatal("starting GraphQL server", zap.Error(http.ListenAndServe(net.JoinHostPort(cfg.Host, cfg.Port), nil)))
			//},
			func(lc fx.Lifecycle, srv *httpserver.Server, cfg httpserver.Config, logger *zap.Logger, shutdowner fx.Shutdowner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							logger.Info(fmt.Sprintf("starting HTTP server {%s}", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))))
							if err := srv.Start(); err != nil {
								logger.Error("error starting HTTP server",
									zap.Error(err),
									zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						if err := srv.Stop(ctx); err != nil {
							logger.Error("error stopping HTTP server",
								zap.Error(err),
								zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))),
							)
						}
						return nil
					},
				})
			}),
	)
}

//func HTTPModule() fx.Option {
//	return fx.Module("http server",
//		fx.Provide(
//			func(cfg *config.Config) httpserver.Config {
//				return cfg.HTTPServer
//			},
//			middleware.NewMiddleware,
//			httphandler.NewHandler,
//			fx.Annotate(
//				func(h *httphandler.Handler) *gin.Engine {
//					return h.InitRoutes()
//				},
//				fx.As(new(http.Handler)),
//			),
//			httpserver.NewServer,
//		),
//		fx.Invoke(
//			func(lc fx.Lifecycle, srv *httpserver.Server, cfg httpserver.Config, logger *zap.Logger, shutdowner fx.Shutdowner) {
//				lc.Append(fx.Hook{
//					OnStart: func(ctx context.Context) error {
//						go func() {
//							logger.Info(fmt.Sprintf("starting HTTP server {%s}", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))))
//							if err := srv.Start(); err != nil {
//								logger.Error("error starting HTTP server",
//									zap.Error(err),
//									zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))),
//								)
//							}
//						}()
//						return nil
//					},
//					OnStop: func(ctx context.Context) error {
//						if err := srv.Stop(ctx); err != nil {
//							logger.Error("error stopping HTTP server",
//								zap.Error(err),
//								zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))),
//							)
//						}
//						return nil
//					},
//				})
//			},
//		),
//	)
//}

func CheckInitializedModules() fx.Option {
	return fx.Module("check modules",
		fx.Invoke(
			func(cfg *config.Config) {},
			func(lg *zap.Logger) {},
			func(storage *valkeyStorage.Storage) {},
			func(db *sqlx.DB) {},
			//func(authR core.Repo) {},
			//func(authUC core.UseCase) {},
		),
	)
}
