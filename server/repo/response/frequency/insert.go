package frequency

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
)

type Connection interface {
	GetConnection(ctx context.Context) (*sql.Conn, *sql.DB, error)
}

type Inserter struct {
	Connection Connection
}

func (i *Inserter) Create(
	ctx context.Context, createdByUserId string, data []byte,
) error {
	id := uuid.New().String()

	conn, db, err := i.Connection.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("GetConnection: %s", err.Error())
	}

	in := goqu.New("mysql", db).Insert(goqu.T("FreqResponse")).Rows(
		goqu.Record{"id": id, "created_by_user_id": createdByUserId, "data_json": data},
	)

	insertSQL, _, err := in.ToSQL()
	if err != nil {
		return fmt.Errorf("in.ToSQL: %s", err.Error())
	}

	fmt.Println("Inserting")
	_, err = conn.ExecContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("ExecContext: %s", err.Error())
	}
	fmt.Println("Done")

	return nil
}
