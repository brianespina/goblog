//go:generate npm run build
package main

import (
	"embed"
	"espinabrian/mdblog/models"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Posts []models.Post

var templatePath = "templates/"

func loadPosts(dir string) (*Posts, error) {
	fnames, _ := os.ReadDir(dir)
	var posts Posts

	for _, file := range fnames {
		name := file.Name()
		slug, _, _ := strings.Cut(name, ".")
		p := models.NewPost()
		p.Load(slug)
		_, err := os.ReadFile("public/images/" + slug + ".jpg")
		if err == nil {
			p.Image = true
		}
		posts = append(posts, *p)
	}
	return &posts, nil
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
	post := models.Post{}
	err := post.Load(vars["slug"])
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

var static embed.FS

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	// Define the directory to serve static files from
	fs := http.FileServer(http.Dir("public"))

	// Handle static files
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/{slug}", viewHandler)
	http.ListenAndServe(":8080", r)
}
