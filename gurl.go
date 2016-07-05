package main

import (
	"database/sql"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/dchest/uniuri"
	_ "github.com/go-sql-driver/mysql"
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

var templates = template.Must(template.ParseFiles("index.html"))

func newName() string {
	id := uniuri.NewLen(LENGTH)
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	query, err := db.Query("select id from url where id=?", id)
	if err != sql.ErrNoRows {
		for query.Next() {
			newName()
		}
	}

	return id
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	id := newName()
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	stm, err := db.Prepare("insert into url values(?, ?)")
	shorten := ADDRESS + "/s/" + id
	_, err = stm.Exec(id, html.EscapeString(url))
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p><b>URL</b>: <a href='"+shorten+"'>"+shorten+"</a></p>")

}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["urlid"]
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var url string
	err = db.QueryRow("select url from url where id=?", html.EscapeString(id)).Scan(&url)
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
	router.HandleFunc("/s/{urlid}", urlHandler)
	router.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Fatal(err)
	}

}
