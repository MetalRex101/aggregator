package network

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func WriteResponse(status int, response Response, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	body, err := json.Marshal(response)
	if err != nil {
		fmt.Println(fmt.Errorf("faield to marshal response: %w", err))
	}

	if _, err := w.Write(body); err != nil {
		fmt.Println(fmt.Errorf("failed to write response: %w", err))
	}
}
