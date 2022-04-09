package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/sea-auca/auc-auth/config"
	"github.com/talkanbaev-artur/shutdown"
	"go.uber.org/zap"
)

var ErrNoCredentials = errors.New("no credentials were supplied for production environment")

func ConnectDatabase(shutdown *shutdown.Shutdown) (rel.Repository, error) {
	conf := config.Config().Database

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	adapter, err := postgres.Open(dsn)
	if err != nil {
		return nil, err
	}
	shutdown.Add(adapter.Close) // add close to the global shutdown
	repo := rel.New(adapter)

	return repo, nil
}

func ConnectPGXDatabase(ctx context.Context) *pgxpool.Pool {
	conf := config.Config().Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		zap.L().Fatal("Failed to configure the PGX driver connection", zap.Error(err))
	}

	logger, _ := zap.NewDevelopment(zap.DebugLevel)
	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	conn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to the Postgres Database with PGX driver", zap.Error(err))
	}
	return conn
}
