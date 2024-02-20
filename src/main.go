package main

import (
	"go-opentelemetry/controllers"
	"go-opentelemetry/initializers"
	"go-opentelemetry/metrics"
	"log"
	"os"
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	otlpEndpoint = os.Getenv("OPTLP_ENDPOINT")
	serviceName = os.Getenv("SERVICE_NAME")
)

func initTracer() func(context.Context) error {
	
	insecureOpt := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			insecureOpt,
			otlptracegrpc.WithEndpoint(otlpEndpoint),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// Run function before main
func init() {
	initializers.ConnectToDB()
}

func main(){
	cleanup := initTracer()
	defer cleanup(context.Background())

	provider := metrics.InitMeter()
	defer provider.Shutdown(context.Background())

	meter := provider.Meter("go-opentelemetry")
	metrics.GenerateMetrics(meter)

	router := gin.Default()
	router.Use(otelgin.Middleware(serviceName))

	router.POST("/api/user/", controllers.AddUser)
	router.GET("/api/users/", controllers.GetUsers)
	router.GET("/api/user/:id", controllers.GetUser)
	router.DELETE("/api/user/:id", controllers.DeleteUser)
	router.PUT("/api/user/:id", controllers.UpdateUser)
	router.PUT("/status", controllers.Status)
	router.Run()
}