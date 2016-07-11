// package gurl is a simple url shortening service
package main

import (
	"database/sql"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	// for random string generation
	"github.com/dchest/uniuri"
	// mysql db driver
	_ "github.com/go-sql-driver/mysql"
	// mux for route handling
	"github.com/gorilla/mux"
)

const (
	//PORT that gurl will listen on
	PORT = ":8000"
	// dbNAME database name
	dbNAME = ""
	// dbPASS database password
	dbPASS = ""
	// dbUSERNAME database username
	dbUSERNAME = ""
	// LENGTH url length
	LENGTH = 6
	// ADDRESS url for shortening service
	ADDRESS = "localhost:8000"

	// DATABASE connection string
	DATABASE = dbUSERNAME + ":" + dbPASS + "@/" + dbNAME + "?charset=utf8"
)

// template file
var templates = template.Must(template.ParseFiles("index.html"))

// newName generates a new name that isn't in the db
// returns string
func newName() string {
	// generate a new name
	id := uniuri.NewLen(LENGTH)
	// open db connection
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// check if name exists in the db
	query, err := db.Query("select id from url where id=?", id)
	// if name exists call newName again
	if err != sql.ErrNoRows {
		for query.Next() {
			newName()
		}
	}

	// return id
	return id
}

// newHandler handles saving a new url into the db
func newHandler(w http.ResponseWriter, r *http.Request) {
	// get form value
	url := r.FormValue("url")

	// call newName
	id := newName()

	// open a db connection
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// prepare insert statement
	stm, err := db.Prepare("insert into url values(?, ?, ?)")

	// make shortened url
	shorten := ADDRESS + "/s/" + id

	// Execute query
	_, err = stm.Exec(id, html.EscapeString(url), time.Now().Format("2016-02-01 12:05:12"))
	if err != nil {
		log.Println(err)
	}
	// return data in html type
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<p><b>URL</b>: <a href='"+shorten+"'>"+shorten+"</a></p>")

}

// urlHandler redirects the user to their unshortened url
func urlHandler(w http.ResponseWriter, r *http.Request) {
	// get id from url
	vars := mux.Vars(r)
	id := vars["urlid"]

	// open db econnection
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// prepare url variable
	var url string
	// query database for address from id
	err = db.QueryRow("select url from url where id=?", html.EscapeString(id)).Scan(&url)
	if err != nil {
		log.Println(err)
	}

	// Redirect user to their url
	http.Redirect(w, r, url, 303)

}

// rootHandler generates the index page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Execute index template
	err := templates.Execute(w, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {
	// new mux router
	router := mux.NewRouter()
	router.HandleFunc("/new", newHandler)
	router.HandleFunc("/s/{urlid}", urlHandler)
	router.HandleFunc("/", rootHandler)
	// listen on PORT and serve router
	err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Fatal(err)
	}

}
