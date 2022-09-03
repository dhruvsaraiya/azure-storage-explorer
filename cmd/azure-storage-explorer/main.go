package main

import (
	"azure-storage-explorer/api"
	"azure-storage-explorer/internal/blob"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(shutdown)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to create logger", zap.Error(err))
	}
	defer logger.Sync()
	logger.Info("Starting service...")

	blobService, err := blob.NewBlobService(ctx, logger)
	if err != nil {
		logger.Fatal("failed to create blob service", zap.Error(err))
	}

	api, err := api.NewAPI(logger, blobService)
	if err != nil {
		logger.Fatal("failed to create api", zap.Error(err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: api,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-shutdown
		logger.Info("Shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("failed to start server", zap.Error(err))
	}

	wg.Wait()
}
