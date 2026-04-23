package web

import (
	"net/http"

	"el-mundo-interior/internal/contact"
	"el-mundo-interior/internal/content"
	"el-mundo-interior/internal/users"
	"el-mundo-interior/internal/web/handlers"
)

func (s *Server) routes() http.Handler {
	// Repositorios: cada uno recibe la conexión a la BD
	postRepo := content.NewPostRepository(s.db)
	userRepo := users.NewUserRepository(s.db)
	contactRepo := contact.NewRepository(s.db)

	// SessionStore: mapa en memoria token → userID, compartido por todos los handlers
	sessions := handlers.NewSessionStore()

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("GET /", handlers.Home)

	// Mundos
	mux.HandleFunc("GET /mundos/{slug}", handlers.WorldBySlug(postRepo, sessions))
	mux.HandleFunc("POST /mundos/{slug}", handlers.CreatePost(postRepo, sessions))

	// Secciones
	mux.HandleFunc("GET /mundos/{slug}/{section}", handlers.WorldSectionBySlug(postRepo, sessions))
	mux.HandleFunc("POST /mundos/{slug}/{section}", handlers.CreateSectionPost(postRepo, sessions))

	// Autenticación
	mux.HandleFunc("GET /registro", handlers.Register(userRepo))
	mux.HandleFunc("POST /registro", handlers.Register(userRepo))
	mux.HandleFunc("GET /login", handlers.Login(userRepo, sessions))
	mux.HandleFunc("POST /login", handlers.Login(userRepo, sessions))
	mux.HandleFunc("POST /logout", handlers.Logout(sessions))

	// Contacto
	mux.HandleFunc("GET /contacto", handlers.Contact(contactRepo))
	mux.HandleFunc("POST /contacto", handlers.Contact(contactRepo))

	// Archivos estáticos
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	return mux
}
