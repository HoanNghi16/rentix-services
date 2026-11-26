package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/HoanNghi16/rentix-services/employee-service/internal/database"
	"github.com/HoanNghi16/rentix-services/employee-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/employee-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/employee-service/internal/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("cannot connect employee database: ", err)
	}
	port := os.Getenv("EMPLOYEE_GRPC_PORT")
	if port == "" {
		port = "8083"
	}
	address := ":" + port
	transportCredentials, err := grpcapi.LoadServerTLSCredentials()
	if err != nil {
		log.Fatal("cannot configure employee gRPC mutual TLS: ", err)
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot listen for employee gRPC: ", err)
	}
	server := grpc.NewServer(grpc.Creds(transportCredentials))
	employeeService := service.NewService(repository.NewRepo(db))
	grpcapi.RegisterEmployeeServiceServer(server, grpcapi.NewEmployeeServer(employeeService))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		server.GracefulStop()
	}()
	log.Printf("employee gRPC server listening on %s", address)
	if err := server.Serve(listener); err != nil {
		log.Fatal(fmt.Errorf("employee gRPC server stopped: %w", err))
	}
}
