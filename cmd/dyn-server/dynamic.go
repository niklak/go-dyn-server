//go:build !static

package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/niklak/go-dyn-server/load"
)

func init() {
	setupHandles = setupDynamicHandles
}

func setupDynamicHandles(r chi.Router, serverPlugins *load.ServerPlugins) {
	log.Printf("[INFO] using dynamic handles!")

	for _, pre := range serverPlugins.Middlewares {
		r.Use(pre)
	}

	for _, pluginHandle := range serverPlugins.Handles {
		for _, method := range pluginHandle.Methods {
			r.MethodFunc(method, pluginHandle.Route, pluginHandle.Handle)
		}
	}
}
