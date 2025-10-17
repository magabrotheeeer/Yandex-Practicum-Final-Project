package server

import (
	"net/http"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/api"
)

func New(webDir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux)

	return mux
}
