package main

import (
	"api-scorekeeper/server"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/bdharris08/scorekeeper"
	"github.com/bdharris08/scorekeeper/score"
	"github.com/bdharris08/scorekeeper/store"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	exampleDSN = "postgres://postgres:xxx@localhost:5432/postgres"
	dsn        = flag.String("dsn", exampleDSN, "dsn for postgres database")
)

func getDSN(dsnFlag *string) (string, error) {
	flag.Parse()
	if *dsn != exampleDSN {
		return *dsn, nil
	}

	d := os.Getenv("DATABASE_URL")
	if d != "" {
		return d, nil
	}

	return "", fmt.Errorf("dsn must be provided by --dsn or en var 'DATABASE_URL'")
}

func main() {
	DSN, err := getDSN(dsn)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	st, _ := store.NewSQLStore(db)
	factory := score.ScoreFactory{
		"trial": func() score.Score { return &score.Trial{} },
	}

	scoreKeeper, err := scorekeeper.New(st, factory)
	if err != nil {
		panic(fmt.Errorf("error creating scoreKeeper: %v", err))
	}

	scoreKeeper.Start()
	defer scoreKeeper.Stop()

	s := server.NewServer(scoreKeeper)
	s.ListenAndServe(":3000")
}
