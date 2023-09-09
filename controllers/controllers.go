package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "index", nil)
}

func NewTable(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("name"))
}
