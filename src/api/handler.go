package api

import (
	"github.com/MetalRex101/affise/src/network"
	"io/ioutil"
	"net/http"
)

type Urls []string

func NewHandler(s Script) *handler {
	return &handler{s: s}
}

type handler struct {
	s Script
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		network.WriteResponse(http.StatusInternalServerError, network.Response{Error: "failed to read request body"}, w)
	}

	resp, scriptErr := h.s.Run(r.Context(), body)
	if scriptErr != nil {
		network.WriteResponse(scriptErr.Code, network.Response{Error: scriptErr.Message}, w)
	} else {
		network.WriteResponse(http.StatusOK, network.Response{Data: resp}, w)
	}
}
