package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/attwad/nikki/server"
	"github.com/attwad/nikki/store"
	"github.com/attwad/nikki/store/initdb"
	"github.com/attwad/nikki/store/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

var (
	addr   = flag.String("addr", "127.0.0.1:8080", "addr:port to listen on")
	dbPath = flag.String("db_path", "db.sqlite", "path to database")
)

func main() {
	flag.Parse()

	var runInitDDL bool
	if _, err := os.Stat(*dbPath); errors.Is(err, os.ErrNotExist) {
		runInitDDL = true
	}

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	defer db.Close()
	if runInitDDL {
		if err := initdb.MustInitDB(ctx, db); err != nil {
			log.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	store := &sqlite.SQLiteStore{
		Queries: store.New(db),
	}
	server.InitMux(mux, "/api", store)
	s := &http.Server{
		Addr:           *addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}
	fmt.Println("Serving at", s.Addr)
	log.Fatal(s.ListenAndServe())
}
