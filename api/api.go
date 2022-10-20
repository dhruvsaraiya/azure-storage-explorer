package api

import (
	"azure-storage-explorer/internal/blob"
	"context"
	"encoding/json"
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

	r.HandleFunc("/containers/{container}/{prefix:.*}", handler(a.handleListBlobs)).Name("listBlobs")
	r.HandleFunc("/containers/{container}", handler(a.handleListBlobs)).Name("listBlobs")
	r.HandleFunc("/containers", handler(a.handleListContainers)).Name("listContainers")

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

func writeJSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, v interface{}) error {
	retJSON, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	_, err = fmt.Fprint(w, string(retJSON))
	if err != nil {
		return err
	}

	return err
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
}

func (a *API) handleListContainers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	containers, err := a.blobService.ListContainers(ctx)
	if err != nil {
		writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
	}
	if err = writeJSONResponse(ctx, w, http.StatusOK, containers); err != nil {
		writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
	}
	return nil
}

func (a *API) handleListBlobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	container := vars["container"]
	prefix := vars["prefix"]
	// If onlyPrefix is true we will return only prefix names,
	// otherwise we will return only blob names
	onlyPrefix := r.URL.Query().Get("op") == "1"

	fmt.Printf("container: %s, prefix: %s, onlyPrefix: %t", container, prefix, onlyPrefix)
	var values []string
	var err error
	if onlyPrefix {
		values, err = a.blobService.ListPrefixes(ctx, container, prefix)
	} else {
		values, err = a.blobService.ListBlobs(ctx, container, prefix)
	}
	if err != nil {
		writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
	}
	if err = writeJSONResponse(ctx, w, http.StatusOK, values); err != nil {
		writeErrorResponse(ctx, w, http.StatusInternalServerError, err)
	}
	return nil
}
