package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	healthchecks "github.com/lengocson131002/go-clean-core/health"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/transport/http"
	"github.com/ocb/mcs-money-transfer/bootstrap"
	"github.com/ocb/mcs-money-transfer/presentation/http/controller"
)

type HttpServer struct {
	Port                    int
	Name                    string
	App                     *fiber.App
	Logger                  logger.Logger
	HealhChecker            healthchecks.HealthChecker
	AccountController       *controller.AccountController
	MoneyTransferController *controller.MoneyTransferController
}

// @title  GOLANG TEMPORAL DEMO
// @version 1.0
// @description GOLANG TEMPORAL DEMO
// @termsOfService http://swagger.io/terms/
// @contact.name LNS
// @contact.email leson131002@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func NewHttpServer(
	cfg *bootstrap.ServerConfig,
	logger logger.Logger,
	healhChecker healthchecks.HealthChecker,
	accountController *controller.AccountController,
	moneyTransferController *controller.MoneyTransferController,
) *HttpServer {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: CustomErrorHandler,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	return &HttpServer{
		Port:                    cfg.HttpPort,
		Name:                    cfg.Name,
		Logger:                  logger,
		App:                     fiberApp,
		HealhChecker:            healhChecker,
		AccountController:       accountController,
		MoneyTransferController: moneyTransferController,
	}
}

func (s *HttpServer) GetStartOptions() []HttpServerStartOption {
	return []HttpServerStartOption{
		WithSwagger(),
		WithLoggings(),
		WithHealthCheck(),
		WithTracing(),
		WithMetrics(),
		WithAccountV1Routes(),
		WithMoneyTransferV1Routes(),
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	// configs options
	opts := s.GetStartOptions()
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			if err := s.App.Shutdown(); err != nil {
				s.Logger.Errorf(ctx, "Failed to shutdown http server: %v", err)
			}
			s.Logger.Info(ctx, "Stop HTTP Server")
		}(ctx)

		<-ctx.Done()
	}()

	hPort := s.Port
	s.Logger.Infof(ctx, "Start HTTP server at port: %v", hPort)
	if err := s.App.Listen(fmt.Sprintf(":%v", hPort)); err != nil {
		s.Logger.Errorf(ctx, "Failed to start http server: %v ", err)
		return err
	}

	return nil
}

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	fRes := http.FailureResponse(err)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Status(fRes.Result.Status).JSON(fRes)
}
