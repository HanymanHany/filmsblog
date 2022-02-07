package routes

import (
	"fmt"
	"net/http"

	"../session"

	"github.com/martini-contrib/render"
)

func GetLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func PostLoginHandler(rnd render.Render, r *http.Request, s *session.Session) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username)
	fmt.Println(password)
	if username == "admin" && password == "IgO6%Vt*Q{" {
		s.Username = username
		s.IsAuthorized = true
		rnd.Redirect("/")
	} else {
		rnd.Redirect("/")
	}
	//rnd.Redirect("/")
}

func LogoutHandler(rnd render.Render, r *http.Request, s *session.Session) {
	s.Username = ""
	s.IsAuthorized = false

	rnd.Redirect("/")
}
