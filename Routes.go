package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"html/template"
)

type Page struct {
	Title string
	Body []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	title := r.URL.Path[len("/view/"):]	
	
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return 
	}

	renderTemplate(w, "view", p)
}

func InitRoutes() *mux.Router {
	r:= mux.NewRouter()
	
	r.HandleFunc("/view/{page:[a-zA-Z0-9]+}", viewHandler).Methods(http.MethodGet)
	// r.HandleFunc("/edit/{page:[a-zA-Z0-9]+}", editHandler).Methods(http.MethodGet)
	// r.HandleFunc("/save/{page:[a-zA-Z0-9]+}", saveHandler).Methods(http.MethodPost)

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}