package routes

import (
	"database/sql"
	"net/http"

	"fmt"
	"regexp"
	"strings"

	"../models"
	"../session"
	"../utils"
	"github.com/martini-contrib/render"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
)

func WriteHandler(rnd render.Render, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	} else {
		model := models.EditPostModel{}
		model.IsAuthorized = s.IsAuthorized
		model.Post = models.Post{}
		rnd.HTML(200, "write", model)
	}
}

func EditHandler(s *session.Session, rnd render.Render, r *http.Request, params martini.Params, db *sql.DB) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}

	idfilms := params["id"]

	post := &models.Post{}
	//QueryRow -работа с единичными записями, сам закрыает коннект
	row := db.QueryRow("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE idfilms = ?", idfilms)

	err := row.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)

	if err != nil {
		rnd.HTML(500, "500", nil)
	} else {

		model := models.EditPostModel{}
		model.IsAuthorized = s.IsAuthorized
		model.Post = *post
		rnd.HTML(200, "write", model)
	}
}

func ViewHandler(s *session.Session, rnd render.Render, r *http.Request, params martini.Params, db *sql.DB) {

	idfilms := params["id"]

	post := &models.Post{}
	//QueryRow -работа с единичными записями, сам закрыает коннект
	row := db.QueryRow("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE idfilms = ?", idfilms)

	err := row.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)

	if err != nil {
		rnd.HTML(404, "404", nil)
	} else {

		model := models.ViewPostModel{}
		model.IsAuthorized = s.IsAuthorized
		model.Post = post
		rnd.HTML(200, "view", model)
	}
}

func SavePostHandler(s *session.Session, rnd render.Render, r *http.Request, db *sql.DB) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}
	idfilms := r.FormValue("id")
	title := r.FormValue("title")
	titleeng := r.FormValue("titleeng")
	previewHtml := r.FormValue("preview")
	contentHtml := r.FormValue("content")
	data := r.FormValue("data")
	year := r.FormValue("year")
	genre := r.FormValue("genre")
	directors := r.FormValue("directors")
	actors := r.FormValue("actors")
	rating := r.FormValue("rating")
	image := r.FormValue("image")
	description := r.FormValue("description")
	keywords := r.FormValue("keywords")
	if idfilms != "" {
		_, err := db.Exec("UPDATE review SET title=?, titleeng=?, previewHtml=?, contentHtml=?, data=?, year=?, genre=?, directors=?, actors=?, rating=?, image=?, description=?, keywords=? WHERE idfilms=?", title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords, idfilms)
		if err != nil {
			panic(err)
			//rnd.HTML(200, "500", nil)
		}
	} else {
		idfilms = utils.GenerateId()
		fmt.Println(idfilms)
		_, err := db.Exec("INSERT INTO review(idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?);", idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords)
		if err != nil {
			panic(err)
			//rnd.HTML(200, "500", nil)
		}
	}

	rnd.Redirect("/")
}

func SimpleSearchFilmsHandler(s *session.Session, rnd render.Render, r *http.Request, db *sql.DB) {
	simplesearch := r.FormValue("simplesearch")
	newsearch := TextZamena(simplesearch)
	if newsearch == "" {
		rnd.HTML(200, "empty", nil)
	} else {
		posts := []*models.Post{}
		rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (title,titleeng) AGAINST (\"" + newsearch + "*\" IN BOOLEAN MODE)")
		if err != nil {
			//panic(err)
			rnd.HTML(500, "500", nil)
		} else {

			for rows.Next() {
				post := &models.Post{}
				err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
				posts = append(posts, post)

			}
			rows.Close()
			if len(posts) == 0 {
				rnd.HTML(200, "empty", nil)
			} else {
				model := models.PostListModel{}
				model.IsAuthorized = s.IsAuthorized
				model.Posts = posts

				rnd.HTML(200, "index", model)
			}
		}
	}
}

