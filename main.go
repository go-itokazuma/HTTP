package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()

	dsn := "user=kazuma.ito dbname=cafe sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	templates := loadTemplates()

	registerHandlers(mux, templates, db)

	log.Fatal(http.ListenAndServe(":8080", methodOverrideMiddleware(mux)))
}
