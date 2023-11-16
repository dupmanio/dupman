package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const connectionTimeout = 5 * time.Second

type OTel struct {
	serviceName    string
	serviceVersion string
	collectorURL   string

	grpcConnection *grpc.ClientConn

	resource    *resource.Resource
	propagation propagation.TextMapPropagator

	tracerProvider  *trace.TracerProvider
	metricsProvider *metric.MeterProvider
}

func NewOTel(ctx context.Context, serviceName, serviceVersion, collectorURL string) (*OTel, error) {
	ot := &OTel{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		collectorURL:   collectorURL,
	}

	if err := ot.setupGRPCConnection(ctx); err != nil {
		return nil, err
	}

	if err := ot.setupResource(); err != nil {
		return nil, err
	}

	ot.setupPropagator()

	if err := ot.setupTraceProvider(ctx); err != nil {
		return nil, err
	}

	if err := ot.setupMeterProvider(ctx); err != nil {
		return nil, err
	}

	return ot, nil
}

func (ot *OTel) setupGRPCConnection(ctx context.Context) error {
	var err error

	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	ot.grpcConnection, err = grpc.DialContext(
		ctx,
		ot.collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return nil
}

func (ot *OTel) setupResource() error {
	var err error

	// @todo: add environment.
	ot.resource, err = resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(ot.serviceName),
			semconv.ServiceVersion(ot.serviceVersion),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to merge resources: %w", err)
	}

	return nil
}

func (ot *OTel) setupPropagator() {
	ot.propagation = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(ot.propagation)
}

func (ot *OTel) setupTraceProvider(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(ot.grpcConnection))
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	ot.tracerProvider = trace.NewTracerProvider(
		trace.WithResource(ot.resource),
		trace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(ot.tracerProvider)

	return nil
}

func (ot *OTel) setupMeterProvider(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(ot.grpcConnection))
	if err != nil {
		return fmt.Errorf("failed to create metric exporter: %w", err)
	}

	ot.metricsProvider = metric.NewMeterProvider(
		metric.WithResource(ot.resource),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	otel.SetMeterProvider(ot.metricsProvider)

	return nil
}

func (ot *OTel) Shutdown(ctx context.Context) error {
	if err := ot.metricsProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to shutdown Metrics Provider: %w", err)
	}

	if err := ot.tracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to shutdown Traces Provider: %w", err)
	}

	if err := ot.grpcConnection.Close(); err != nil {
		return fmt.Errorf("unable to close GRPC Connection: %w", err)
	}

	return nil
}
