package load

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"plugin"
)

const (
	handleDir     = "handlers"
	middlewareDir = "middlewares"
)

type PluginHandle struct {
	Route   string
	Methods []string
	Handle  func(w http.ResponseWriter, r *http.Request)
}

type ServerPlugins struct {
	Handles     []PluginHandle
	Middlewares []func(http.Handler) http.Handler
}

func filePathWalkDir(root string) ([]string, error) {

	files, err := filepath.Glob(path.Join(root, "*.so"))
	return files, err
}

func readHandle(fp string) (ph PluginHandle, err error) {
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

func listHandles(root string) (handles []PluginHandle, err error) {
	files, err := filePathWalkDir(root)
	if err != nil {
		return
	}
	for _, fp := range files {
		handle, err := readHandle(fp)
		if err != nil {
			return nil, err
		}

		handles = append(handles, handle)
	}
	return
}

func readMiddleware(fp string) (pm func(http.Handler) http.Handler, err error) {
	p, err := plugin.Open(fp)
	if err != nil {
		return
	}
	rawMiddleware, err := p.Lookup("Middleware")
	if err != nil {
		return
	}
	pm = rawMiddleware.(func(http.Handler) http.Handler)
	return
}

func listMiddlewares(root string) (middlewares []func(http.Handler) http.Handler, err error) {
	files, err := filePathWalkDir(root)
	if err != nil {
		return
	}
	for _, fp := range files {
		middleware, err := readMiddleware(fp)
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, middleware)
	}
	return
}

func NewServerPlugins(root string) (serverPlugins ServerPlugins, err error) {

	handles, err := listHandles(path.Join(root, handleDir))
	if err != nil {
		return
	}

	preMiddlewares, err := listMiddlewares(path.Join(root, middlewareDir))
	if err != nil {
		return
	}

	serverPlugins = ServerPlugins{
		Handles:     handles,
		Middlewares: preMiddlewares,
	}
	return
}

//TODO: add a web handler that accepts a handler
