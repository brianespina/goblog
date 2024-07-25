package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/mux"
)

type Post struct {
	Slug  string
	Title string
	Body  template.HTML
}
type Posts struct {
	Slugs []string
}
type MetaData struct {
	Title string `yaml:"title"`
}

var templatePath = "templates/"
var blogPath = "blog/"

func mdToHtml(post *Post, body []byte) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.ToHTML(body, p, renderer)
	post.Body = template.HTML(html)
}
func loadPost(slug string) (*Post, error) {
	filename := blogPath + slug + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	meta := &MetaData{}
	rest, err := frontmatter.Parse(bytes.NewReader(body), meta)
	if err != nil {
		return nil, err
	}
	post := &Post{Slug: "title", Title: meta.Title}
	mdToHtml(post, rest)
	return post, nil
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
		throwInternalServerError(err)
		return
	}
	t, err := template.ParseFiles(templatePath+"layout.html", templatePath+"home.html")
	if err != nil {
		throwInternalServerError(err)
		return
	}
	err = t.Execute(w, posts)
	if err != nil {
		throwInternalServerError(err)
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
		throwInternalServerError(err)
		return
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/images/favicon.ico")
}
func throwInternalServerError(err error) {
	log.Fatal(err.Error(), "internal server error")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/{slug}", viewHandler)
	http.ListenAndServe(":8080", r)
}
