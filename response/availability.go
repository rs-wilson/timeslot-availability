package response

import "encoding/json"

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

func (a *AvailabilityResponse) ToJson() ([]byte, error) {
	js, err := json.Marshal(a)
	return js, err
}
