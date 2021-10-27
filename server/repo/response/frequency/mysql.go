package frequency

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func NewMySQLAmplitudeRepo() *MySQLAmplitudeRepo {
	repo := &MySQLAmplitudeRepo{}
	repo.Inserter = &Inserter{Connection: repo}
	return repo
}

type MySQLAmplitudeRepo struct {
	*Inserter
	ipAddr string
	conn   *sql.Conn
	db     *sql.DB
}

func (m *MySQLAmplitudeRepo) GetConnection(ctx context.Context) (*sql.Conn, *sql.DB, error) {
	if m.conn != nil {
		err := m.conn.PingContext(ctx)
		if err != nil {
			log.Printf("Ping failed: %s", err.Error())
		} else {
			return m.conn, nil, nil
		}
	}

	c, err := m.db.Conn(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("db.Conn: %s", err.Error())
	}
	m.conn = c

	return m.conn, m.db, nil
}

func (m *MySQLAmplitudeRepo) Initialize(ctx context.Context) error {
	f, err := os.Open("/dbcreds.txt")
	if err != nil {
		return fmt.Errorf("creds: %s", err.Error())
	}

	rd := bufio.NewReader(f)
	dbUser, _, err := rd.ReadLine()
	if err != nil {
		return fmt.Errorf("ReadLine(dbUser): %s", err.Error())
	}

	dbPass, _, err := rd.ReadLine()
	if err != nil {
		return fmt.Errorf("ReadLine(dbPass): %s", err.Error())
	}

	cfg := mysql.NewConfig()
	cfg.User = string(dbUser)
	cfg.Passwd = string(dbPass)
	cfg.Net = "tcp"
	cfg.Addr = "190.92.153.141"
	cfg.DBName = "car_av_db"

	m.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("sql.Open: %s", err.Error())
	}
	return nil
}
