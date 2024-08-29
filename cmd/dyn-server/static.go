//go:build static

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/niklak/go-dyn-server/load"
)

var trueClientIP = "True-Client-IP"
var xForwardedFor = "X-Forwarded-For"
var xRealIP = "X-Real-IP"

// HeadersResponse is the response for the headers endpoint
type HeadersResponse struct {
	Headers http.Header `json:"headers"`
}

// IpResponse is the response for the ip endpoint
type IpResponse struct {
	Origin string `json:"origin"`
}

func init() {
	setupHandles = setupStaticHandles
}

func setupStaticHandles(r chi.Router, serverPlugins *load.ServerPlugins) {
	log.Printf("[INFO] using static handles!\n")

	r.Use(Cors)

	r.Get("/headers", http.HandlerFunc(HeadersHandle))
	r.Get("/ip", http.HandlerFunc(IpHandle))
}

// HeadersHandle returns only the request headers. Check `HeadersResponse`.
func HeadersHandle(w http.ResponseWriter, r *http.Request) {
	writeJsonResponse(w, http.StatusOK, HeadersResponse{Headers: getRequestHeader(r)})
}

// IpHandle returns only the IP address. Check `IpResponse`.
func IpHandle(w http.ResponseWriter, r *http.Request) {

	writeJsonResponse(w, http.StatusOK, IpResponse{Origin: getIP(r)})
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

func getRequestHeader(r *http.Request) http.Header {
	h := r.Header.Clone()
	h.Set("Host", r.Host)
	return h
}

func writeJsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}

//Cors middleware

// Cors middleware to handle Cross Origin Resource Sharing (CORS).
func Cors(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
				for _, v := range acrh {
					w.Header().Add("Access-Control-Allow-Headers", v)
				}
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
