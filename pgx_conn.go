package gohelpers

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PgxConn struct {
	ReadConn  *pgx.Conn
	WriteConn *pgx.Conn
}

func NewPgxConfig(ctx context.Context, readURL string, writeURL string) (*PgxConn, error) {

	readConnConfig, err := pgx.ParseConfig(readURL)

	if err != nil {
		return nil, err
	}

	readConn, err := pgx.ConnectConfig(ctx, readConnConfig)

	if err != nil {
		return nil, err
	}

	writeConnConfig, err := pgx.ParseConfig(writeURL)

	if err != nil {
		return nil, err
	}

	writeConn, err := pgx.ConnectConfig(ctx, writeConnConfig)

	if err != nil {
		return nil, err
	}

	return &PgxConn{
		ReadConn:  readConn,
		WriteConn: writeConn,
	}, nil
}
