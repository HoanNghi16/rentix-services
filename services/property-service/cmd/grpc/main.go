package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/HoanNghi16/rentix-services/property-service/internal/database"
	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
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
	grpcserver.RegisterPropertyServiceServer(server, grpcserver.NewPropertyServer(service.NewService(repository.NewRepo(db))))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		server.GracefulStop()
	}()

	log.Printf("property gRPC server listening on %s", address)
	if err := server.Serve(listener); err != nil {
		log.Fatal(fmt.Errorf("property gRPC server stopped: %w", err))
	}
}
