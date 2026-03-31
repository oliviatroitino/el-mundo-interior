package web

import (
	"net/http"
)

type Server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return &Server{
		addr: addr,
		mux:  mux,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.mux)
}
