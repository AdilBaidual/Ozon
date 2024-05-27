package delivery

import (
	"Service/config"
	"Service/internal/core"
	"Service/internal/core/delivery/graph"
	"Service/internal/middleware"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Handler struct {
	engine   *gin.Engine
	mw       *middleware.Middleware
	resolver *graph.Resolver
	coreUC   core.UseCase
}

func NewHandler(mw *middleware.Middleware, resolver *graph.Resolver, coreUC core.UseCase) *Handler {
	return &Handler{
		mw:       mw,
		resolver: resolver,
		coreUC:   coreUC,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	h.engine = router

	router.Use(cors.Default())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(h.mw.LoggingMiddleware())
	router.Use(h.mw.ValidatePasetoToken())
	router.Use(graph.DataLoader(h.coreUC))

	router.POST("/query", h.graphqlHandler())
	router.GET("/query", h.graphqlHandler())
	router.GET("/", h.playgroundHandler())

	return h.engine
}

func (h *Handler) graphqlHandler() gin.HandlerFunc {
	graphqlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: h.resolver}))
	// graphqlHandler.Use(extension.FixedComplexityLimit(5))
	graphqlHandler.AddTransport(&transport.Websocket{})
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
