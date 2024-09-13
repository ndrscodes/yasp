package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ndrscodes/yasp/internal/handlers"
	"github.com/ndrscodes/yasp/internal/static"
	"github.com/ndrscodes/yasp/internal/templates"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tpl, err := template.ParseFS(templates.Files, "layout/*")
	if err != nil {
		log.Fatalf("Unable to parse templates: %v", err)
		os.Exit(1)
	}

	mux := createMux(tpl)

	log.Fatalln(http.ListenAndServe(":8000", mux))
}

func createMux(root *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	hoh, err := handlers.NewHomeHandler(root)
	if err != nil {
		log.Fatalf("Unable to create home handler: %v", err)
		os.Exit(1)
	}

	mux.HandleFunc("GET /", hoh.HandleGet)

	fs := http.FileServer(http.FS(static.Static))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	return mux
}
