package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
	"log"
)

type Connection interface {
	GetConnection(ctx context.Context) (*sql.Conn, *sql.DB, error)
}

type Ensurer struct {
	Connection Connection
}

func (i *Ensurer) Ensure(
	ctx context.Context, googleId string, googleEmailAddress string,
) (userId string, err error) {
	conn, db, err := i.Connection.GetConnection(ctx)
	if err != nil {
		return "", fmt.Errorf("GetConnection: %s", err.Error())
	}

	getExisting := goqu.New("mysql", db).
		Select("user_id").
		From(goqu.T("GoogleToUser")).
		Where(goqu.C("google_id").Eq(googleId)).
		Limit(1)

	query, _, err := getExisting.ToSQL()
	if err != nil {
		return "", fmt.Errorf("ToSQL: %s", err.Error())
	}

	r, err := conn.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("QueryContext: %s", err.Error())
	}

	var result []byte
	for r.Next() {
		err = r.Scan(&result)
		if err != nil {
			return "", fmt.Errorf("scan: %s", err.Error())
		}
		log.Printf("Read %d bytes from FR\n", len(result))
	}

	if result != nil {
		return string(result), nil
	}

	userId = uuid.New().String()

	tx, err := goqu.New("mysql", db).BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("BeginTx: %s", err.Error())
	}

	in := tx.Insert(goqu.T("Users")).Rows(
		goqu.Record{"id": userId, "email_address": googleEmailAddress},
	)

	insertSQL, _, err := in.ToSQL()
	if err != nil {
		return "", fmt.Errorf("in.ToSQL: %s", err.Error())
	}

	log.Println("Inserting")
	_, err = tx.ExecContext(ctx, insertSQL)
	if err != nil {
		return "", fmt.Errorf("ExecContext: %s", err.Error())
	}

	googlIn := tx.Insert(goqu.T("GoogleToUser")).Rows(
		goqu.Record{"google_id": googleId, "user_id": userId},
	)

	gInsertSQL, _, err := googlIn.ToSQL()
	if err != nil {
		return "", fmt.Errorf("gIn.ToSQL: %s", err.Error())
	}

	log.Println("Inserting")
	_, err = tx.ExecContext(ctx, gInsertSQL)
	if err != nil {
		return "", fmt.Errorf("ExecContext: %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("tx.Commit: %s", err.Error())
	}
	log.Println("Done")

	return userId, nil
}
