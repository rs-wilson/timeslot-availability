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
	t.Run("Check clear availability", func(t *testing.T) {
		CheckAvailable(t, startTime, slotDur, true)
	})

	t.Run("Check reserving available slot", func(t *testing.T) {
		ReserveAvailable(t, startTime, slotDur)
	})

	t.Run("Ensure reservation is blocked", func(t *testing.T) {
		CheckAvailable(t, startTime, slotDur, false)
	})

	t.Run("Check overlapping time slots", func(t *testing.T) {
		CheckAvailable(t, startTime, slotDur, false)
	})

	t.Run("Ensure error on unavailable reservation", func(t *testing.T) {
		ReserveUnavailable(t, offsetTime, slotDur)
	})

	t.Run("Ensure error on freeing not a full slot", func(t *testing.T) {
		FreeNoSlot(t, offsetTime, slotDur)
	})

	t.Run("Ensure no error on freeing a full slot", func(t *testing.T) {
		FreeSlot(t, startTime, slotDur)
	})

	t.Run("Ensure a deleted slot is actually free", func(t *testing.T) {
		CheckAvailable(t, startTime, slotDur, true)
	})

	t.Run("Ensure no error on reserving over a deleted slot", func(t *testing.T) {
		ReserveAvailable(t, offsetTime, slotDur)
	})
}

func CheckAvailable(t *testing.T, slot time.Time, slotDur time.Duration, expected bool) {
	slotReq := request.NewTimeslotRequest(slot, slotDur) //convert to appropriate timestamps

	addrWithParams := fmt.Sprintf("%s?start_timestamp=%d&duration=%d", fullTestAddress, slotReq.StartTimestamp, slotReq.Duration)
	res, err := http.Get(addrWithParams)
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
