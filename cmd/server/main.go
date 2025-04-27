package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "webhook-service/internal/api"
    "webhook-service/internal/cache"
    "webhook-service/internal/config"
    "webhook-service/internal/queue"
    "webhook-service/internal/store/mongo"
    "webhook-service/internal/store/postgres"
    "webhook-service/internal/worker"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    pgDB, err := postgres.NewPostgres(cfg.PostgresURL)
    if err != nil {
        log.Fatalf("postgres connect: %v", err)
    }
    subRepo := postgres.NewSubscriptionRepo(pgDB)

    mongoClient, err := mongo.NewMongoClient(cfg.MongoURL)
    if err != nil {
        log.Fatalf("mongo connect: %v", err)
    }
    logRepo := mongo.NewLogRepo(mongoClient, cfg.MongoDBName)

    redisClient := cache.NewRedis(cfg.RedisURL)

    rabbitConn, err := queue.NewRabbitMQ(cfg.RabbitURL)
    if err != nil {
        log.Fatalf("rabbitmq connect: %v", err)
    }
    publisher, err := queue.NewPublisher(rabbitConn)
    if err != nil {
        log.Fatalf("publisher init: %v", err)
    }

    deliveryWorker := worker.NewWorker(subRepo, logRepo, publisher, redisClient, cfg)
    go deliveryWorker.Start()

    r := gin.Default()
    api.RegisterSubscriptionRoutes(r, subRepo)
    api.RegisterIngestRoutes(r, publisher)

    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: r,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen error: %v", err)
        }
    }()
    log.Printf("HTTP server running on port %s", cfg.Port)

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("server shutdown failed: %v", err)
    }
}
