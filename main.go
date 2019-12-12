package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var db *sql.DB

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Instantiate the configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("/home/sakhtar/Documents/workspaces/golang/src/github.com/shamimgeek/url-shortner/")
	viper.ReadInConfig()

	// Instantiate the database
	var err error
	dsn := viper.GetString("mysql_user") + ":" + viper.GetString("mysql_password") + "@tcp(" + viper.GetString("mysql_host") + ":3306)/" + viper.GetString("mysql_database") + "?collation=utf8mb4_unicode_ci&parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Instantiate the mux router
	r := mux.NewRouter()
	r.HandleFunc("/links", ShortenHandler).Queries("url", "")
	r.HandleFunc("/{shortid:[a-zA-Z0-9]+}", ShortenedUrlHandler)
	r.HandleFunc("/", CatchAllHandler)

	// Assign mux as the HTTP handler
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

// Shortens a given URL passed through in the request.
// If the URL has already been shortened, returns the existing URL.
// Writes the short URL in plain text to w.
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the url parameter has been sent along (and is not empty)
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Get the short URL out of the config
	if !viper.IsSet("short_url") {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	short_url := viper.GetString("short_url")

	// Check if url already exists in the database
	var shortid string
	err := db.QueryRow("SELECT `shortid` FROM `tiny_urls` WHERE `url` = ?", url).Scan(&shortid)
	if err == nil {
		// The URL already exists! Return the shortened URL.
		w.Write([]byte(short_url + "/" + shortid))
		return
	}

	// generate a shortid and validate it doesn't
	// exist until we find a valid one.
	var exists = true
	for exists == true {
		shortid = generateshortid()
		err, exists = shortidExists(shortid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Insert it into the database
	stmt, err := db.Prepare("INSERT INTO `tiny_urls` (`shortid`, `url`, `date`, `hits`) VALUES (?, ?, NOW(), ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(shortid, url, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(short_url + "/" + shortid))
}

// generateshortid will generate a random shortid to be used as shorten link.
func generateshortid() string {
	// It doesn't exist! Generate a new shortid for it
	// From: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	s := make([]rune, 6)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}

	return string(s)
}

// shortidExists will check whether the shortid already exists in the database
func shortidExists(shortid string) (e error, exists bool) {
	err := db.QueryRow("SELECT EXISTS(SELECT * FROM `tiny_urls` WHERE `shortid` = ?)", shortid).Scan(&exists)
	if err != nil {
		return err, false
	}

	return nil, exists
}

// Handles a requested short URL.
// Redirects with a 301 header if found.
func ShortenedUrlHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Check if a shortid exists
	vars := mux.Vars(r)
	shortid, ok := vars["shortid"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// 2. Check if the shortid exists in the database
	var url string
	err := db.QueryRow("SELECT `url` FROM `tiny_urls` WHERE `shortid` = ?", shortid).Scan(&url)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 3. If the shortid (and thus the URL) exist, update the hit counter
	stmt, err := db.Prepare("UPDATE `tiny_urls` SET `hits` = `hits` + 1 WHERE `shortid` = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(shortid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Finally, redirect the user to the URL
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// Catches all other requests to the short URL domain.
// If a default URL exists in the config redirect to it.
func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the redirect URL out of the config
	if !viper.IsSet("default_url") {
		// The reason for using StatusNotFound here instead of StatusInternalServerError
		// is because this is a catch-all function. You could come here via various
		// ways, so showing a StatusNotFound is friendlier than saying there's an
		// error (i.e. the configuration is missing)
		http.NotFound(w, r)
		return
	}

	// 2. If it exists, redirect the user to it
	http.Redirect(w, r, viper.GetString("default_url"), http.StatusMovedPermanently)
}
