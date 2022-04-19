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
	exampleDSN    = "postgres://postgres:xxx@localhost:5432/postgres"
	dsn           = flag.String("dsn", exampleDSN, "dsn for postgres database")
	routes        = flag.Bool("routes", false, "Generate router documentation")
	listenAddress = flag.String("listen", ":3000", "address to listen on")
)

// primitive config setup
// A production app might use a custom config package
// 	that enforces precedence between configuration sources
//  which could be flags, env vars,
//	or even remote sources like parameter store or secrets manager.
type config struct {
	dsn           string
	listenAddress string
}

func getConfig() config {
	conf := config{
		listenAddress: ":3000",
	}

	flag.Parse()

	switch {
	case *dsn != exampleDSN:
		conf.dsn = *dsn
	case os.Getenv("DATABASE_URL") != "":
		conf.dsn = os.Getenv("DATABASE_URL")
	}

	switch {
	case *listenAddress != ":3000":
		conf.listenAddress = *listenAddress
	case os.Getenv("LISTEN") != "":
		conf.listenAddress = os.Getenv("LISTEN")
	}

	return conf
}

func main() {
	conf := getConfig()

	db, err := sql.Open("pgx", conf.dsn)
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

	if *routes {
		fmt.Println(s.Usage())
		return
	}

	s.ListenAndServe(conf.listenAddress)
}
