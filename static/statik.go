package static

import (
	"net/http"

	_ "github.com/kimkit/appdaemon/static/statik"
	"github.com/rakyll/statik/fs"
)

func NewHandler(prefix string) (http.Handler, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	return http.FileServer(statikFS), nil
}
