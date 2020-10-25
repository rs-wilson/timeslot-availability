package request

import (
	"encoding/json"
	"time"
)

func TimeslotRequestFromJson(js []byte) (*TimeslotRequest, error) {
	tsReq := &TimeslotRequest{}
	err := json.Unmarshal(js, tsReq)
	return tsReq, err
}

func NewTimeslotRequest(start time.Time, d time.Duration) *TimeslotRequest {
	return &TimeslotRequest{
		StartTimestamp: start.Unix(),
		Duration:       int64(d.Seconds()), //ignore sub-seconds
	}
}

type TimeslotRequest struct {
	StartTimestamp int64 `json:"start_timestamp"`
	Duration       int64 `json:"duration"`
}

func (r *TimeslotRequest) ToJson() (string, error) {
	js, err := json.Marshal(r)
	return string(js), err
}

func (r *TimeslotRequest) ToTime() (time.Time, time.Duration) {
	startTime := time.Unix(r.StartTimestamp, 0)
	dur := time.Duration(r.Duration)
	return startTime, dur
}