func ComplexSearchFilmsHandler(s *session.Session, rnd render.Render, r *http.Request, db *sql.DB) {
	namefilms := r.FormValue("namefilm")
	actorss := r.FormValue("actor")
	directorss := r.FormValue("director")
	years := r.FormValue("year")
	genres := r.FormValue("genre")
	//удаляем все лишнее из полей
	namefilm := TextZamena(namefilms)
	actors := TextZamena(actorss)
	directors := TextZamena(directorss)
	year := TextZamena(years)
	genre := TextZamena(genres)
	//создаем 2 среза
	var name []string
	var value []string

	i := 0
	if namefilm == "" && actors == "" && directors == "" && year == "" && genre == "" {
		rnd.HTML(200, "empty", nil)
	} else {
		if namefilm != "" {
			i++
			name = append(name, "title,titleeng")
			value = append(value, namefilm)
		}
		if actors != "" {
			i++
			name = append(name, "actors")
			value = append(value, actors)
		}
		if directors != "" {
			i++
			name = append(name, "directors")
			value = append(value, directors)
		}
		if year != "" {
			i++
			name = append(name, "year")
			value = append(value, year)
		}
		if genre != "" {
			i++
			name = append(name, "genre")
			value = append(value, genre)
		}

		if i == 1 {
			names := name[0]
			values := value[0]
			posts := []*models.Post{}
			rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (" + names + ") AGAINST (\"" + values + "*\" IN BOOLEAN MODE)")
			if err != nil {
				//panic(err)
				rnd.HTML(500, "500", nil)
			} else {

				for rows.Next() {
					post := &models.Post{}
					err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
					posts = append(posts, post)
				}
				rows.Close()
				if len(posts) == 0 {
					rnd.HTML(200, "empty", nil)
				} else {
					model := models.PostListModel{}
					model.IsAuthorized = s.IsAuthorized
					model.Posts = posts

					rnd.HTML(200, "index", model)
				}
			}
		}

		if i == 2 {
			names := name[0]
			names2 := name[1]
			values := value[0]
			values2 := value[1]
			posts := []*models.Post{}
			rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (" + names + ") AGAINST (\"" + values + "*\" IN BOOLEAN MODE) AND MATCH (" + names2 + ") AGAINST (\"" + values2 + "*\" IN BOOLEAN MODE)")
			if err != nil {
				//panic(err)
				rnd.HTML(500, "500", nil)
			} else {

				for rows.Next() {
					post := &models.Post{}

					err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
					posts = append(posts, post)

				}
				rows.Close()
				if len(posts) == 0 {
					rnd.HTML(200, "empty", nil)
				} else {
					model := models.PostListModel{}
					model.IsAuthorized = s.IsAuthorized
					model.Posts = posts

					rnd.HTML(200, "index", model)
				}
			}
		}

		if i == 3 {
			names := name[0]
			names2 := name[1]
			names3 := name[2]

			values := value[0]
			values2 := value[1]
			values3 := value[2]

			posts := []*models.Post{}
			rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (" + names + ") AGAINST (\"" + values + "*\" IN BOOLEAN MODE) AND MATCH (" + names2 + ") AGAINST (\"" + values2 + "*\" IN BOOLEAN MODE) AND MATCH (" + names3 + ") AGAINST (\"" + values3 + "*\" IN BOOLEAN MODE)")

			if err != nil {
				//panic(err)
				rnd.HTML(500, "500", nil)
			} else {

				for rows.Next() {
					post := &models.Post{}

					err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
					posts = append(posts, post)

				}
				rows.Close()
				if len(posts) == 0 {
					rnd.HTML(200, "empty", nil)
				} else {
					model := models.PostListModel{}
					model.IsAuthorized = s.IsAuthorized
					model.Posts = posts

					rnd.HTML(200, "index", model)
				}
			}
		}

		if i == 4 {
			names := name[0]
			names2 := name[1]
			names3 := name[2]
			names4 := name[3]

			values := value[0]
			values2 := value[1]
			values3 := value[2]
			values4 := value[3]

			posts := []*models.Post{}
			rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (" + names + ") AGAINST (\"" + values + "*\" IN BOOLEAN MODE) AND MATCH (" + names2 + ") AGAINST (\"" + values2 + "*\" IN BOOLEAN MODE) AND MATCH (" + names3 + ") AGAINST (\"" + values3 + "*\" IN BOOLEAN MODE) AND MATCH (" + names4 + ") AGAINST (\"" + values4 + "*\" IN BOOLEAN MODE)")

			if err != nil {
				//panic(err)
				rnd.HTML(500, "500", nil)
			} else {

				for rows.Next() {
					post := &models.Post{}

					err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
					posts = append(posts, post)

				}
				rows.Close()
				if len(posts) == 0 {
					rnd.HTML(200, "empty", nil)
				} else {
					model := models.PostListModel{}
					model.IsAuthorized = s.IsAuthorized
					model.Posts = posts

					rnd.HTML(200, "index", model)
				}
			}
		}

		if i == 5 {
			names := name[0]
			names2 := name[1]
			names3 := name[2]
			names4 := name[3]
			names5 := name[4]

			values := value[0]
			values2 := value[1]
			values3 := value[2]
			values4 := value[3]
			values5 := value[4]

			posts := []*models.Post{}
			rows, err := db.Query("SELECT idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords FROM review WHERE MATCH (" + names + ") AGAINST (\"" + values + "*\" IN BOOLEAN MODE) AND MATCH (" + names2 + ") AGAINST (\"" + values2 + "*\" IN BOOLEAN MODE) AND MATCH (" + names3 + ") AGAINST (\"" + values3 + "*\" IN BOOLEAN MODE) AND MATCH (" + names4 + ") AGAINST (\"" + values4 + "*\" IN BOOLEAN MODE) AND MATCH (" + names5 + ") AGAINST (\"" + values5 + "*\" IN BOOLEAN MODE)")

			if err != nil {
				//panic(err)
				rnd.HTML(500, "500", nil)
			} else {

				for rows.Next() {
					post := &models.Post{}

					err = rows.Scan(&post.Idfilms, &post.Title, &post.TitleEng, &post.PreviewHtml, &post.ContentHtml, &post.Data, &post.Year, &post.Genre, &post.Directors, &post.Actors, &post.Rating, &post.Image, &post.Description, &post.Keywords)
					posts = append(posts, post)

				}
				rows.Close()
				if len(posts) == 0 {
					rnd.HTML(200, "empty", nil)
				} else {
					model := models.PostListModel{}
					model.IsAuthorized = s.IsAuthorized
					model.Posts = posts

					rnd.HTML(200, "index", model)
				}
			}
		}
	}
}

