package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HoanNghi16/rentix-services/property-service/internal/database"
	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/property-service/internal/messaging"
	"github.com/HoanNghi16/rentix-services/property-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/property-service/internal/service"
	"github.com/HoanNghi16/rentix-services/property-service/internal/transport/grpcserver"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	db := database.ConnectDB()
	port := os.Getenv("PROPERTY_GRPC_PORT")
	if port == "" {
		log.Fatal("PROPERTY_GRPC_PORT is required")
	}
	address := ":" + port

	transportCredentials, err := grpcapi.LoadServerTLSCredentials()
	if err != nil {
		log.Fatal("cannot configure property gRPC mutual TLS: ", err)
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot listen for property gRPC: ", err)
	}

	server := grpc.NewServer(grpc.Creds(transportCredentials))
	propertyService := service.NewService(repository.NewRepo(db))
	grpcserver.RegisterPropertyServiceServer(server, grpcserver.NewPropertyServer(propertyService))
	publisher, err := messaging.NewKafkaPublisherFromEnv()
	if err != nil {
		log.Fatal("cannot configure Kafka publisher: ", err)
	}
	defer publisher.Close()

	workerContext, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		runOutboxPublisher(workerContext, propertyService, publisher)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		cancelWorker()
		server.GracefulStop()
		<-workerDone
	}()

	log.Printf("property gRPC server listening on %s", address)
	if err := server.Serve(listener); err != nil {
		log.Fatal(fmt.Errorf("property gRPC server stopped: %w", err))
	}
}

func runOutboxPublisher(ctx context.Context, propertyService *service.PropService, publisher *messaging.KafkaPublisher) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := propertyService.PublishOutbox(ctx, publisher); err != nil && ctx.Err() == nil {
				log.Printf("outbox publisher iteration failed: %v", err)
			}
		}
	}
}
