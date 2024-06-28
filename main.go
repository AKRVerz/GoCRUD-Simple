package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
type Todo struct {
	gorm.Model
	ID       uint
	Title    string
	ImageURL string
	DueDate  string
	Done     bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect DB")
	}
	db.AutoMigrate(&Todo{})

	tmpl := template.Must(template.ParseFiles("template/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			todo := r.FormValue("todo")
			imageURL := r.FormValue("imageURL")
			dueDate := r.FormValue("dueDate")
			db.Create(&Todo{Title: todo, ImageURL: imageURL, DueDate: dueDate, Done: false})
		}

		var todos []Todo
		db.Find(&todos)
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos:     todos,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/done/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/done/")
		var todo Todo
		db.First(&todo, id)
		todo.Done = true
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/delete/")
		db.Delete(&Todo{}, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.ListenAndServe(":8080", nil)
}
