package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jantytgat/go-sql-queryrepo/pkg/queryrepo"
)

//go:embed assets/migrations/*
var migrationFs embed.FS

//go:embed assets/statements/*
var statementsFS embed.FS

var db *sql.DB

func main() {
	var err error

	fmt.Println("Create in-memory database and run migrations:")
	// Create in-memory demo database
	if err = initialize(); err != nil {
		panic(err)
	}

	fmt.Println("Fetch demo list query from filesystem:")
	var query string
	if query, err = queryrepo.LoadQueryFromFs(statementsFS, "assets/statements", "demo", "list"); err != nil {
		panic(err)
	}
	fmt.Println("Query:", query)

	if query, err = queryrepo.LoadQueryFromFs(statementsFS, "assets/statements", "demo", "insert"); err != nil {
		panic(err)
	}
	fmt.Println("Query:", query)
	fmt.Println("")

	fmt.Println("Run insert statement:")
	var stmtInsert *sql.Stmt
	if stmtInsert, err = queryrepo.PrepareFromFs(db, statementsFS, "assets/statements", "demo", "insert"); err != nil {
		panic(err)
	}

	var res int
	if err = stmtInsert.QueryRow("item2").Scan(&res); err != nil {
		panic(err)
	}
	fmt.Println("Successfully inserted into database, returned id:", res)

	fmt.Println("")
	fmt.Println("Run list statement:")
	var stmtQuery *sql.Stmt
	if stmtQuery, err = queryrepo.PrepareFromFs(db, statementsFS, "assets/statements", "demo", "list"); err != nil {
		fmt.Println("Error preparing statement")
	}

	if stmtQuery == nil {
		panic(errors.New("statement not found"))
	}

	var rows *sql.Rows
	if rows, err = stmtQuery.Query(); err != nil {
		fmt.Println("Error querying statement")
	}
	if rows != nil {
		for rows.Next() {
			var id int
			var name string
			if err = rows.Scan(&id, &name); err != nil {
				fmt.Println("Error scanning row")
			}
			fmt.Println("Output", id, name)
		}
	}

}

func initialize() error {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		fmt.Println("Error opening database")
		return err
	}

	var src source.Driver
	if src, err = iofs.New(migrationFs, "assets/migrations"); err != nil {
		fmt.Println("Error opening migrations source")
		return err
	}

	var driver database.Driver
	if driver, err = sqlite.WithInstance(db, &sqlite.Config{}); err != nil {
		fmt.Println("Error opening migrations destination")
		return err
	}

	var m *migrate.Migrate
	if m, err = migrate.NewWithInstance("fs", src, "sqlite", driver); err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
