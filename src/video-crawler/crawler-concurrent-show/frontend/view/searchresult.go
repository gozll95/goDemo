package view

import (
	"html/template"
	"io"

	"video-crawler/crawler-concurrent-show/frontend/model"
)

type SearchReusltView struct {
	template *template.Template
}

func CreateSearchResultView(filename string) SearchReusltView {
	return SearchReusltView{
		template: template.Must(
			template.ParseFiles(filename)),
	}
}

func (s SearchReusltView) Render(w io.Writer, data model.SearchResult) error {
	return s.template.Execute(w, data)
}