func DeleteHandler(s *session.Session, rnd render.Render, r *http.Request, params martini.Params, db *sql.DB) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	} else {
		idfilms := params["id"]
		if idfilms == "" {
			rnd.Redirect("/")
			return
		}
		_, err := db.Exec("DELETE FROM review WHERE idfilms=?", idfilms)
		if err != nil {
			panic(err)
		}
		rnd.Redirect("/")
	}

}

func GetHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}
func ContactsHandler(rnd render.Render) {
	rnd.HTML(200, "contacts", nil)
}
func NewsHandler(rnd render.Render) {
	rnd.HTML(200, "news", nil)
}

func PageNotFound(rnd render.Render) {
	rnd.HTML(404, "404", nil)
}
func TheEnd(rnd render.Render) {
	rnd.HTML(500, "500", nil)
}

/*func EmptyHandler(rnd render.Render) {
	rnd.HTML(200, "empty", nil)
}*/
func SearchFilmsHandler(rnd render.Render) {
	rnd.HTML(200, "search", nil)
}
func СutText(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

func TextZamena(stroka string) string {
	//регулярка уберет все знаки пунктуации
	var re = regexp.MustCompile(`[[:punct:]]`)
	// замена всех пробелов звездочками
	var replacer = strings.NewReplacer(" ", "*", "\f", "*")
	//обрезаем строку до 60 символов
	str := СutText(stroka, 60)
	//Убираем все знаки пунктуации и т.д.
	str1 := re.ReplaceAllString(str, "")
	//Убираем все лишние пробелы, оставляем только по 1 пробел. Все заносится в массив []
	str2 := strings.Fields(str1)
	//преобразуем из []string в string
	str3 := []string(str2)
	str4 := strings.Join(str3, " ")
	//все пробелы заменяем *
	str5 := replacer.Replace(str4)
	return str5
}
