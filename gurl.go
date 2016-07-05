package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	//PORT that gurl will listen on
	PORT       = ":8000"
	dbNAME     = ""
	dbPASS     = ""
	dbUSERNAME = ""

	DATABASE = dbUSERNAME + ":" + dbPASS + "@/" + dbNAME + "?charset=utf8"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["urlid"]
	delkey := vars["delkey"]

}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["urlid"]
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	http.Redirect(w, r, url, 303)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/new", newHandler)
	router.HandleFunc("/del/{urlid}/{delkey}", delHandler)
	router.HandleFunc("/s/{urlid}", urlHandler)
	router.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Fatal(err)
	}

}
