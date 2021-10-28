package users

import (
	"context"
	"database/sql"
)

type Connection interface {
	GetConnection(ctx context.Context) (*sql.Conn, error)
}

type Ensurer struct {
	Connection Connection
}

func (i *Ensurer) Ensure(
	ctx context.Context, googleId string,
) error {

	// TODO: Get user by google ID and return

	// TODO: If none exist, create user and google->user mapping (in one TX?)

	return nil
}
