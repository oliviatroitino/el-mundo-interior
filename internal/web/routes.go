package web

import (
	"net/http"

	"el-mundo-interior/internal/contact"
	"el-mundo-interior/internal/content"
	"el-mundo-interior/internal/users"
	"el-mundo-interior/internal/web/handlers"
)

func (s *Server) routes() http.Handler {
	// Repositories: each one receives the DB connection.
	postRepo := content.NewPostRepository(s.db)
	userRepo := users.NewUserRepository(s.db)
	contactRepo := contact.NewRepository(s.db)

	// SessionStore: in-memory map token -> userID.
	sessions := handlers.NewSessionStore()

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("GET /", handlers.Home(sessions))

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

	// Contacto
	mux.HandleFunc("GET /contacto", handlers.Contact(contactRepo))
	mux.HandleFunc("POST /contacto", handlers.Contact(contactRepo))

	// Archivos estáticos — servidos con caché stale-while-revalidate
	mux.Handle("GET /css/", withCache(http.StripPrefix("/css/", http.FileServer(http.Dir("css")))))
	mux.Handle("GET /assets/", withCache(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("assets/uploads"))))

	return mux
}

// withCache añade un header Cache-Control que permite servir contenido
// cacheado mientras se revalida en segundo plano (stale-while-revalidate).
func withCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=86400")
		next.ServeHTTP(w, r)
	})
}
