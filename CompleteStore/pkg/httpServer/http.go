package httpserver

import (
	"fmt"
	"net/http"
)

type Server struct {
	port string
	*http.Server
}

func New(port int) *Server {
	return &Server{
		port: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Start() error {
	//	http.HandleFunc("/store", s.KeyStoreApi)
	return http.ListenAndServe(s.port, nil)
}

func Shutdown() {

}
