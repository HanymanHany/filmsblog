package routes

import (
	"database/sql"

	"../models"
	"../session"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
)

func IndexHandler(r *http.Request, rnd render.Render, s *session.Session, db *sql.DB) {

	posts := []*models.Post{}
	rows, err := db.Query("SELECT  idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image FROM review ORDER BY id DESC ")
	if err != nil {
		rnd.HTML(500, "500", nil)
	} else {

		for rows.Next() {
			post := &models.Post{}
			err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image)
			posts = append(posts, post)

		}

		rows.Close()

		model := models.PostListModel{}
		model.IsAuthorized = s.IsAuthorized
		model.Posts = posts

		rnd.HTML(200, "index", model)
	}

}

func ReviewHandler(rnd render.Render, r *http.Request, s *session.Session, db *sql.DB) {

	newposts := r.FormValue("new")
	allposts := r.FormValue("all")
	posts := []*models.Post{}

	if newposts == "Новинки 2018" {

		rows, err := db.Query("SELECT  idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE year=2018 ORDER BY id DESC ")
		if err != nil {
			rnd.HTML(500, "500", nil)
		} else {

			for rows.Next() {
				post := &models.Post{}
				err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
				posts = append(posts, post)

			}
			rows.Close()
			model := models.PostListModel{}
			model.IsAuthorized = s.IsAuthorized
			model.Posts = posts

			rnd.HTML(200, "index", model)
		}
	} else if allposts == "Все обзоры" {
		rows, err := db.Query("SELECT  idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review ORDER BY id DESC ")
		if err != nil {
			rnd.HTML(500, "500", nil)
		} else {

			for rows.Next() {
				post := &models.Post{}
				err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
				posts = append(posts, post)

			}

			rows.Close()
			model := models.PostListModel{}
			model.IsAuthorized = s.IsAuthorized
			model.Posts = posts

			rnd.HTML(200, "index", model)
		}

	}

}
