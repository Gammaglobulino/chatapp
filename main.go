package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// handler for HTML go templates
type templateHandler struct {
	once     sync.Once //single instance of the template
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{
		filename: "chat.html",
	})
	http.Handle("room", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
