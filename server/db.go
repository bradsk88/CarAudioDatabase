package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func dbTest(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("/dbcreds.txt")
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Creds: %s", err.Error()))
		return
	}

	rd := bufio.NewReader(f)
	dbUser, _, err := rd.ReadLine()
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Read(user): %s", err.Error()))
		return
	}

	dbPass, _, err := rd.ReadLine()
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Read(pass): %s", err.Error()))
		return
	}

	cfg := mysql.NewConfig()
	cfg.User = string(dbUser)
	cfg.Passwd = string(dbPass)
	cfg.Net = "tcp"
	cfg.Addr = "190.92.153.141"
	cfg.DBName = "car_av_db"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Open: %s", err.Error()))
		return
	}
	c, err := db.Conn(r.Context())
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Conn: %s", err.Error()))
		return
	}

	in := goqu.Dialect("mysql").Insert(goqu.T("Test")).Rows(
		goqu.Record{"Name": "Brad"},
	)

	insertSQL, _, err := in.ToSQL()
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Insert.ToSQL: %s", err.Error()))
		return
	}

	fmt.Println(insertSQL)

	_, err = c.ExecContext(r.Context(), insertSQL)
	if err != nil {
		w.WriteHeader(500)
		log.Println(fmt.Sprintf("Exec: %s", err.Error()))
		return
	}

}
