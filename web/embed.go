package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist
var frontendFS embed.FS

// FrontendFS returns the frontend filesystem
func FrontendFS() http.FileSystem {
	fsys, _ := fs.Sub(frontendFS, "frontend/dist")
	return http.FS(fsys)
}
