package content

import "time"

// Post representa una publicación de usuario almacenada en la base de datos.
// UserName se rellena con un JOIN al consultar, para no necesitar una query extra.
type Post struct {
	ID          int
	UserID      int
	UserName    string // nombre del autor, obtenido con JOIN
	WorldSlug   string
	SectionSlug string
	Title       string
	Body        string
	Location    string
	CreatedAt   time.Time
}

// PostRepository define las operaciones de acceso a datos para posts.
type PostRepository interface {
	// GetByWorld devuelve todos los posts de un mundo, ordenados por fecha descendente.
	GetByWorld(worldSlug string) ([]Post, error)
	// GetBySection devuelve los posts de una sección concreta.
	GetBySection(worldSlug, sectionSlug string) ([]Post, error)
	// Create inserta un nuevo post y devuelve su ID.
	Create(post Post) (int, error)
}
