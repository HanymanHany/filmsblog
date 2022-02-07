package models

type Post struct {
	Idfilms     string
	Title       string
	TitleEng    string
	PreviewHtml string
	ContentHtml string
	Data        string
	Year        string
	Genre       string
	Directors   string
	Actors      string
	Rating      string
	Image       string
	Description string
	Keywords    string
}

func NewPost(idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords string) *Post {
	return &Post{idfilms, title, titleeng, previewHtml, contentHtml, data, year, genre, directors, actors, rating, image, description, keywords}

}
