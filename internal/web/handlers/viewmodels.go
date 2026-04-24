// Package handlers contains view models passed to templates.
package handlers

import "el-mundo-interior/internal/content"

// NavItem is an item inside a nav dropdown.
type NavItem struct {
	Href   string
	Label  string
	Icon   string // image path, empty when not used
	Active bool   // true when this is the current page
	Method string // empty or GET for links, POST for actions
}

// NavDropdown groups nav items under one summary.
type NavDropdown struct {
	Label       string
	SummaryIcon string // image path shown in <summary>, empty when not used
	Class       string // CSS class for <details>, defaults to "nav__dropdown"
	Items       []NavItem
}

// NavLink is a plain (non-dropdown) nav link.
type NavLink struct {
	Href  string
	Label string
}

// NavData groups all elements needed to render the navbar.
type NavData struct {
	HomeHref     string        // e.g. "/" on internal pages, empty on home
	NavDropdowns []NavDropdown // dropdowns de navegación (Mundos, Secciones)
	Links        []NavLink     // links planos (Valoraciones, Suscripciones)
	UserDropdown *NavDropdown  // dropdown de usuario, siempre al final
}

// Post represents a user publication.
type Post struct {
	User      string
	Title     string
	Text      string
	Location  string
	MediaPath string
	Date      string
}

// HomePlanetItem contains world data used in home planet cards.
type HomePlanetItem struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	IsReverse   bool
}

// ReviewItem represents a review shown on home.
type ReviewItem struct {
	Stars  string // e.g. "★★★★★"
	Text   string
	Author string
}

// PlanItem represents a subscription plan on home.
type PlanItem struct {
	ID       string // HTML id, empty when not used
	Name     string
	Features []string
}

// HomePageData contains data for /.
type HomePageData struct {
	Nav      NavData
	LoggedIn bool
	Worlds   []HomePlanetItem
	Reviews  []ReviewItem
	Plans    []PlanItem
}

// WorldPageData contains data for /mundos/{slug}.
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

// SectionPageData contains data for /mundos/{slug}/{section}.
type SectionPageData struct {
	World      content.World
	Section    content.WorldSection
	Nav        NavData
	Questions  []string
	MyPosts    []Post
	OtherPosts []Post
}

// RegisterPageData contains data for /registro.
type RegisterPageData struct {
	Nav   NavData
	Error string
	Name  string
	Email string
}

// LoginPageData contains data for /login.
type LoginPageData struct {
	Nav   NavData
	Error string
	Email string
}

// ContactPageData contains data for /contacto.
type ContactPageData struct {
	Nav     NavData
	Success bool
	Error   string
}
