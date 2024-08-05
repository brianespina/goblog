package models

import (
	"bytes"
	"html/template"
	"os"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Title  string
	Body   template.HTML
	Exerpt string
	Slug   string
	Image  bool
}

var blogPath = "blog/"

func NewPost() *Post {
	return &Post{
		Image: false,
	}
}
func (p *Post) Load(slug string) error {
	filename := blogPath + slug + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	meta := &struct {
		Title  string `yaml:"title"`
		Exerpt string `yaml:"exerpt"`
	}{
		Title:  "",
		Exerpt: "",
	}
	rest, err := frontmatter.Parse(bytes.NewReader(body), meta)
	if err != nil {
		return err
	}
	p.Title = meta.Title
	p.Exerpt = meta.Exerpt
	p.Slug = slug
	p.mdToHtml(rest)
	return nil
}

func (p *Post) mdToHtml(body []byte) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	newParse := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.ToHTML(body, newParse, renderer)
	p.Body = template.HTML(html)
}
