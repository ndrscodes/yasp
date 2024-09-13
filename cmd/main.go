package main

import (
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ndrscodes/yasp/internal/handlers"
	"github.com/ndrscodes/yasp/internal/web"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tpl, err := template.ParseFS(web.Files, "templates/layout/*")
	if err != nil {
		log.Fatalf("Unable to parse templates: %v", err)
	}

	mux := createMux(tpl)

	log.Println("Starting server on port 8000...")
	log.Fatalln(http.ListenAndServe(":8000", mux))
}

func createMux(root *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	hoh, err := handlers.NewHomeHandler(root)
	if err != nil {
		log.Fatalf("Unable to create home handler: %v", err)
	}

	mux.HandleFunc("GET /", hoh.HandleGet)

	pub, err := fs.Sub(web.Files, "static")
	if err != nil {
		log.Fatalf("Unable to navigate to embedded static directory: %v", err)
	}

	fs := http.FileServer(http.FS(pub))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	return mux
}
