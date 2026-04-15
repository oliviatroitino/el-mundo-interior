package handlers

import (
	"el-mundo-interior/internal/content"
)

// buildWorldDropdown construye el dropdown de mundos para el nav.
// activeSlug marca el mundo actualmente activo (vacío si no aplica).
// Cuando hay un mundo activo, el label del dropdown muestra su nombre y SummaryIcon su icono.
func buildWorldDropdown(activeSlug string) NavDropdown {
	d := NavDropdown{Label: "Mundos"}
	for _, w := range content.OrderedWorlds() {
		if w.Slug == activeSlug {
			d.Label = w.Title
			d.SummaryIcon = w.Icon
		}
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
// Cuando hay una sección activa, el label del dropdown muestra su nombre.
func buildSectionDropdown(sections []content.WorldSection, activeSlug string) NavDropdown {
	d := NavDropdown{Label: "Secciones"}
	for _, s := range sections {
		if s.Slug == activeSlug {
			d.Label = s.Title
		}
		d.Items = append(d.Items, NavItem{
			Href:   s.Path,
			Label:  s.Title,
			Active: s.Slug == activeSlug,
		})
	}
	return d
}

// buildUserDropdown construye el dropdown de usuario (icono de persona) para páginas internas.
func buildUserDropdown() NavDropdown {
	return NavDropdown{
		Class:       "nav__userdropdown",
		SummaryIcon: "/assets/icons/users.svg",
		Items: []NavItem{
			{Href: "#", Label: "Perfil"},
			{Href: "#", Label: "Cerrar sesión"},
		},
	}
}
