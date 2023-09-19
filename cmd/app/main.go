package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/khasmag06/effective-mobile-test/config"
	"github.com/khasmag06/effective-mobile-test/internal/app"
	"log"
)

// @Title FIOService API
// @Description Сервис хранения персональных данных.
// @Version 1.0

// @BasePath /api
// @Host localhost:8080

func main() {
	// configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// migrate
	m, err := migrate.New("file://migrations", cfg.PG.URL)
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %v", err)
	}

	// run
	app.Run(cfg)

}
