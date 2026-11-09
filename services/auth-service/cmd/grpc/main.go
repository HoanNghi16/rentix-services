package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/database"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("cannot connect database: ", err)
	}

	port := os.Getenv("AUTH_GRPC_PORT")
	if port == "" {
		log.Fatal("AUTH_GRPC_PORT is required")
	}
	address := ":" + port
	transportCredentials, err := grpcapi.LoadServerTLSCredentials()
	if err != nil {
		log.Fatal("cannot configure gRPC mutual TLS: ", err)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot listen for gRPC: ", err)
	}

	server := grpc.NewServer(grpc.Creds(transportCredentials))
	grpcapi.RegisterAuthServiceServer(server, grpcapi.NewAuthServer(service.NewService(repository.NewRepo(db))))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		server.GracefulStop()
	}()

	log.Printf("auth gRPC server listening on %s", address)
	if err := server.Serve(listener); err != nil {
		log.Fatal(fmt.Errorf("gRPC server stopped: %w", err))
	}
}
