package servicesetup

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func SetupTracing(logger *zap.Logger) {
	jg, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Process:       jaeger.Process{ServiceName: "social-todo-app"},
	})
	if err != nil {
		logger.Fatal("error happen when setting up Jeager tracer", zap.Error(err))
	}

	trace.RegisterExporter(jg)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})
}
