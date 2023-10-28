// What we're trying to do here:
// Create a server - use gorilla/mux for handling the requests
// Wrap the thing neatly into a *net/http.Server - so that we can use ListenAndServe to easily serve
// the data.

// Our httpServer wraps our Log object and exposes POST and GET calls on the root (/) handler to
// Produce and Consume the logs respectively.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type logHandler struct {
	log *Log
}

// Create a new log handler that will be used to bind the http routes to produce and consume
// records from the append-only-log
func newLogHandler() *logHandler {
	return &logHandler{
		log: NewLog(),
	}
}

// The respective types
type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (l *logHandler) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	off, err := l.log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	res := ProduceResponse{Offset: off}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
func (l *logHandler) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	record, err := l.log.Read(req.Offset)
	// We have a defined error here, so let's use that
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	res := ConsumeResponse{Record: record}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func NewHTTPServer(addr string) *http.Server {
	lH := newLogHandler()
	r := mux.NewRouter()

	r.HandleFunc("/", lH.handleProduce).Methods("POST")
	r.HandleFunc("/", lH.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
