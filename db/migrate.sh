#/usr/bin/bash

POSTGRES_HOST='localhost'
POSTGRES_PORT='5432'
POSTGRES_USERNAME='postgres'
POSTGRES_PASSWORD='postgres'
POSTGRES_DB='auc_auth_dev'
DSN=$(printf 'postgres://%s:%s@%s:%s/%s?sslmode=disable' $POSTGRES_USERNAME $POSTGRES_PASSWORD $POSTGRES_HOST $POSTGRES_PORT $POSTGRES_DB)

function migrate() {
    rel migrate -adapter=github.com/go-rel/rel/adapter/postgres -driver=github.com/lib/pq -dsn="${DSN}"
}

function rollback {
    rel rollback -adapter=github.com/go-rel/rel/adapter/postgres -driver=github.com/lib/pq -dsn="${DSN}"
}

$@