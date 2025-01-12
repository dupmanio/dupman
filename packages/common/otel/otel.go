package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
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

	grpcConnection *grpc.ClientConn

	resource    *resource.Resource
	propagation propagation.TextMapPropagator

	tracerProvider  *sdkTrace.TracerProvider
	metricsProvider *sdkMetric.MeterProvider
	loggerProvider  *sdkLog.LoggerProvider

	Tracer trace.Tracer
	Meter  metric.Meter
	Logger log.Logger
}

func NewOTel(
	ctx context.Context,
	env string,
	serviceName string,
	serviceVersion string,
	collectorURL string,
) (*OTel, error) {
	ot := &OTel{
		env:            env,
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		collectorURL:   collectorURL,
	}

	if err := ot.setupGRPCConnection(); err != nil {
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

	if err := ot.setupLogProvider(ctx); err != nil {
		return nil, err
	}

	ot.Tracer = ot.tracerProvider.Tracer(fmt.Sprintf("dupman.io/service/%s", serviceName))
	ot.Meter = ot.metricsProvider.Meter(fmt.Sprintf("dupman.io/service/%s", serviceName))
	ot.Logger = ot.loggerProvider.Logger(fmt.Sprintf("dupman.io/service/%s", serviceName))

	return ot, nil
}

func (ot *OTel) setupGRPCConnection() error {
	var err error

	ot.grpcConnection, err = grpc.NewClient(
		ot.collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func (ot *OTel) setupLogProvider(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(ot.grpcConnection))
	if err != nil {
		return fmt.Errorf("failed to create log exporter: %w", err)
	}

	ot.loggerProvider = sdkLog.NewLoggerProvider(
		sdkLog.WithResource(ot.resource),
		sdkLog.WithProcessor(sdkLog.NewBatchProcessor(logExporter)),
	)

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

func (ot *OTel) GetZapCore() *otelzap.Core {
	return otelzap.NewCore(
		fmt.Sprintf("dupman.io/service/%s", ot.serviceName),
		otelzap.WithLoggerProvider(ot.loggerProvider),
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

// @todo: refactor/re-plan logging.

func (ot *OTel) LogInfoEvent(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	ot.InfoEvent(ctx, message, attributes...)
}

func (ot *OTel) InfoEvent(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent(message, trace.WithAttributes(attributes...))
	}
}

func (ot *OTel) LogErrorEvent(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	ot.ErrorEvent(ctx, message, err, attributes...)
}

func (ot *OTel) ErrorEvent(ctx context.Context, message string, err error, attributes ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err, trace.WithAttributes(attributes...))
		span.SetStatus(codes.Error, message)
	}
}

func (ot *OTel) GetSpanForFunctionCall(
	ctx context.Context,
	skipCaller int,
	attributes ...attribute.KeyValue,
) (context.Context, trace.Span) {
	function, functionAttributes := GetFunctionCallAttributes(skipCaller)

	//nolint: spancheck
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

	if err := ot.loggerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to shutdown Log Provider: %w", err)
	}

	if err := ot.grpcConnection.Close(); err != nil {
		return fmt.Errorf("unable to close GRPC Connection: %w", err)
	}

	return nil
}
