package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/ndrscodes/yasp/internal/db"
	"github.com/ndrscodes/yasp/internal/handlers"
	"github.com/ndrscodes/yasp/internal/web/static"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conn, err := createDb()
	if err != nil {
		slog.Error("Unable to create database connection: %v", "error", err)
		os.Exit(1)
	}

	migrator, err := db.NewMigrator(conn, true)
	if err != nil {
		slog.Error("Unable to create migrator: %v", "error", err)
		os.Exit(1)
	}

	err = migrator.Up()
	if err != nil {
		slog.Error("Unable to migrate database: %v", "error", err)
		os.Exit(1)
	}

	sys, err := db.NewSystemRepository(conn)
	if err != nil {
		slog.Error("Unable to create system repository: %v", "error", err)
		os.Exit(1)
	}

	s, err := sys.GetAll()
	if err != nil {
		slog.Error("Unable to fetch systems", "error", err)
		os.Exit(1)
	}
	log.Printf("Fetched %v systems from the database", s)

	log.Println("Starting server on port 8000...")

	mux := createMux(conn)
	log.Fatalln(http.ListenAndServe(":8000", mux))
}

func createDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:yasp@localhost:5432/yasp?sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("established database connection")

	return db, nil
}

func createMux(d *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	sys, err := db.NewSystemRepository(d)
	if err != nil {
		log.Fatalf("Unable to create system repository: %v", err)
	}

	hoh, err := handlers.NewHomeHandler(&sys)
	if err != nil {
		log.Fatalf("Unable to create home handler: %v", err)
	}

	mux.HandleFunc("GET /", hoh.HandleGet)

	fs := http.FileServer(http.FS(static.Static))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	return mux
}
