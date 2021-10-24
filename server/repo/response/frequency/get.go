package frequency

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"log"
)

type Getter struct {
	Connection Connection
}

func (i *Inserter) Get(
	ctx context.Context, id string,
) error {
	conn, err := i.Connection.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("GetConnection: %s", err.Error())
	}

	in := goqu.Dialect("mysql").
		Select("data_json").
		From(goqu.T("FreqResponse")).
		Where(goqu.C("id").Eq(id)).
		Limit(1)

	sql, _, err := in.ToSQL()
	if err != nil {
		return fmt.Errorf("ToSQL: %s", err.Error())
	}

	r, err := conn.QueryContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("QueryContext: %s", err.Error())
	}

	var result string
	for r.Next() {
		err = r.Scan(&result)
		if err != nil {
			return fmt.Errorf("scan: %s", err.Error())
		}
		log.Println(result)
		log.Printf("Read %d bytes from FR\n", len(result))
	}

	return nil
}
