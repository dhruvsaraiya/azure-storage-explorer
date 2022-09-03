package api

import (
	"azure-storage-explorer/internal/blob"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type API struct {
	router      *mux.Router
	logger      *zap.Logger
	blobService *blob.BlobService
}

func NewAPI(logger *zap.Logger, blobService *blob.BlobService) (*API, error) {
	r := mux.NewRouter()
	r.HandleFunc("/_ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "PONG\n")
	})

	return &API{
		router:      r,
		logger:      logger,
		blobService: blobService,
	}, nil
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
