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

	a := &API{
		router:      r,
		logger:      logger,
		blobService: blobService,
	}

	r.HandleFunc("/list-containers", handler(a.handleListContainers))

	return a, nil
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// handleFunc is an http.HandlerFunc that can return an error
type handleFunc func(w http.ResponseWriter, r *http.Request) error

func handler(hf handleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := hf(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a *API) handleListContainers(w http.ResponseWriter, r *http.Request) error {
	a.blobService.ListContainers()
	return nil
}
