package users

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"log"
)

var (
	DisplayNameAlreadyExistsErr = fmt.Errorf("display name is already taken")
)

type DisplayNameClaimer struct {
	Connection Connection
}

func (i *Ensurer) ClaimDisplayName(
	ctx context.Context, userID string, displayName string,
) error {
	conn, db, err := i.Connection.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("GetConnection: %s", err.Error())
	}

	getExisting := goqu.New("mysql", db).
		Select("user_id").
		From(goqu.T("User")).
		Where(goqu.C("displayname").Eq(displayName)).
		Limit(1)

	query, _, err := getExisting.ToSQL()
	if err != nil {
		return fmt.Errorf("ToSQL: %s", err.Error())
	}

	r, err := conn.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("QueryContext: %s", err.Error())
	}

	for r.Next() {
		return DisplayNameAlreadyExistsErr
	}

	in := goqu.New("mysql", db).Update(goqu.T("Users")).Set(
		goqu.Record{"displayname": displayName},
	).Where(goqu.C("id").Eq(userID))

	updateSQL, _, err := in.ToSQL()
	if err != nil {
		return fmt.Errorf("in.ToSQL: %s", err.Error())
	}

	log.Println("Inserting")
	_, err = conn.ExecContext(ctx, updateSQL)
	if err != nil {
		return fmt.Errorf("ExecContext: %s", err.Error())
	}

	return nil
}
