package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "gennadiy"
	password = "123"
	dbname   = "appeals"
)

type Appeal struct {
	authorName, authorLocation, authorMail, appealDate, appealText string
}

var connStr string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)

func createListTable() {

	_, createTableErr := db.Exec(`
		CREATE TABLE IF NOT EXISTS appeals_list (
			id SERIAL PRIMARY KEY,
			authorName VARCHAR(255),
			authorLocation VARCHAR(255),
			authorMail VARCHAR(255),
			appealDate VARCHAR(255),
			appealText VARCHAR(255)
		)
	`)
	if createTableErr != nil {
		log.Fatal(createTableErr)

	} else {
		fmt.Println("Table created successfully")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../templates/home.html")
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "../../templates/form.html")

	} else if r.Method == http.MethodPost {

		if pingErr := db.Ping(); pingErr != nil {
			fmt.Println("Error at writing to db")
			log.Fatal(pingErr)
		} else {
			responseData := Appeal{
				authorName:     r.FormValue("authorName"),
				authorLocation: r.FormValue("authorLocation"),
				authorMail:     r.FormValue("authorMail"),
				appealDate:     time.Now().Format("2002-01-02 15:23"),
				appealText:     r.FormValue("appealText"),
			}

			_, insertErr := db.Exec(`
				INSERT INTO appeals_list (authorName, authorLocation, authorMail, appealDate, appealText)
				VALUES ($1, $2, $3, $4, $5)
			`, responseData.authorName,
				responseData.authorLocation,
				responseData.authorMail,
				responseData.appealDate,
				responseData.appealText,
			)

			if insertErr != nil {
				fmt.Println("Error at insert")
				http.ServeFile(w, r, "../../templates/form.html")
				log.Fatal(insertErr)

			} else {
				http.ServeFile(w, r, "../../templates/list.html")
			}

		}

	}

}
func listHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../templates/list.html")
}

func main() {

	createListTable(db)

	//fmt.Println("hello")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/list", listHandler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server at port 8082")
	defer db.Close()
}
