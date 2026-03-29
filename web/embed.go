package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist
var FrontendFS embed.FS

// GetFrontendFS returns the frontend filesystem as http.FileSystem
func GetFrontendFS() http.FileSystem {
	fsys, _ := fs.Sub(FrontendFS, "frontend/dist")
	return http.FS(fsys)
}
