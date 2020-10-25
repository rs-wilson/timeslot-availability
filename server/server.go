package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewTimeslotServer(listenAddr string) *TimeslotServer {
	// Setup server router
	r := chi.NewRouter()
	server := &TimeslotServer{
		addr: listenAddr,
		mux:  r,
	}

	// Define endpoints & handlers
	r.Post("/v1/timeslot", server.AvailabilityHandler)
	r.Put("/v1/timeslot", server.ReserveHandler)
	r.Delete("/v1/timeslot", server.FreeHandler)

	return server
}

type TimeslotServer struct {
	addr string
	mux  *chi.Mux
}

func (ts *TimeslotServer) ListenAndServe() error {
	return http.ListenAndServe(ts.addr, ts.mux)
}

func (ts *TimeslotServer) AvailabilityHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "internal error", 500)
}

func (ts *TimeslotServer) ReserveHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "internal error", 500)
}

func (ts *TimeslotServer) FreeHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "internal error", 500)
}
