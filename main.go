package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Cafes struct {
	ID      int
	Name    string
	Address string
	Open    string
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

	rows, err := db.Query("SELECT id, name, address, open FROM cafes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // 必ずCloseする

	for rows.Next() {
		var id int
		var name string
		var address string
		var open string

		err := rows.Scan(&id, &name, &address, &open)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID:%d Name:%s Address:%s Open:%s\n", id, name, address, open)
	}

	// rows.Next()後のエラーチェックも忘れずに
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	t, err := template.New("foo").Parse(`{{define "T"}}<html><body>Hello, {{.}}!</body></html>{{end}}{{define "Cafe"}}<html><body>Cafe ID: {{.}}</body></html>{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, World!")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = t.ExecuteTemplate(w, "T", "World")

	})

	list, err := template.New("foo").Parse(`{{define "List"}}
	<html><body>
		<h1>カフェリスト:</h1>
		<p><a href="/cafes/new">店舗を追加する</a></p>
		<ul>
		{{range .}}
			<li><a href="/cafes/{{.ID}}">{{.ID}}: {{.Name}}</a></li>
		{{end}}
		</ul>
	</body></html>
	{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}
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

	cafe, err := template.New("foo").Parse(`{{define "Cafe"}}
	<html><body>
		{{.ID}}: {{.Name}}<br>
		{{.Address}}<br>
		{{.Open}}
		<p><a href="/cafes">一覧に戻る</a></p>
	</body></html>{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	http.HandleFunc("GET /cafes/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var c Cafes
		err = db.QueryRow("SELECT id, name, address, open FROM cafes WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Address, &c.Open)
		if err != nil {
			http.Error(w, "Cafe not found", http.StatusNotFound)
			return
		}

		//fmt.Fprintf(w, "User ID: %d\n", id)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = cafe.ExecuteTemplate(w, "Cafe", c)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	})

	newTmpl, err := template.New("foo").Parse(`{{define "New"}}
	<html><body>
		<h1>新しいカフェを追加</h1>
		<form action="/cafes" method="post">
			<p>店舗名: <input type="text" name="name"></p>
			<p>住所: <input type="text" name="address"></p>
			<p>営業時間: <input type="text" name="open"></p>
			<p><input type="submit" value="追加"></p>
		</form>
		<p><a href="/cafes">一覧に戻る</a></p>
	</body></html>
	{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	http.HandleFunc("GET /cafes/new", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = newTmpl.ExecuteTemplate(w, "New", nil)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("POST /cafes", func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("name")
		address := r.FormValue("address")
		open := r.FormValue("open")

		_, err := db.Exec("INSERT INTO cafes (name, address, open) VALUES ($1, $2, $3)", name, address, open)
		if err != nil {
			http.Error(w, "Database insert error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/cafes", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
