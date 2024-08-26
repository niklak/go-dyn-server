package main

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"plugin"
)

type PluginHandle struct {
	Route   string
	Methods []string
	Handle  func(w http.ResponseWriter, r *http.Request)
}

type ServerPlugins struct {
	Handles         []PluginHandle
	PreMiddlewares  []func(http.Handler) http.Handler
	PostMiddlewares []func(http.Handler) http.Handler
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, fpath)
		}
		return nil
	})
	return files, err
}

func loadHandle(fp string) (ph PluginHandle, err error) {
	p, err := plugin.Open(fp)
	if err != nil {
		return
	}

	rawHandle, err := p.Lookup("Handle")
	if err != nil {
		return
	}
	ph.Handle = rawHandle.(func(w http.ResponseWriter, r *http.Request))

	rawRoute, err := p.Lookup("Route")
	if err != nil {
		return
	}
	ph.Route = *(rawRoute.(*string))

	if ph.Route == "" {
		err = errors.New("loadHandle: `Route` value is not found")
		return
	}

	rawMethods, err := p.Lookup("Methods")
	if err != nil {
		return
	}
	ph.Methods = *(rawMethods.(*[]string))

	if len(ph.Methods) == 0 {
		ph.Methods = []string{http.MethodGet}
	}

	return
}

func loadHandles(files []string) (handles []PluginHandle, err error) {
	handles = []PluginHandle{}
	for _, fp := range files {
		handle, err := loadHandle(fp)
		if err != nil {
			return nil, err
		}

		handles = append(handles, handle)
	}
	return
}

func loadPlugins(root string) (serverPlugins ServerPlugins, err error) {
	handlersPath := path.Join(root, "handlers")

	files, err := filePathWalkDir(handlersPath)
	if err != nil {
		return
	}

	serverPlugins = ServerPlugins{}
	serverPlugins.Handles, err = loadHandles(files)
	return
}

//TODO: add a web handler that accepts a handler
