package middleware

import (
	"Service/internal/auth"
	"Service/pkg/paseto"
	valkeyStorage "Service/pkg/storage/valkey"
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type Middleware struct {
	logger       *zap.Logger
	authUC       auth.UseCase
	paseto       *paseto.Paseto
	tokenStorage *valkeyStorage.Storage
}

func NewMiddleware(logger *zap.Logger, authUC auth.UseCase, paseto *paseto.Paseto, tokenStorage *valkeyStorage.Storage) *Middleware {
	return &Middleware{
		logger:       logger,
		authUC:       authUC,
		paseto:       paseto,
		tokenStorage: tokenStorage,
	}
}

func (mw *Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()

		spanContext := trace.SpanContextFromContext(ctx.Request.Context())
		requestLogger := mw.logger.With(zap.String("request_id", spanContext.TraceID().String()))
		ctx.Set("logger", requestLogger)

		ctx.Next()

		duration := time.Since(start)

		logInfos := []zap.Field{zap.String("method", ctx.Request.URL.Path), zap.String("processing time", duration.String())}
		if ctx.Errors != nil {
			logInfos = append(logInfos, zap.String("errors", ctx.Errors.String()))
		}

		requestLogger.Info("Request info", logInfos...)
	}
}

func (m *Middleware) ValidatePasetoToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger, ok := ctx.Value("logger").(*zap.Logger)
		if !ok {
			logger = m.logger
		}
		access := ctx.Request.Header.Get("Access")
		if access == "" {
			logger.Info("empty access")
			ctx.Next()
			return
		}

		bl, err := m.tokenStorage.Get(access + "_blacklist")
		if bl != nil || err != nil {
			logger.Info("access token was deleted")
			ctx.Next()
			return
		}

		uuid, err := m.paseto.ValidateToken(access)
		if err != nil {
			logger.Info("validate access token", zap.Error(err))
			ctx.Next()
			return
		}

		c := context.WithValue(ctx.Request.Context(), "uuid", uuid)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}
