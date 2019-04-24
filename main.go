package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pop/herokuenv"
	"pop/suggestion"

	"github.com/gorilla/mux"
)

var tmpl *template.Template
var names map[string]int
var submitPage []byte

func init() {
	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	submitPageFile, openErr := os.Open("templates/submit.html")
	if openErr != nil {
		log.Fatal(openErr)
	}
	defer submitPageFile.Close()
	var err error
	submitPage, err = ioutil.ReadAll(submitPageFile)
	if err != nil {
		log.Fatal(err)
	}

	names = map[string]int{}
}

func main() {
	router := newRouter()

	serveURI := "0.0.0.0:" + herokuenv.Port
	log.Printf("Serving at %s\n", serveURI)
	http.ListenAndServe(serveURI, router)
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("GET")

	staticFileDir := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDir))

	router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	router.HandleFunc("/submit", submitPageHandler).Methods("GET")
	router.HandleFunc("/submit", submitHandler).Methods("POST")

	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	go log.Println("Index page hit")
	suggestionList := suggestion.NewSuggestions(&names)
	tmpl.Execute(w, suggestionList)
}

func submitPageHandler(w http.ResponseWriter, r *http.Request) {
	go log.Println("Submit page hit")

	w.Write(submitPage)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	go log.Println("Submit endpoint hit")

	err := r.ParseForm()
	if err != nil {
		log.Println("Error retrieving form value; returned internal server error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")

	names[name]++

	http.Redirect(w, r, "/", http.StatusFound)
}
