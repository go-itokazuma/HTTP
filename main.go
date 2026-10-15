package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, World!")
		err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
