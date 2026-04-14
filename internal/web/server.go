// PROMPT: Dame una estructura básica de un servidor HTTP en Go, 
// con una función para iniciar el servidor y otra para definir las rutas. 
// El servidor debe escuchar en un puerto específico y manejar solicitudes HTTP básicas.

// Package web contiene la lógica del servidor HTTP y manejo de rutas.
package web

import (
	"net/http"
	"time"
)

// Server encapsula la configuración del servidor HTTP.
type Server struct {
	addr string           // Dirección en la que escucha (ej: ":8080")
	http *http.Server    // Instancia del servidor HTTP de Go
}

func NewServer(addr string) *Server {
	server := &Server{addr: addr}
	server.http = &http.Server{
		Addr:              addr,
		Handler:           server.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return server
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
