package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Slug string
	Body template.HTML
}

var templatePath = "templates/"
var blogPath = "blog/"

func mdToHtml(md *[]byte) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(*md)

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	*md = markdown.Render(doc, renderer)
}

func loadPost(slug string) (*Post, error) {
	filename := blogPath + slug + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	mdToHtml(&body)
	return &Post{Slug: "title", Body: template.HTML(body)}, nil
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := os.ReadDir("blog")
	if err != nil {
		log.Fatal(err)
		return
	}
	t, _ := template.ParseFiles(templatePath+"layout.html", templatePath+"home.html")
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	post, err := loadPost("my-first-post")
	if err != nil {
		log.Fatal(err)
		return
	}
	t, err := template.ParseFiles(templatePath+"layout.html", templatePath+"page.html")
	if err != nil {
		log.Fatal("error in parsing files")
		return
	}
	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/blog/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
