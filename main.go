package main

import (
	"../chat/trace"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"html/template"
	"log"
	"net/http"
	"os"
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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Port address for the application")
	flag.Parse() // parse the flags
	gomniauth.SetSecurityKey("Gammaglobulino")
	gomniauth.WithProviders(facebook.New("2440352826280138",
		"523f7cb418c27e6adf7a0cb1169bb4ab",
		"http://localhost:8080/auth/callback/facebook"))
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{
		filename: "chat.html"}))
	http.Handle("/login", &templateHandler{
		filename: "login.html",
	})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
