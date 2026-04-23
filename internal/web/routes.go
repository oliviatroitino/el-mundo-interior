package web

import (
	"net/http"

	"el-mundo-interior/internal/content"
	"el-mundo-interior/internal/users"
	"el-mundo-interior/internal/web/handlers"
)

func (s *Server) routes() http.Handler {
	// Repositories: each one receives the DB connection.
	postRepo := content.NewPostRepository(s.db)
	userRepo := users.NewUserRepository(s.db)

	// SessionStore: in-memory map token -> userID.
	sessions := handlers.NewSessionStore()

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("GET /", handlers.Home)

	// Worlds
	mux.HandleFunc("GET /mundos/{slug}", handlers.WorldBySlug(postRepo, sessions))
	mux.HandleFunc("POST /mundos/{slug}", handlers.CreatePost(postRepo, sessions))

	// Sections
	mux.HandleFunc("GET /mundos/{slug}/{section}", handlers.WorldSectionBySlug(postRepo, sessions))
	mux.HandleFunc("POST /mundos/{slug}/{section}", handlers.CreateSectionPost(postRepo, sessions))

	// Auth
	mux.HandleFunc("GET /registro", handlers.Register(userRepo))
	mux.HandleFunc("POST /registro", handlers.Register(userRepo))
	mux.HandleFunc("GET /login", handlers.Login(userRepo, sessions))
	mux.HandleFunc("POST /login", handlers.Login(userRepo, sessions))
	mux.HandleFunc("POST /logout", handlers.Logout(sessions))

	// Contact
	mux.HandleFunc("GET /contacto", handlers.Contact())
	mux.HandleFunc("POST /contacto", handlers.Contact())

	// Static files
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	return mux
}
