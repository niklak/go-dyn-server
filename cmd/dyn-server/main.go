package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/niklak/go-dyn-server/load"
)

const logPrefix string = "DYN-SERVER"

type config struct {
	Host         string        `env:"HOST"`
	Port         int           `env:"PORT" envDefault:"8080"`
	Addr         string        `env:"ADDR,expand" envDefault:"$HOST:${PORT}"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"120s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"120s"`
	PluginRoot   string        `env:"PLUGIN_ROOT"`
}

func setupPluginHandles(r chi.Router, serverPlugins *load.ServerPlugins) {

	for _, pre := range serverPlugins.Middlewares {
		r.Use(pre)
	}

	for _, pluginHandle := range serverPlugins.Handles {
		for _, method := range pluginHandle.Methods {
			r.MethodFunc(method, pluginHandle.Route, pluginHandle.Handle)
		}
	}
}

func main() {
	var err error

	cfg := config{}
	opts := env.Options{Prefix: "SERVER_"}
	if err = env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("[ERROR] %s: %v\n", logPrefix, err)
	}

	log.Printf("[INFO] %s: CONFIG: %+v\n", logPrefix, cfg)

	serverPlugins, err := load.NewServerPlugins(cfg.PluginRoot)
	if err != nil {
		log.Fatalf("[ERROR] %s: %v\n", logPrefix, err)
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("[INFO] GOOS: %s, GOARCH: %s\n", os.Getenv("GOOS"), os.Getenv("GOARCH"))

	r := chi.NewRouter()

	setupPluginHandles(r, &serverPlugins)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[INFO][%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      r,
	}

	go func() {
		log.Printf("[INFO] %s: START SERVING ON %s\n", logPrefix, cfg.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("[WARNING] %s: %v\n", logPrefix, err)
			}
		}
		log.Printf("[INFO] %s: STOPPED SERVING\n", logPrefix)
	}()

	<-stop
	log.Printf("[INFO] %s: shutting down...\n", logPrefix)
	if err = srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("[ERROR] %s: shutdown %v\n", logPrefix, err)
	}

	log.Printf("[INFO] %s: gracefully stopped\n", logPrefix)
}
