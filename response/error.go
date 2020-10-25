package response

import "encoding/json"

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

func (e *ErrorResponse) ToJson() (string, error) {
	js, err := json.Marshal(e)
	return string(js), err
}
