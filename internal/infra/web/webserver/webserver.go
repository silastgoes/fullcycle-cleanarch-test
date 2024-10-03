package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]*HandleMethod
	WebServerPort string
}

type HandleMethod struct {
	Handler http.HandlerFunc
	Method  string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]*HandleMethod),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers[path] = &HandleMethod{
		Handler: handler,
		Method:  method,
	}
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Method(handler.Method, path, handler.Handler)

	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
