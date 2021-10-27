package frequency

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/bradsk88/CarAudioDatabase/server/model/frequency"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"log"
)

type Getter struct {
	Connection Connection
}

func (i *Inserter) Get(
	ctx context.Context, id string,
) ([]model.DataPoint, error) {
	conn, db, err := i.Connection.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetConnection: %s", err.Error())
	}

	in := goqu.New("mysql", db).
		Select("data_json").
		From(goqu.T("FreqResponse")).
		Where(goqu.C("id").Eq(id)).
		Limit(1)

	sql, _, err := in.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("ToSQL: %s", err.Error())
	}

	r, err := conn.QueryContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("QueryContext: %s", err.Error())
	}

	var result []byte
	for r.Next() {
		err = r.Scan(&result)
		if err != nil {
			return nil, fmt.Errorf("scan: %s", err.Error())
		}
		log.Printf("Read %d bytes from FR\n", len(result))
	}

	var out []model.DataPoint
	err = json.Unmarshal(result, &out)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %s", err.Error())
	}

	return out, nil
}
