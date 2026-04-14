// Package handlers contiene los view models: structs de datos que se pasan
// a cada template para que pueda renderizarse. No contienen lógica, solo datos.
package handlers

import "el-mundo-interior/internal/content"

// WorldNavItem representa un mundo en la barra de navegación superior.
// Se usa en todas las páginas que muestran el menú de mundos.
type WorldNavItem struct {
	Slug   string
	Title  string
	Icon   string
	Active bool // true si es el mundo que se está visitando actualmente
}

// HomePageData contiene los datos para la página de inicio (/).
type HomePageData struct {
	Worlds []content.World // los 6 mundos para mostrar las cards
}

// WorldPageData contiene los datos para la página de un mundo (/mundos/{slug}).
type WorldPageData struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	Sections    []content.WorldSection
	Worlds      []WorldNavItem // para el nav
}

// SectionPageData contiene los datos para la página de una sección (/mundos/{slug}/{section}).
type SectionPageData struct {
	World   content.World
	Section content.WorldSection
	Worlds  []WorldNavItem // para el nav
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
