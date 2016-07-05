package main

import (
	"database/sql"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
)

const (
	//PORT that gurl will listen on
	PORT       = ":8000"
	dbNAME     = ""
	dbPASS     = ""
	dbUSERNAME = ""
	LENGTH     = 6
	ADDRESS    = "localhost:8000"

	DATABASE = dbUSERNAME + ":" + dbPASS + "@/" + dbNAME + "?charset=utf8"
)

var templates = template.Must(template.ParseFiles("assets/index.html"))

func newName() string {
	id := uniuri.NewLen(LENGTH)
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	var dbID string
	err = db.QueryRow("select id from urls where id=?", id).Scan(&dbID)
	if err != sql.ErrNoRows {
		newName()
	}

	return id
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	expiry := r.FormValue("expiry")
	id := newName()
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	stm, err := db.Prepare("insert into urls values(?, ?, ?)")
	shorten := ADDRESS + "/s/" + id
	_, err = stm.Exec(id, html.EscapeString(url), html.EscapeString(expiry))
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p><b>URL</b>: <a href='"+shorten+"'>"+shorten+"</a></p>")

}

func delHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["urlid"]
	delkey := vars["delkey"]

}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["urlid"]
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var url, expiry string
	err = db.QueryRow("select url, expiry from url where id=?", html.EscapeString(id)).Scan(&url, &expiry)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, url, 303)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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
