package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rayten225/kinofast/filmObj"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./films.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS films (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		img_url TEXT,
		title TEXT,
		description TEXT,
		time TEXT,
		genre TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func getFilms(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT img_url, title, description, time, genre FROM films")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var films []filmObj.Film
	for rows.Next() {
		var film filmObj.Film
		err := rows.Scan(&film.ImgUrl, &film.Title, &film.Description, &film.Time, &film.Genre)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		films = append(films, film)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}

func main() {
	initDB()
	insertSampleData()

	http.HandleFunc("/films", getFilms)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func insertSampleData() {
	films := []filmObj.Film{
		filmObj.NewFilm("url1", "Title1", "Description1", "120 min", "Genre1"),
		filmObj.NewFilm("url2", "Title2", "Description2", "90 min", "Genre2"),
	}

	for _, film := range films {
		_, err := db.Exec("INSERT INTO films (img_url, title, description, time, genre) VALUES (?, ?, ?, ?, ?)",
			film.ImgUrl, film.Title, film.Description, film.Time, film.Genre)
		if err != nil {
			log.Fatal(err)
		}
	}
}
