package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	//PORT that gurl will listen on
	PORT = ":8000"
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
