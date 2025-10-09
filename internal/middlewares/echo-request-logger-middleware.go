package middlewares

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoRequestLoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			errorMessage := ""

			if v.Error != nil {
				errorMessage = v.Error.Error()
			}

			log := logger.With(
				slog.String("request_id", v.RequestID),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("host", v.Host),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("user_agent", v.UserAgent),
				slog.Int("status_code", v.Status),
				slog.String("error", errorMessage),
				slog.Float64("latency_ms", float64(v.Latency.Microseconds())/1000),
			)

			if c.Response().Status >= 500 {
				log.Error("request error")
				return nil
			}

			log.Info("request success")
			return nil
		},
	})
}
