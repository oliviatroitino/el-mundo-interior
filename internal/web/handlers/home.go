package handlers

import (
	"el-mundo-interior/internal/content"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		Worlds: content.OrderedWorlds(),
	}
	render(w, "templates/pages/home.tmpl", data)
}
