package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rs-wilson/timeslot-availability/server"
	"github.com/rs-wilson/timeslot-availability/storage"
)

const ServerPort = "6543" //TODO: make configurable

func main() {
	err := run()
	if err != nil {
		log.Printf("Error returned from main.run: %s", err.Error())
	}
}

func run() error {
	//setup logging
	err := setupLogging("log.txt") //TODO: make configurable
	if err != nil {
		return errors.Wrap(err, "failed to setup logging")
	}

	//run & error check
	// setup server
	addr := fmt.Sprintf(":%s", ServerPort)
	store := storage.NewInMemoryTimeslotStore()
	ts := server.NewTimeslotServer(addr, store)

	// run server
	return ts.ListenAndServe()
}

func setupLogging(filename string) error {
	// Create server output file
	logfile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create server log file")
	}

	// Setup multi-writer for logging to both locations
	mw := io.MultiWriter(os.Stderr, logfile)
	log.SetOutput(mw)
	return nil
}
