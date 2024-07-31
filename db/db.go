package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	context    context.Context
	connection *pgx.Conn
}

func (s *Store) Connect() *pgx.Conn {
	s.context = context.Background()
	conn, err := pgx.Connect(s.context, "somedatabaseurl")
	if err != nil {
		log.Fatal("could not connect to database")
		os.Exit(1)
	}
	s.connection = conn
	return conn
}

func (s *Store) Close() {
	s.connection.Close(s.context)
}
