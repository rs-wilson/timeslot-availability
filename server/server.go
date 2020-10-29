package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs-wilson/timeslot-availability/request"
	"github.com/rs-wilson/timeslot-availability/response"
)

func NewTimeslotServer(listenAddr string, tsStore TimeslotStore) *TimeslotServer {
	// Setup server router
	r := chi.NewRouter()
	server := &TimeslotServer{
		addr:  listenAddr,
		store: tsStore,
		mux:   r,
	}

	// Define endpoints & handlers
	r.Get("/v1/timeslot", server.AvailabilityHandler)
	r.Put("/v1/timeslot", server.ReserveHandler)
	r.Delete("/v1/timeslot", server.FreeHandler)

	return server
}

type TimeslotStore interface {
	IsAvailable(time.Time, time.Duration) bool
	Reserve(time.Time, time.Duration) (bool, error)
	Delete(time.Time, time.Duration) (bool, error)
}

type TimeslotServer struct {
	addr  string
	store TimeslotStore

	mux *chi.Mux
}

func (ts *TimeslotServer) ListenAndServe() error {
	return http.ListenAndServe(ts.addr, ts.mux)
}

func (ts *TimeslotServer) AvailabilityHandler(w http.ResponseWriter, req *http.Request) {
	slot, dur, err := request.ExtractSlotQueryParams(req)
	if err != nil {
		log.Printf("invalid request recieved: %s", err)
		js := response.NewErrorResponse("invalid request").ToJson()
		http.Error(w, js, 400)
		return
	}

	isAvailable := ts.store.IsAvailable(slot, dur)
	aRes := response.NewAvailabilityResponse(isAvailable)
	w.Write([]byte(aRes.ToJson())) //200 OK
}

func (ts *TimeslotServer) ReserveHandler(w http.ResponseWriter, req *http.Request) {
	tsReq, ok := extractTimeslotRequest(w, req)
	if !ok {
		return
	}
	slot, dur := tsReq.ToTime()

	isAvailable, err := ts.store.Reserve(slot, dur)
	if !isAvailable {
		js := response.NewErrorResponse("the requested time slot is not available").ToJson()
		http.Error(w, js, 409)
		return
	}
	if err != nil {
		log.Printf("error reserving timeslot: %s", err.Error())
		js := response.NewErrorResponse("internal error").ToJson()
		http.Error(w, js, 500)
		return
	}

	w.WriteHeader(http.StatusOK) //200
}

func (ts *TimeslotServer) FreeHandler(w http.ResponseWriter, req *http.Request) {
	tsReq, ok := extractTimeslotRequest(w, req)
	if !ok {
		return
	}
	slot, dur := tsReq.ToTime()

	ok, err := ts.store.Delete(slot, dur)
	if err != nil {
		log.Printf("error deleting timeslot: %s", err.Error())
		js := response.NewErrorResponse("internal error").ToJson()
		http.Error(w, js, 500)
	}

	if !ok {
		js := response.NewErrorResponse("rquested timeslot is not a full timeslot entry").ToJson()
		http.Error(w, js, 404)
	}

	w.WriteHeader(http.StatusNoContent) //204
}

// returns a timeslot request and an OK value. if false, this method has writen an error.
func extractTimeslotRequest(w http.ResponseWriter, req *http.Request) (*request.TimeslotRequest, bool) {
	bb, err := ioutil.ReadAll(req.Body)
	if err != nil {
		js := response.NewErrorResponse("improper timeslot request body").ToJson()
		http.Error(w, js, 400)
		return nil, false
	}

	tsReq, err := request.TimeslotRequestFromJson(bb)
	if err != nil {
		js := response.NewErrorResponse("improper timeslot request body").ToJson()
		http.Error(w, js, 400)
		return nil, false
	}

	return tsReq, true
}
