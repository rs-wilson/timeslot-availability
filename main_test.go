package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/rs-wilson/timeslot-availability/request"
	"github.com/rs-wilson/timeslot-availability/response"
	"github.com/stretchr/testify/require"
)

var fullTestAddress = fmt.Sprintf("http://localhost:%s/v1/timeslot", ServerPort)

func TestMain_EndToEnd(t *testing.T) {
	// Ensure clean test run
	exec.Command("pkill timeslot-server").Run()

	// server cmd
	cmd := exec.Command("./timeslot-server")
	t.Cleanup(func() {
		if cmd.Process != nil {
			cmd.Process.Kill() //ensure it's dead, even on a bad test run
		}
	})

	// Start server, but don't wait. logs will be in the log.txt file
	cmd.Start()

	// Give server time to start up
	time.Sleep(1 * time.Second)

	// Test Vars
	startTime := time.Now()
	offsetTime := startTime.Add(time.Minute * 30)
	slotDur := time.Hour

	// Perform & verify
	CheckAvailable(t, startTime, slotDur, true)

	// Ensure reserving works
	ReserveAvailable(t, startTime, slotDur)

	// Ensure the reservation is now blocked
	CheckAvailable(t, startTime, slotDur, false)

	// Ensure overlap works
	CheckAvailable(t, startTime, slotDur, false)

	// Error check on unavailable reservation
	ReserveUnavailable(t, offsetTime, slotDur)

	// Error check on freeing a non-slot
	FreeNoSlot(t, offsetTime, slotDur)

	// Ensure we can free our original reservation
	FreeSlot(t, startTime, slotDur)

	// Ensure the free actually worked
	CheckAvailable(t, startTime, slotDur, true)

	// Ensure we can reserve the new time
	ReserveAvailable(t, offsetTime, slotDur)
}

func CheckAvailable(t *testing.T, slot time.Time, slotDur time.Duration, expected bool) {
	js := request.NewTimeslotRequest(slot, slotDur).ToJson()
	reqBody := strings.NewReader(string(js))

	res, err := http.Post(fullTestAddress, "application/json", reqBody)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	availableRes, err := response.AvailabilityResponseFromJson(bodyBytes)
	require.NoError(t, err)

	require.Equal(t, expected, availableRes.Available)
}

func ReserveAvailable(t *testing.T, slot time.Time, slotDur time.Duration) {
	js := request.NewTimeslotRequest(slot, slotDur).ToJson()
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("PUT", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 200)
}

func ReserveUnavailable(t *testing.T, slot time.Time, slotDur time.Duration) {
	js := request.NewTimeslotRequest(slot, slotDur).ToJson()
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("PUT", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 409)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	errRes, err := response.ErrorResponseFromJson(bodyBytes)
	require.NoError(t, err)

	require.NotEmpty(t, errRes.Message)
}

func FreeNoSlot(t *testing.T, slot time.Time, slotDur time.Duration) {
	js := request.NewTimeslotRequest(slot, slotDur).ToJson()
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("DELETE", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 404)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	errRes, err := response.ErrorResponseFromJson(bodyBytes)
	require.NoError(t, err)

	require.NotEmpty(t, errRes.Message)
}

func FreeSlot(t *testing.T, slot time.Time, slotDur time.Duration) {
	js := request.NewTimeslotRequest(slot, slotDur).ToJson()
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("DELETE", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 204)
}
