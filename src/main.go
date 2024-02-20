package main

import (
	"github.com/eshiettjoseph/go-opentelemetry/src/controllers"
	"github.com/eshiettjoseph/go-opentelemetry/src/initializers"
	"log"
	"os"
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	otlpEndpoint = os.Getenv("OPTLP_ENDPOINT")
	serviceName = os.Getenv("SERVICE_NAME")
)

func initTracer() func(context.Context) error {
	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()

	// Update default OTLP reciver endpoint
	endpointOpt := otlptracehttp.WithEndpoint(otlpEndpoint)

	exporter, err := otlptracehttp.New(
		context.Background(), 
		insecureOpt, endpointOpt,
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