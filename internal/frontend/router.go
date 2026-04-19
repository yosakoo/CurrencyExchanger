package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web
var webFiles embed.FS

func NewRouter(mux *http.ServeMux) {
	sub, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatalf("frontend: failed to sub web FS: %v", err)
	}

	mux.Handle("GET /", http.FileServerFS(sub))
}
