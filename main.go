package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type Cafes struct {
	ID   int
	Name string
}

var Db *sql.DB

func main() {
	dsn := "user=kazuma.ito dbname=cafe sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 実際に接続できるか確認
	if err = db.Ping(); err != nil {
		log.Fatal("DB接続失敗:", err)
	}

	log.Println("DB接続成功")

	rows, err := db.Query("SELECT id, name FROM cafes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // 必ずCloseする

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID:%d Name:%s\n", id, name)
	}

	// rows.Next()後のエラーチェックも忘れずに
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	t, err := template.New("foo").Parse(`{{define "T"}}<html><body>Hello, {{.}}!</body></html>{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, World!")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = t.ExecuteTemplate(w, "T", "World")

	})

	list, err := template.New("foo").Parse(`{{define "List"}}<html><body>カフェリスト:<ul>{{range .}}<li>{{.ID}}: {{.Name}}</li>{{end}}</ul></body></html>{{end}}`)
	http.HandleFunc("GET /cafes", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM cafes")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close() // 必ずCloseする

		var cafeList []Cafes
		for rows.Next() {
			var id int
			var name string

			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			cafeList = append(cafeList, Cafes{ID: id, Name: name})
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = list.ExecuteTemplate(w, "List", cafeList)

	})

	http.HandleFunc("GET /cafes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = t.ExecuteTemplate(w, "T", "World")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
