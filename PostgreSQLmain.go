package main

import (
	"fmt"
	"log"

	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Cafe struct {
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
}
