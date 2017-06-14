package common

import (
	"encoding/json"
	"net/http"
)

// InvalidOpenstackToken Error
type InvalidOpenstackToken struct{}

func (e InvalidOpenstackToken) Error() string {
	return "provided openstack token is not valid"
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// WriteErrorToResponse writes given params to http response
func WriteErrorToResponse(w http.ResponseWriter, code int, message string,
	details string) {

	w.WriteHeader(code)
	payload, _ := json.Marshal(&errorResponse{code, message, details})
	w.Write(payload)
}
