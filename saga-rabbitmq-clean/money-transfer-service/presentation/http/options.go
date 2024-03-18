package http

import (
	"os"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

type HttpServerStartOption func(*HttpServer) error

func WithLoggings() HttpServerStartOption {
	return func(s *HttpServer) error {
		s.App.Use(fiberLog.New(fiberLog.Config{
			Next:         nil,
			Done:         nil,
			Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stdout,
		}))
		return nil
	}
}

func WithSwagger() HttpServerStartOption {
	return func(s *HttpServer) error {
		s.App.Get("/swagger/*", swagger.HandlerDefault)
		return nil
	}
}

func WithHealthCheck() HttpServerStartOption {
	return func(s *HttpServer) error {
		s.App.Get("/liveliness", func(c *fiber.Ctx) error {
			result := s.HealhChecker.LivenessCheck()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})

		s.App.Get("/readiness", func(c *fiber.Ctx) error {
			result := s.HealhChecker.RedinessCheck()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})
		return nil
	}
}

func WithTracing() HttpServerStartOption {
	return func(s *HttpServer) error {
		s.App.Use(otelfiber.Middleware())
		return nil
	}
}

func WithMetrics() HttpServerStartOption {
	return func(s *HttpServer) error {
		prometheus := fiberprometheus.New(s.Name)
		prometheus.RegisterAt(s.App, "/metrics")
		s.App.Use(prometheus.Middleware)
		return nil
	}
}
