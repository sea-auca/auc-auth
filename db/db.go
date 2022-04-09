package db

import (
	"errors"
	"fmt"

	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
	_ "github.com/lib/pq"
	"github.com/sea-auca/auc-auth/config"
	"github.com/talkanbaev-artur/shutdown"
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
