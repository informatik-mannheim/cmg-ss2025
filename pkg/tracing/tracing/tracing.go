package tracing

import (
	"context"
	"net/http"
	"strconv"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type customErrorHandler struct{}

func (h customErrorHandler) Handle(err error) {
	logging.Error(err.Error())
}

var _ otel.ErrorHandler = customErrorHandler{} // type safety check

var tracer trace.Tracer = otel.Tracer("default-tracer")

// Init sets up tracing with OTLP (Jaeger via OTLP collector endpoint)
func Init(serviceName, otlpEndpoint string) (func(context.Context) error, error) {
	logging.Init("tracing")

	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(), // I would say this is fine because it should be in a local network, so no TLS is needed
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(customErrorHandler{})
	tracer = tp.Tracer(serviceName)

	return tp.Shutdown, nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		type statusRecorder struct {
			http.ResponseWriter
			status int
		}
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rec, r.WithContext(ctx))

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.code", strconv.Itoa(rec.status)),
		)
	})
}
