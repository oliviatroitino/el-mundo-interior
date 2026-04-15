// Package handlers contiene los view models: structs de datos que se pasan
// a cada template para que pueda renderizarse. No contienen lógica, solo datos.
package handlers

import "el-mundo-interior/internal/content"

// NavItem es un elemento dentro de un desplegable del nav.
type NavItem struct {
	Href   string
	Label  string
	Icon   string // ruta a imagen, vacío si no aplica
	Active bool   // true si es la página actualmente activa
}

// NavDropdown es un desplegable del nav con su título y sus items.
// SummaryIcon permite mostrar una imagen en el botón del desplegable (e.g. planeta activo).
// Class es la clase CSS del <details>; si está vacío se usa "nav__dropdown".
type NavDropdown struct {
	Label       string
	SummaryIcon string // ruta a imagen mostrada en el <summary>, vacío si no aplica
	Class       string // clase CSS del <details>, vacío → "nav__dropdown"
	Items       []NavItem
}

// NavLink es un enlace simple (no desplegable) en el nav.
type NavLink struct {
	Href  string
	Label string
}

// NavData agrupa todos los elementos de la barra de navegación.
// Cada handler construye el suyo propio para que el nav sea reutilizable.
// HomeHref: si está relleno, se muestra el logo como enlace a esa URL.
type NavData struct {
	HomeHref  string // p.e. "/" en páginas internas; vacío en la home
	Dropdowns []NavDropdown
	Links     []NavLink
}

// Post representa una publicación de usuario.
type Post struct {
	User     string
	Title    string
	Text     string
	Location string
	Date     string
}

// HomePlanetItem contiene los datos de un mundo para la lista de planetas en la home.
// IsReverse indica si la card debe mostrarse con imagen a la derecha (layout alternado).
type HomePlanetItem struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	IsReverse   bool
}

// ReviewItem representa una valoración de usuario en la home.
type ReviewItem struct {
	Stars  string // e.g. "★★★★★"
	Text   string
	Author string
}

// PlanItem representa un plan de suscripción en la home.
type PlanItem struct {
	ID       string // id HTML, vacío si no aplica
	Name     string
	Features []string
}

// HomePageData contiene los datos para la página de inicio (/).
type HomePageData struct {
	Nav     NavData
	Worlds  []HomePlanetItem
	Reviews []ReviewItem
	Plans   []PlanItem
}

// WorldPageData contiene los datos para la página de un mundo (/mundos/{slug}).
type WorldPageData struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	Sections    []content.WorldSection
	Nav         NavData
	MyPosts     []Post
	OtherPosts  []Post
}

// SectionPageData contiene los datos para la página de una sección (/mundos/{slug}/{section}).
type SectionPageData struct {
	World   content.World
	Section content.WorldSection
	Nav     NavData
}

// RegisterPageData contiene los datos para el formulario de registro (/registro).
// Error se muestra si el envío del formulario falla.
// Name y Email se usan para repoblar los campos si hay error, para no perder lo escrito.
type RegisterPageData struct {
	Error string
	Name  string
	Email string
}

// LoginPageData contiene los datos para el formulario de inicio de sesión (/login).
type LoginPageData struct {
	Error string
	Email string
}

// ContactPageData contiene los datos para el formulario de contacto (/contacto).
type ContactPageData struct {
	Success bool   // true si el mensaje se envió correctamente
	Error   string // mensaje de error si algo falló
}
