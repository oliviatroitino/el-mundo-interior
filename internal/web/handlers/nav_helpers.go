package handlers

import (
	"el-mundo-interior/internal/content"
)

// buildWorldDropdown construye el dropdown de mundos para el nav.
// activeSlug marca el mundo actualmente activo (vacío si no aplica).
func buildWorldDropdown(activeSlug string) NavDropdown {
	d := NavDropdown{Label: "Mundos"}
	for _, w := range content.OrderedWorlds() {
		d.Items = append(d.Items, NavItem{
			Href:   "/mundos/" + w.Slug,
			Label:  w.Title,
			Icon:   w.Icon,
			Active: w.Slug == activeSlug,
		})
	}
	return d
}

// buildSectionDropdown construye el dropdown de secciones de un mundo para el nav.
// activeSlug marca la sección actualmente activa (vacío si no aplica).
func buildSectionDropdown(sections []content.WorldSection, activeSlug string) NavDropdown {
	d := NavDropdown{Label: "Secciones"}
	for _, s := range sections {
		d.Items = append(d.Items, NavItem{
			Href:   s.Path,
			Label:  s.Title,
			Active: s.Slug == activeSlug,
		})
	}
	return d
}

// authLinks devuelve los enlaces de autenticación comunes (Entrar / Registro).
func authLinks() []NavLink {
	return []NavLink{
		{Href: "/login", Label: "Entrar"},
		{Href: "/registro", Label: "Registro"},
	}
}
