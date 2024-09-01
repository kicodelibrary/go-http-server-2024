package api

import "encoding/json"

// Response is a response message.
type Response struct {
	Message string `json:"message"`
}

// NewJSONResponse creates a new JSON response.
func NewJSONResponse(message string) []byte {
	r := Response{
		Message: message,
	}
	b, err := json.Marshal(r)
	if err != nil {
		panic(err) // There should be no error here.
	}
	return b
}
