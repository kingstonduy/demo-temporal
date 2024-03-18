package bootstrap

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean-core/config"
	"github.com/lengocson131002/go-clean-core/trace"
	"github.com/lengocson131002/go-clean-core/trace/otel"
)

func GetTracer(srvCfg *ServerConfig, cfg config.Configure) trace.Tracer {
	endpoint := cfg.GetString("TRACE_ENDPOINT")
	tracer, err := otel.NewOpenTelemetryTracer(context.Background(), srvCfg.Name, endpoint)
	if err != nil {
		panic(fmt.Errorf("failed to create tracer object: %w", err))
	}
	return tracer
}
