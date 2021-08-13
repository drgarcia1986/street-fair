package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router *mux.Router
	port   int
	log    *logrus.Logger
}

func (s *Server) Run() error {
	server := &http.Server{
		Handler: s.Router,
		Addr:    fmt.Sprintf(":%d", s.port),
	}
	s.log.WithField("port", s.port).Info("Starting Server")
	return server.ListenAndServe()
}

func NewServer(port int, log *logrus.Logger) *Server {
	return &Server{
		port:   port,
		log:    log,
		Router: mux.NewRouter().StrictSlash(true),
	}
}
