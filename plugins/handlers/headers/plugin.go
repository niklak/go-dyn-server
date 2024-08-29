package main

import (
	"net/http"

	"github.com/niklak/go-dyn-server/helpers"
)

var Route = "/headers"

var Methods = []string{"GET"}

// HeadersResponse is the response for the headers endpoint
type HeadersResponse struct {
	Headers http.Header `json:"headers"`
}

func getRequestHeader(r *http.Request) http.Header {
	h := r.Header.Clone()
	h.Set("Host", r.Host)
	return h
}

// Handle returns only the request headers. Check `HeadersResponse`.
func Handle(w http.ResponseWriter, r *http.Request) {
	helpers.JsonResponse(w, http.StatusOK, HeadersResponse{Headers: getRequestHeader(r)})
}
