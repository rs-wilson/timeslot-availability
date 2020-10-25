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

	// Start server
	cmd := exec.Command("./timeslot-server")
	cmd.Start()
	defer func() {
		out, err := cmd.CombinedOutput()
		logs := string(out)
		if err != nil {
			logs = "{failed to retrieve logs}"
		}
		fmt.Printf("\nSERVER LOGS:\n%s\n\n", logs)
		cmd.Process.Kill() //ensure it's dead
	}()

	// Give server time to start up
	time.Sleep(1 * time.Second)

	// Test Vars
	startTime := time.Now()
	offsetTime := startTime.Add(time.Minute * 30)
	slotDur := time.Hour

	// Perform & verify
	CheckAvailable(t, startTime, slotDur, true)

	ReserveAvailable(t, startTime, slotDur)

	CheckAvailable(t, offsetTime, slotDur, false)

	ReserveUnavailable(t, offsetTime, slotDur)

	FreeNoSlot(t, offsetTime, slotDur)

	FreeSlot(t, startTime, slotDur)

	CheckAvailable(t, offsetTime, slotDur, true)

	FreeNoSlot(t, offsetTime, slotDur)
}

func CheckAvailable(t *testing.T, slot time.Time, slotDur time.Duration, expected bool) {
	js, err := request.NewTimeslotRequest(slot, slotDur).ToJson()
	require.NoError(t, err)
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
	js, err := request.NewTimeslotRequest(slot, slotDur).ToJson()
	require.NoError(t, err)
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("PUT", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 200)
}

func ReserveUnavailable(t *testing.T, slot time.Time, slotDur time.Duration) {
	js, err := request.NewTimeslotRequest(slot, slotDur).ToJson()
	require.NoError(t, err)
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("PUT", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 400)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	errRes, err := response.ErrorResponseFromJson(bodyBytes)
	require.NoError(t, err)

	require.NotEmpty(t, errRes.Message)
}

func FreeNoSlot(t *testing.T, slot time.Time, slotDur time.Duration) {
	js, err := request.NewTimeslotRequest(slot, slotDur).ToJson()
	require.NoError(t, err)
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("DELETE", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 400)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	errRes, err := response.ErrorResponseFromJson(bodyBytes)
	require.NoError(t, err)

	require.NotEmpty(t, errRes.Message)
}

func FreeSlot(t *testing.T, slot time.Time, slotDur time.Duration) {
	js, err := request.NewTimeslotRequest(slot, slotDur).ToJson()
	require.NoError(t, err)
	reqBody := strings.NewReader(string(js))

	req, err := http.NewRequest("DELETE", fullTestAddress, reqBody)
	require.NoError(t, err)
	req.Header.Set("Context-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 204)
}
