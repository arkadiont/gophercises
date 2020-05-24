package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// Option struct
type Option struct {
	Text string `json:"text"`
	Link string `json:"arc"`
}

// Chapter struct
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Book map
type Book map[string]Chapter

// GetBook from a reader
func GetBook(r io.Reader) (b Book, err error) {
	err = json.NewDecoder(r).Decode(&b)
	return
}

type handler struct {
	b Book
	t *template.Template
	f func(*http.Request) string
}

// HandlerOption for configure handler
type HandlerOption func(*handler)

// WithTemplate add template
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithParseFunc add parse func
func WithParseFunc(f func(*http.Request) string) HandlerOption {
	return func(h *handler) {
		h.f = f
	}
}

// NewHandler create new handler
func NewHandler(b Book, ops ...HandlerOption) http.Handler {
	h := handler{
		b: b,
		t: template.Must(template.New("").Parse(defaultHTML)),
		f: defaultParse,
	}
	for _, op := range ops {
		op(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.f(r)
	if chapter, ok := h.b[key]; ok {
		if err := h.t.Execute(w, chapter); err != nil {
			log.Printf("err populating template : %s", err)
			http.Error(w, "Ups, something wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

func defaultParse(r *http.Request) string {
	p := strings.TrimSpace(r.URL.Path)
	if p == "" || p == "/" {
		p = "/intro"
	}
	return p[1:]
}

var defaultHTML = `
<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
    <ul>
        {{range .Options}}
            <li> <a href="/{{.Link}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
`
