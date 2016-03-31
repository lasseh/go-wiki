package main

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  template.HTML
}

func loadPage(page string) (*Page, error) {
	filename := "git/" + page + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: page, Body: template.HTML(blackfriday.MarkdownCommon([]byte(body)))}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if len(page) == 0 {
		http.Redirect(w, r, "/Index", 301)
		return
	}
	p, err := loadPage(page)
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}
	renderTemplate(w, "view", p)
}

var templates = template.Must(template.ParseFiles("view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":1337", nil)
}
