package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"assecor.assessment.test/internal/data"
	"assecor.assessment.test/internal/validator"
	_ "github.com/duckdb/duckdb-go/v2"
)

const version = "1.0.0"

type config struct {
	port int
	dsn  string
}

// Application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.dsn, "dsn", "", "Data source name")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	if len(cfg.dsn) > 0 {
		v := validator.New()
		app.models.LoadFromCsv(v, cfg.dsn)
		if !v.Valid() {
			for k, m := range v.Errors {
				app.logger.Printf("[%s]: %s\n", k, m)
			}
		}
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}

// Create internal database
func openDB(_ config) (*sql.DB, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	query := `
		CREATE SEQUENCE seq_personid START 1;
		CREATE TABLE IF NOT EXISTS persons (
			id INTEGER PRIMARY KEY DEFAULT NEXTVAL('seq_personid'),	
			name TEXT NOT NULL,
			lastname TEXT NOT NULL,
			zipcode TEXT NOT NULL,
			city TEXT NOT NULL,
			color INTEGER NOT NULL)`
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDB()
	_, err = db.ExecContext(ctxDB, query)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
