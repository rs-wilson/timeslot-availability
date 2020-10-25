package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rs-wilson/timeslot-availability/server"
)

const ServerPort = "6543" //TODO: make configurable

func main() {
	//setup logging
	log.SetOutput(os.Stderr) //TODO: make configurable

	//run & error check
	err := run()
	if err != nil {
		log.Printf("Error returned from main.run: %s", err.Error())
	}
}

func run() error {
	// setup server
	addr := fmt.Sprintf(":%s", ServerPort)
	ts := server.NewTimeslotServer(addr)

	// run server
	return ts.ListenAndServe()
}
