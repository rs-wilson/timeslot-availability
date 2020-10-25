package response

import (
	"encoding/json"
	"log"
)

func AvailabilityResponseFromJson(js []byte) (*AvailabilityResponse, error) {
	aRes := &AvailabilityResponse{}
	err := json.Unmarshal(js, aRes)
	return aRes, err
}

func NewAvailabilityResponse(a bool) *AvailabilityResponse {
	return &AvailabilityResponse{Available: a}
}

type AvailabilityResponse struct {
	Available bool `json:"available"`
}

func (a *AvailabilityResponse) ToJson() string {
	js, err := json.Marshal(a)
	jsStr := string(js)
	if err != nil {
		log.Printf("failed to marshal availability response: %s", err.Error())
		jsStr = "internal error"
	}
	return jsStr
}
