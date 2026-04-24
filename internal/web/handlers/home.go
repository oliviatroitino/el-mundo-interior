package handlers

import (
	"el-mundo-interior/internal/content"
	"net/http"
)

func Home(sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ordered := content.OrderedWorlds()
		items := make([]HomePlanetItem, len(ordered))
		for i, world := range ordered {
			items[i] = HomePlanetItem{
				Slug:        world.Slug,
				Title:       world.Title,
				Description: world.Description,
				Icon:        world.Icon,
				IsReverse:   i%2 != 0,
			}
		}

		_, userName, loggedIn := sessions.GetUser(r)
		nav := NavData{
			NavDropdowns: []NavDropdown{buildWorldDropdown("")},
			Links: []NavLink{
				{Href: "#other-users", Label: "Valoraciones"},
				{Href: "#subscriptions", Label: "Suscripciones"},
			},
		}
		if loggedIn {
			ud := buildUserDropdown(userName)
			nav.UserDropdown = &ud
		}
		data := HomePageData{
			LoggedIn: loggedIn,
			Nav:      nav,
			Worlds: items,
			Reviews: []ReviewItem{
				{
					Stars:  "★★★★★",
					Text:   "Introspecta me ha ayudado a organizar mis pensamientos y emociones. Es como tener un espacio personal para reflexionar y crecer.",
					Author: "LunaMente",
				},
				{
					Stars:  "★★★★★",
					Text:   "Me encanta la variedad de mundos que ofrece Introspecta. Cada uno me permite explorar diferentes aspectos de mi vida.",
					Author: "Creativa22",
				},
				{
					Stars:  "★★★★★",
					Text:   "La sección de suscripciones es genial, me motiva a mantenerme constante en mi crecimiento personal.",
					Author: "Reflexionario",
				},
			},
			Plans: []PlanItem{
				{
					ID:   "free-plan",
					Name: "Gratis",
					Features: []string{
						"Acceso a todos los mundos",
						"Crear hasta 3 entradas por día",
						"Guardar reflexiones y notas personales",
						"Acceso desde móvil y ordenador",
					},
				},
				{
					Name: "Plus",
					Features: []string{
						"Todo lo incluido en Gratis",
						"Entradas ilimitadas en todos los mundos",
						"Estadísticas de hábitos y progreso personal",
						"Recordatorios diarios y semanales",
						"Exportar tus diarios y reflexiones",
					},
				},
				{
					Name: "Pro",
					Features: []string{
						"Todo lo incluido en Plus",
						"Panel completo de crecimiento personal",
						"Análisis de progreso por mundos",
						"Acceso anticipado a nuevas herramientas",
					},
				},
			},
		}

		render(w, "templates/pages/home.tmpl", data)
	}
}
