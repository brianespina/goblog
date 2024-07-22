package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/mux"
)

type Post struct {
	Slug string
	Body template.HTML
}
type Posts struct {
	Slugs []string
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

func loadPosts(dir string) (*Posts, error) {
	posts, _ := os.ReadDir(dir)
	var names []string
	for _, post := range posts {
		name := post.Name()
		slug, _, _ := strings.Cut(name, ".")
		names = append(names, slug)
	}
	return &Posts{Slugs: names}, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := loadPosts("blog")
	if err != nil {
		throwInternalServerError(w, err)
		return
	}
	t, err := template.ParseFiles(templatePath+"layout.html", templatePath+"home.html")
	if err != nil {
		throwInternalServerError(w, err)
		return
	}
	err = t.Execute(w, posts)
	if err != nil {
		throwInternalServerError(w, err)
		return
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := loadPost(vars["slug"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles(templatePath+"layout.html", templatePath+"page.html")
	if err != nil {
		log.Fatal("error in parsing files")
		return
	}
	err = t.Execute(w, post)
	if err != nil {
		throwInternalServerError(w, err)
		return
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/images/favicon.ico")
}
func throwInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/{slug}", viewHandler)
	http.ListenAndServe(":8080", r)
}
