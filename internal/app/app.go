package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/sixojke/internal/config"
	delivery "github.com/sixojke/internal/delivery/http"
	"github.com/sixojke/internal/repository"
	"github.com/sixojke/internal/server"
	"github.com/sixojke/internal/service"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("config error: %v", err))
	}

	redis, err := repository.NewRedisDB(cfg.Redis)
	if err != nil {
		log.Fatal(fmt.Sprintf("redis connection error: %v", err))
	}
	defer redis.Close()
	log.Info("[REDIS] Connection successful")

	repo := repository.NewRepository(&repository.Deps{
		Redis:  redis,
		Config: cfg,
	})
	services := service.NewService(&service.Deps{
		Config: cfg,
		Repo:   repo,
	})
	handler := delivery.NewHandler(services)

	srv := server.NewServer(cfg.HTTPServer, handler.Init())
	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %v\n", err)
		}
	}()
	log.Info(fmt.Sprintf("[SERVER] Started :%v", cfg.HTTPServer.Port))

	shutdown(srv, redis)
}

func shutdown(srv *server.Server, redis *redis.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 2 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}

	redis.Close()
}
