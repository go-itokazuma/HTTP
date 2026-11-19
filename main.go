package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Cafes struct {
	ID      int
	Name    string
	Address string
	Open    string
}

func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				override := strings.ToUpper((r.Form.Get("_method")))
				if override == http.MethodPut || override == http.MethodDelete {
					r.Method = override
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func loadTemplates() map[string]*template.Template {
	files := map[string]string{
		"List":   "templates/list/list.html.tmpl",
		"Detail": "templates/detail/detail.html.tmpl",
		"New":    "templates/new/new.html.tmpl",
		"Edit":   "templates/edit/edit.html.tmpl",
	}

	templates := make(map[string]*template.Template)
	for name, path := range files {
		templates[name] = template.Must(template.ParseFiles(path))
	}
	return templates
}

func getTemplate(templates map[string]*template.Template, name string) (*template.Template, error) {
	tmpl, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}

func handleNewCafe(templates map[string]*template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl, err := getTemplate(templates, "New")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "New", nil)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	}
}

func handleCreateCafe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		address := r.FormValue("address")
		open := r.FormValue("open")

		_, err := db.Exec("INSERT INTO cafes (name, address, open) VALUES ($1, $2, $3)", name, address, open)
		if err != nil {
			http.Error(w, "Database insert error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/cafes", http.StatusSeeOther)
	}
}

func handleEditCafe(templates map[string]*template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl, err := getTemplate(templates, "Edit")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "Edit", c)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	}
}

func handleListCafes(templates map[string]*template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM cafes ORDER BY id")
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close() // 必ずCloseする

		var cafeList []Cafes
		for rows.Next() {
			var id int
			var name string

			err := rows.Scan(&id, &name)
			if err != nil {
				http.Error(w, "Database query error", http.StatusInternalServerError)
				return
			}
			cafeList = append(cafeList, Cafes{ID: id, Name: name})
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl, err := getTemplate(templates, "List")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "List", cafeList)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	}
}

func handleDetailCafe(templates map[string]*template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl, err := getTemplate(templates, "Detail")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "Detail", c)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	}
}

func handleUpdateCafe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		address := r.FormValue("address")
		open := r.FormValue("open")

		_, err = db.Exec("UPDATE cafes SET name = $1, address = $2, open = $3 WHERE id = $4", name, address, open, id)
		if err != nil {
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/cafes/%d", id), http.StatusSeeOther)
	}
}

func handleDeleteCafe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM cafes WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Database delete error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/cafes", http.StatusSeeOther)
	}
}

func main() {
	mux := http.NewServeMux()

	dsn := "user=kazuma.ito dbname=cafe sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	templates := loadTemplates()

	mux.HandleFunc("GET /cafes", handleListCafes(templates, db))
	mux.HandleFunc("GET /cafes/{id}", handleDetailCafe(templates, db))
	mux.HandleFunc("GET /cafes/new", handleNewCafe(templates))
	mux.HandleFunc("POST /cafes", handleCreateCafe(db))
	mux.HandleFunc("GET /cafes/{id}/edit", handleEditCafe(templates, db))
	mux.HandleFunc("PUT /cafes/{id}", handleUpdateCafe(db))
	mux.HandleFunc("DELETE /cafes/{id}", handleDeleteCafe(db))

	log.Fatal(http.ListenAndServe(":8080", methodOverrideMiddleware(mux)))
}
