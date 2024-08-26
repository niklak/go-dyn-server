package main

import (
	"net/http"
	"strings"

	"github.com/niklak/go-dyn-server/helpers"
)

const (
	trueClientIP  = "True-Client-IP"
	xForwardedFor = "X-Forwarded-For"
	xRealIP       = "X-Real-IP"
)

var Route = "/ip"

var Methods = []string{"GET"}

// IpHandle returns only the IP address. Check `IpResponse`.

// IpResponse is the response for the ip endpoint
type IpResponse struct {
	Origin string `json:"origin"`
}

func getIP(r *http.Request) string {
	var ip string

	if realIP := r.Header.Get(xRealIP); realIP != "" {
		ip = realIP
	} else if clientIP := r.Header.Get(trueClientIP); clientIP != "" {
		ip = clientIP
	} else if forwardedFor := r.Header.Get(xForwardedFor); forwardedFor != "" {
		ip = strings.TrimSpace(strings.SplitN(forwardedFor, ",", 2)[0])
	} else {
		ip = r.RemoteAddr
	}
	return ip
}

func Handle(w http.ResponseWriter, r *http.Request) {
	data := IpResponse{Origin: getIP(r)}
	helpers.JsonResponse(w, http.StatusOK, data)
}
