package documents

type PostDocument struct {
	Id              string `bson:"_id,omitempty"`
	Title           string
	TitleEng        string
	PreviewHtml     string
	ContentHtml     string
	ContentMarkdown string
	Data            string
	Year            string
	Genre           string
	Directors       string
	Actors          string
	Rating          string
	Image           string
}
