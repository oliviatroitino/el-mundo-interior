// Package web contiene la lógica del servidor HTTP y manejo de rutas.
package web

import (
	"net/http"

	// Importamos el paquete de handlers donde están las funciones que responden a cada ruta
	"el-mundo-interior/internal/web/handlers"
)

// routes() configura todas las rutas de la aplicación.
// Retorna un http.Handler (gestor de rutas) que procesa las peticiones HTTP.
func (s *Server) routes() http.Handler {  // Pertenece a Server (server.go), pero se define la lógica acá por ser rutas
	// Creamos un nuevo gestor de rutas (multiplexer)
	mux := http.NewServeMux()

	// Rutas principales de la aplicación
	mux.HandleFunc("GET /", handlers.Home)              // Página de inicio
	mux.HandleFunc("GET /mundos/{slug}", handlers.WorldBySlug) // Página individual de cada mundo
	mux.HandleFunc("GET /mundos/{slug}/{section}", handlers.WorldSectionBySlug) // Página de cada sección dentro de un mundo

	// Rutas para servir archivos estáticos (CSS e imágenes)
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	return mux
}
