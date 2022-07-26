package api

import (
	"address_issuer/pkg/client"
	"address_issuer/pkg/config"
	"address_issuer/pkg/utils"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed build
var embeddedFiles embed.FS

func NewRouter(config *config.Config, client *client.Client, cache *utils.BigCache) *mux.Router {
	r := mux.NewRouter()

	path := r.PathPrefix("/api/v1/address-issuer").Subrouter()
	path.HandleFunc("/health", health(config))
	path.HandleFunc("/connections", listConnections(config, client))
	path.HandleFunc("/credentials", listCredentials(config, client))
	path.HandleFunc("/credential", getCredential(config, client, cache))
	path.HandleFunc("/credential-email", getCredentialByEmail(config, client, cache))
	path.HandleFunc("/topic/{topic}/", webhookEvents(config, client, cache))

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
