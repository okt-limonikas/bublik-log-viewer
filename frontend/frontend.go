package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/okt-limonikas/bublik-log-viewer/internal/constants"
)

//go:embed build/*
var BuildFs embed.FS

func BuildHTTPFS() http.FileSystem {
	build, err := fs.Sub(BuildFs, constants.BuildPath)

	if err != nil {
		log.Fatal(err)
	}

	return http.FS(build)
}
