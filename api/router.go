package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"address-issuer/pkg/acapy"
	"address-issuer/pkg/config"
	"address-issuer/pkg/utils"

	"github.com/gorilla/mux"
)

//go:embed build
var embeddedFiles embed.FS

func NewRouter(config *config.Config, acapy *acapy.Client, cache *utils.BigCache) *mux.Router {
	r := mux.NewRouter()

	path := r.PathPrefix("/api/v2/address-issuer").Subrouter()
	path.HandleFunc("/credential", getCredential(config, acapy, cache))
	path.HandleFunc("/email-credential", getCredentialByEmail(config, acapy, cache))
	path.HandleFunc("/topic/{topic}/", webhookEvents(config, acapy, cache))

	r.PathPrefix("/").Handler(http.FileServer(getFileSystem()))

	return r
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		fmt.Println(err)
	}

	return http.FS(fsys)
}
