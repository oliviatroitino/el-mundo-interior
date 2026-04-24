package handlers

import (
	"el-mundo-interior/internal/content"
)

// buildWorldDropdown builds the worlds dropdown for the nav.
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

// buildSectionDropdown builds the sections dropdown for a world.
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

// buildUserDropdown builds the user dropdown for logged-in users.
func buildUserDropdown(userName string) NavDropdown {
	return NavDropdown{
		Class:       "nav__userdropdown",
		SummaryIcon: "/assets/icons/users.svg",
		Items: []NavItem{
			{Href: "#", Label: "¡Hola, " + userName + "!"},
			{Href: "/logout", Label: "Cerrar sesión", Method: "POST"},
		},
	}
}

