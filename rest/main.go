package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reshmavatkar/kv-store/rest/client"
	"github.com/reshmavatkar/kv-store/rest/handler"
)

func main() {
	grpcAddr := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051" // fallback
	}
	storeClient, err := client.NewStoreClient(grpcAddr)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer storeClient.Close()

	h := handler.NewHandler(storeClient)

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())
	//TODO: add middleware to handle recovery, logging, and other validation
	router.PUT("/store", h.PutValue)
	router.GET("/store/:key", h.GetValue)
	router.DELETE("/store/:key", h.DeleteValue)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("REST API server running on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down REST API...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
