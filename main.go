package main

import (
	"database/sql"
	"fmt"
	"html/template"

	"./routes"
	"./session"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	_ "github.com/go-sql-driver/mysql"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :80")

	//основные настройки к базе
	dsn := "root:1qazXSW@@tcp(localhost:3306)/films?"
	//указывем кодировку
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)

	db.SetMaxOpenConns(10) //количество подключений к бд

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Map(db)

	m.Use(session.Middleware)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", routes.IndexHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Get("/logout", routes.LogoutHandler)
	m.Get("/contacts", routes.ContactsHandler)
	m.Get("/news", routes.NewsHandler)
	//m.Post("/empty", routes.EmptyHandler)
	m.Get("/theend", routes.TheEnd)
	m.Get("/search", routes.SearchFilmsHandler)
	m.Get("/write", routes.WriteHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/view/:id", routes.ViewHandler)
	m.Get("/delete/:id", routes.DeleteHandler)

	m.Post("/SavePost", routes.SavePostHandler)
	m.Post("/simplesearch", routes.SimpleSearchFilmsHandler)
	m.Post("/complexsearch", routes.ComplexSearchFilmsHandler)
	m.Post("/search", routes.SearchFilmsHandler)
	m.Post("/review", routes.ReviewHandler)
	m.Post("/gethtml", routes.GetHtmlHandler)
	m.Post("/contacts", routes.ContactsHandler)
	m.Post("/news", routes.NewsHandler)
	//m.Post("/empty", routes.EmptyHandler)
	m.Post("/theend", routes.TheEnd)
	m.Post("/login", routes.PostLoginHandler)
	m.NotFound(routes.PageNotFound)
	m.RunOnAddr(":9000")
}
