package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pop/herokuenv"
	"pop/shutdown"
	"pop/store"
	"pop/suggestion"
	"sync"

	"github.com/gorilla/mux"
)

var tmpl *template.Template
var names map[string]int
var namesLock sync.RWMutex
var wg sync.WaitGroup
var submitPage []byte
var saveEnv bool

func init() {
	log.SetOutput(os.Stderr)

	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	submitPageFile, openErr := os.Open("templates/submit.html")
	defer submitPageFile.Close()
	if openErr != nil {
		log.Fatal(openErr)
	}
	var err error
	submitPage, err = ioutil.ReadAll(submitPageFile)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if !herokuenv.DatabaseURIExists() {
		log.Fatalln("$DATABASE_URL must be set")
		return
	}

	db, err := sql.Open("postgres", herokuenv.DatabaseURI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Ping succeded")

	store.SetStore(store.NewDataStore(db))
	defer store.DataStore.Close()

	if err := store.DataStore.InitTable(); err != nil {
		panic(err)
	}
	log.Println("Table initialization succeded")

	s, err := store.DataStore.Load()
	if err != nil {
		log.Fatalf("Loading from the database failed: %v\n", err)
	}
	log.Println("Loading from the database succeded")
	namesLock.Lock()
	names = s
	namesLock.Unlock()

	router := newRouter()

	serveURI := "0.0.0.0:" + herokuenv.Port
	server := http.Server{
		Addr:    serveURI,
		Handler: router,
	}
	wg.Add(1)
	go shutdown.WaitShutdown(&server, &wg)
	log.Printf("Listening at %s\n", serveURI)
	server.ListenAndServe()

	log.Println("Waiting for goroutines to finish")
	wg.Wait()
	log.Println("Wait group empty; exiting")
	// TODO: extra save needed?
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
	log.Println("Index page hit")

	namesLock.RLock()
	suggestionList := suggestion.NewSuggestions(&names)
	namesLock.RUnlock()

	tmpl.Execute(w, suggestionList)
}

func submitPageHandler(w http.ResponseWriter, r *http.Request) {
	go log.Println("Submit page hit")

	w.Write(submitPage)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Submit endpoint hit")

	err := r.ParseForm()
	if err != nil {
		log.Println("Error retrieving form value; returned internal server error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")

	namesLock.Lock()
	names[name]++
	namesLock.Unlock()

	wg.Add(1)
	go store.ConcurrentSave(names, store.DataStore, &wg)

	http.Redirect(w, r, "/", http.StatusFound)
}
