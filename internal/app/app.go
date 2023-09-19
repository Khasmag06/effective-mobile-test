package app

import (
	"context"
	"fmt"
	"github.com/khasmag06/effective-mobile-test/config"
	"github.com/khasmag06/effective-mobile-test/internal/controller/api"
	"github.com/khasmag06/effective-mobile-test/internal/repo/people/cache"
	peopleRepo "github.com/khasmag06/effective-mobile-test/internal/repo/people/postgres"
	"github.com/khasmag06/effective-mobile-test/internal/service/people"
	"github.com/khasmag06/effective-mobile-test/internal/webapi"
	"github.com/khasmag06/effective-mobile-test/pkg/httpserver"
	"github.com/khasmag06/effective-mobile-test/pkg/kafka"
	"github.com/khasmag06/effective-mobile-test/pkg/logger"
	"github.com/khasmag06/effective-mobile-test/pkg/postgres"
	"github.com/khasmag06/effective-mobile-test/pkg/redis"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l, err := logger.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}
	defer func() { _ = l.Sync() }()

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.PG)
	if err != nil {
		l.Fatalf("failed to connect to postgres db: %s", err)
	}
	defer db.Close()

	redisDB, err := redis.ConnectRedis(ctx, cfg.Redis)
	if err != nil {
		l.Fatalf("failed to connect to postgres redis db: %s", err)
	}

	repo := peopleRepo.New(db.Pool)
	peopleCache := cache.New(redisDB, repo, l)
	service := people.New(peopleCache)

	fioInfoApi := webapi.New(cfg.PersonApi, service, l)

	kafkaClient, err := kafka.NewKafkaClient(cfg.Kafka.BrokerURLs)
	if err != nil {
		l.Fatalf("failed to create kafka client: %v", err)
	}
	defer kafkaClient.Close()
	consumeTopic := cfg.Kafka.FioTopic
	messages, errors := kafkaClient.ConsumeFromTopic(consumeTopic)

	go func() {
		for {
			select {
			case msg := <-messages:
				l.Infof("Received message from topic %s: %s", consumeTopic, string(msg.Value))

				if err := fioInfoApi.AddFioData(msg.Value); err != nil {
					l.Error(err.Error())
					errorMessage := fmt.Sprintf("%s: %s", err.Error(), string(msg.Value))
					kafkaClient.SendMessageToTopic(cfg.Kafka.FioFailedTopic, []byte(errorMessage))
					continue
				}
			case err := <-errors:
				l.Errorf("failed to read message from topic %s: %v", consumeTopic, err)
			}
		}
	}()

	// HTTP Server
	l.Info("Starting api server...")
	handler := api.NewHandler(service, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Errorf("app - Run - httpServer.Notify: %w", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Errorf("app - Run - httpServer.Shutdown: %w", err)
	}
}
