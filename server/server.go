package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// Page structure
type Page struct {
	Title string
	Body  []byte
}

var stat = map[string]int{"view": 0, "edit": 0, "save": 0}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	title := r.URL.Path[len("/view/"):]

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	stat["view"]++

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	stat["edit"]++

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	title := r.URL.Path[len("/save/"):]

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat["save"]++

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func statHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stat)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// InitServer here
func InitServer() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/view/{page:[a-zA-Z0-9]+}", viewHandler).Methods(http.MethodGet)
	r.HandleFunc("/edit/{page:[a-zA-Z0-9]+}", editHandler).Methods(http.MethodGet)
	r.HandleFunc("/save/{page:[a-zA-Z0-9]+}", saveHandler).Methods(http.MethodPost)
	r.HandleFunc("/stat", statHandler).Methods(http.MethodGet)

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}
