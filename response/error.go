package response

import (
	"encoding/json"
	"log"
)

func ErrorResponseFromJson(js []byte) (*ErrorResponse, error) {
	eRes := &ErrorResponse{}
	err := json.Unmarshal(js, eRes)
	return eRes, err
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func (e *ErrorResponse) ToJson() string {
	js, err := json.Marshal(e)
	jsStr := string(js)
	if err != nil {
		log.Printf("failed to marshal error response: %s", err.Error())
		jsStr = "internal error"
	}
	return jsStr
}
