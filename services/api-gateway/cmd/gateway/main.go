package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/authclient"
	"github.com/HoanNghi16/rentix-services/api-gateway/internal/handler"
	"github.com/HoanNghi16/rentix-services/api-gateway/internal/propertyclient"
	"github.com/gin-gonic/gin"
)

func main() {
	auth, err := authclient.NewFromEnv()
	if err != nil {
		log.Fatal("cannot configure auth gRPC client: ", err)
	}
	defer auth.Close()
	property, err := propertyclient.NewFromEnv()
	if err != nil {
		log.Fatal("cannot configure property gRPC client: ", err)
	}
	defer property.Close()

	accessSecret := os.Getenv("ACCESS_SECRET")
	if strings.TrimSpace(accessSecret) == "" {
		log.Fatal("ACCESS_SECRET is required to protect property routes")
	}

	gin.EnableJsonDecoderDisallowUnknownFields()
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	handler.RegisterRoutes(router, auth)
	handler.RegisterPropertyRoutes(router, property, accessSecret)
	server := &http.Server{
		Addr:              envOrDefault("GATEWAY_HTTP_ADDR", ":8080"),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	log.Printf("API gateway listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("gateway HTTP server stopped: ", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
