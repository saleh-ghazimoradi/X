package main

import (
	"github.com/saleh-ghazimoradi/X/config"
	"github.com/saleh-ghazimoradi/X/migrations"
	"github.com/saleh-ghazimoradi/X/utils"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("failed to load config", "err", err.Error())
		os.Exit(1)
	}

	postgresql := utils.NewPostgresql(
		utils.WithHost(cfg.Postgresql.Host),
		utils.WithPort(cfg.Postgresql.Port),
		utils.WithUser(cfg.Postgresql.User),
		utils.WithPassword(cfg.Postgresql.Password),
		utils.WithName(cfg.Postgresql.Name),
		utils.WithTimeout(cfg.Postgresql.Timeout),
		utils.WithSSLMode(cfg.Postgresql.SSLMode),
		utils.WithMaxOpenConn(cfg.Postgresql.MaxOpenConn),
		utils.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
		utils.WithMaxIdleConn(cfg.Postgresql.MaxIdleConn),
	)

	postgresqlDB, err := postgresql.Connect()
	if err != nil {
		logger.Error("failed to connect to postgresql", "err", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := postgresqlDB.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	migrate, err := migrations.NewMigrate(postgresqlDB, postgresql.Name)
	if err != nil {
		logger.Error("failed to load migrations", "err", err.Error())
		os.Exit(1)
	}

	if err := migrate.UP(); err != nil {
		logger.Error("failed to up migrations", "err", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := migrate.Close(); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

}
