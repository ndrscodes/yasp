package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ndrscodes/yasp/internal/handlers"
	"github.com/ndrscodes/yasp/internal/util/templates"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tr, err := templates.NewTemplateRegistry("templates/*/*.html")
	if err != nil {
		slog.Error("unable to set up template registry", "error", err)
		panic(err)
	}

	mux := createMux(tr)

	log.Fatalln(http.ListenAndServe(":8000", mux))
}

func createMux(tr *templates.TemplateRegistry) *http.ServeMux {
	mux := http.NewServeMux()
	hoh := handlers.NewHomeHandler(tr, "index.html")
	mux.HandleFunc("GET /", hoh.HandleGet)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	return mux
}
