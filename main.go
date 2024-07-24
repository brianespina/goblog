package main

import (
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
	Slug string
	Body template.HTML
	Meta MetaData
}
type Posts struct {
	Slugs []string
}
type MetaData struct {
	Title string `yaml:"title"`
}

var templatePath = "templates/"
var blogPath = "blog/"

func mdToHtml(md []byte) (post *Post) {
	getMetaData(md, post)
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(post.Body))

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	post.Body = template.HTML(markdown.Render(doc, renderer))
	return
}
func getMetaData(md []byte, post *Post) {
	meta := &MetaData{}
	rest, err := frontmatter.Parse(strings.NewReader(string(md)), meta)
	if err != nil {
		log.Fatal("Error getting meta data")
	}
	post.Meta = *meta
	post.Body = template.HTML(rest)
}
func loadPost(slug string) (*Post, error) {
	filename := blogPath + slug + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//	return &Post{Slug: "title", Body: template.HTML(body)}, nil
	return mdToHtml(body), nil
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
