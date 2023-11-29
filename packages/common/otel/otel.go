package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

const connectionTimeout = 5 * time.Second

type OTel struct {
	env            string
	serviceName    string
	serviceVersion string
	collectorURL   string
	logger         *zap.Logger

	grpcConnection *grpc.ClientConn

	resource    *resource.Resource
	propagation propagation.TextMapPropagator

	tracerProvider  *sdkTrace.TracerProvider
	metricsProvider *sdkMetric.MeterProvider

	Tracer trace.Tracer
	Meter  metric.Meter
}

func NewOTel(env, serviceName, serviceVersion, collectorURL string, logger *zap.Logger) (*OTel, error) {
	ctx := context.Background()
	ot := &OTel{
		env:            env,
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		collectorURL:   collectorURL,
		logger:         logger,
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

	ot.Tracer = ot.tracerProvider.Tracer(fmt.Sprintf("dupman.io/service/%s", serviceName))
	ot.Meter = ot.metricsProvider.Meter(fmt.Sprintf("dupman.io/service/%s", serviceName))

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

	ot.resource, err = resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(ot.serviceName),
			semconv.ServiceVersion(ot.serviceVersion),
			semconv.DeploymentEnvironment(ot.env),
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

	ot.tracerProvider = sdkTrace.NewTracerProvider(
		sdkTrace.WithResource(ot.resource),
		sdkTrace.WithBatcher(traceExporter),
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

	ot.metricsProvider = sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(ot.resource),
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(metricExporter)),
	)

	otel.SetMeterProvider(ot.metricsProvider)

	return nil
}

func (ot *OTel) GetOTelGinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(
		ot.serviceName,
		otelgin.WithTracerProvider(ot.tracerProvider),
	)
}

func (ot *OTel) GetGormPlugin(dbName string) gorm.Plugin {
	return tracing.NewPlugin(
		tracing.WithDBName(dbName),
		tracing.WithoutQueryVariables(),
		tracing.WithTracerProvider(ot.tracerProvider),
	)
}

func (ot *OTel) InstrumentRedis(redisClient redis.UniversalClient) error {
	if err := redisotel.InstrumentTracing(redisClient, redisotel.WithTracerProvider(ot.tracerProvider)); err != nil {
		return fmt.Errorf("unable to instrument redis tracing: %w", err)
	}

	if err := redisotel.InstrumentMetrics(redisClient, redisotel.WithMeterProvider(ot.metricsProvider)); err != nil {
		return fmt.Errorf("unable to instrument redis metrics: %w", err)
	}

	return nil
}

// LogInfoEvent combines InfoEvent and InfoLog actions.
func (ot *OTel) LogInfoEvent(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	ot.InfoEvent(ctx, message, attributes...)
	ot.InfoLog(ctx, message, attributes...)
}

func (ot *OTel) InfoEvent(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent(message, trace.WithAttributes(attributes...))
	}
}

func (ot *OTel) InfoLog(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	ot.logInfoOrError(ctx, message, nil, attributes...)
}

func (ot *OTel) logInfoOrError(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	fields := ot.convertAttributesToZapFields(
		append(
			attributes,
			ot.getCommonLogAttributes(ctx)...,
		),
	)

	if err == nil {
		ot.logger.Info(message, fields...)

		return
	}

	ot.logger.Error(message, append(fields, zap.Error(err))...)
}

// LogErrorEvent combines ErrorEvent and ErrorLog actions.
func (ot *OTel) LogErrorEvent(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	ot.ErrorEvent(ctx, message, err, attributes...)
	ot.ErrorLog(ctx, message, err, attributes...)
}

func (ot *OTel) ErrorEvent(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err, trace.WithAttributes(attributes...))
		span.SetStatus(codes.Error, message)
	}
}

func (ot *OTel) ErrorLog(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	ot.logInfoOrError(ctx, message, err, attributes...)
}

func (ot *OTel) getCommonLogAttributes(ctx context.Context) []attribute.KeyValue {
	spanContext := trace.SpanFromContext(ctx).SpanContext()

	return []attribute.KeyValue{
		TraceID(spanContext.TraceID().String()),
		SpanID(spanContext.SpanID().String()),
	}
}

func (ot *OTel) convertAttributesToZapFields(attributes []attribute.KeyValue) []zapcore.Field {
	fields := make([]zapcore.Field, len(attributes))
	for i, attr := range attributes {
		fields[i] = zap.Any(string(attr.Key), attr.Value.AsInterface())
	}

	return fields
}

func (ot *OTel) GetSpanForFunctionCall(
	ctx context.Context,
	skipCaller int,
	attributes ...attribute.KeyValue,
) (context.Context, trace.Span) {
	function, functionAttributes := GetFunctionCallAttributes(skipCaller)

	return ot.Tracer.Start(
		ctx,
		function,
		trace.WithAttributes(functionAttributes...),
		trace.WithAttributes(attributes...),
	)
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
