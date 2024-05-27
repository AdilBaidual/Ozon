// nolint: all // not compatible with Uber.Fx
package app

import (
	"Service/config"
	"Service/internal/auth"
	authRepo "Service/internal/auth/repository"
	authUC "Service/internal/auth/usecase"
	"Service/internal/core"
	"Service/internal/core/delivery"
	graph2 "Service/internal/core/delivery/graph"
	coreRepoPostgres "Service/internal/core/repository/postgres"
	coreRepoRedis "Service/internal/core/repository/redis"
	coreUC "Service/internal/core/usecase"
	"Service/internal/middleware"
	"Service/pkg/httpserver"
	"Service/pkg/jaeger"
	"Service/pkg/paseto"
	postgresStorage "Service/pkg/storage/postgres"
	redisStorage "Service/pkg/storage/redis"
	valkeyStorage "Service/pkg/storage/valkey"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
		RedisModule(),
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

func RedisModule() fx.Option {
	return fx.Module("redis",
		fx.Provide(
			func(cfg *config.Config) (*redis.Client, error) {
				return redisStorage.NewRedisClient(&redis.Options{
					Addr:     net.JoinHostPort(cfg.Redis.Host, strconv.Itoa(cfg.Redis.Port)),
					Password: "",
					DB:       0,
				})
			},
		),
	)
}

func RepositoryModule() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			fx.Annotate(
				authRepo.NewRepository,
				fx.As(new(auth.Repo)),
			),
			func(db *sqlx.DB, client *redis.Client) core.Repo {
				if config.InMemory {
					return coreRepoRedis.NewRepository(client)
				}
				return coreRepoPostgres.NewRepository(db)
			},
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
			coreUC.NewUseCase,
			func(uc *coreUC.UC) core.UseCase {
				return uc
			},
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

func CheckInitializedModules() fx.Option {
	return fx.Module("check modules",
		fx.Invoke(
			func(cfg *config.Config) {},
			func(lg *zap.Logger) {},
			func(storage *valkeyStorage.Storage) {},
			func(db *sqlx.DB) {},
			func(redisClient *redis.Client) {},
			func(authR auth.Repo, coreR core.Repo) {},
			func(authUC auth.UseCase, coreUC core.UseCase) {},
		),
	)
}
