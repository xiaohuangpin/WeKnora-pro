package tracing

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	AppName = "WeKnoraApp"
)

type Tracer struct {
	Cleanup func(context.Context) error
}

var tracer trace.Tracer

// InitTracer initializes OpenTelemetry tracer
func InitTracer() (*Tracer, error) {
	// Create resource description
	labels := []attribute.KeyValue{
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String(AppName),
	}
	res := resource.NewWithAttributes(semconv.SchemaURL, labels...)
	var err error

	// First try to create OTLP exporter (can connect to Jaeger, Zipkin, etc.)
	var traceExporter sdktrace.SpanExporter
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		// Use gRPC exporter
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
		traceExporter, err = otlptrace.New(context.Background(), client)
		if err != nil {
			return nil, err
		}
	} else {
		// If no OTLP endpoint is set, default to standard output
		traceExporter, err = stdouttrace.New()
		if err != nil {
			return nil, err
		}
	}

	// Create batch SpanProcessor
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	sampler := sdktrace.AlwaysSample()

	// Create and register TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create Tracer for project use
	tracer = tp.Tracer(AppName)

	// Return cleanup function
	return &Tracer{
		Cleanup: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
				return err
			}
			return nil
		},
	}, nil
}

// GetTracer gets global Tracer
func GetTracer() trace.Tracer {
	return tracer
}

// Create context with span
func ContextWithSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, name, opts...)
}
