package delivery

import (
	"Service/config"
	"Service/internal/graph"
	"Service/internal/middleware"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	engine   *gin.Engine
	mw       *middleware.Middleware
	resolver *graph.Resolver
}

func NewHandler(mw *middleware.Middleware, resolver *graph.Resolver) *Handler {
	return &Handler{
		mw:       mw,
		resolver: resolver,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	h.engine = router

	router.Use(cors.Default())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(h.mw.LoggingMiddleware())
	router.Use(h.mw.ValidatePasetoToken())

	router.POST("/query", h.graphqlHandler())
	router.GET("/", h.playgroundHandler())
	router.GET("/ping", func(ctx *gin.Context) {
		tracer := otel.Tracer(config.ServiceName)
		_, span := tracer.Start(ctx.Request.Context(), "pong")
		defer span.End()

		logger, ok := ctx.Value("logger").(*zap.Logger)
		if !ok {
			fmt.Println("logger not found")
		} else {
			logger.Info("logger found!")
		}

		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	return h.engine
}

func (h *Handler) graphqlHandler() gin.HandlerFunc {
	graphqlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: h.resolver}))

	return func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *Handler) playgroundHandler() gin.HandlerFunc {
	graphqlHandler := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	}
}
