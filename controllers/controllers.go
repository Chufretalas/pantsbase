package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	m "github.com/Chufretalas/pantsbase/models"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "index", nil)
}

func NewTable(w http.ResponseWriter, r *http.Request) {
	tableName := r.FormValue("name")
	columnIndexes := strings.Split(r.FormValue("column_indexes"), " ")
	fmt.Println(columnIndexes)
	columns := make([]m.Column, 0, len(columnIndexes))
	for _, index := range columnIndexes {
		name := r.FormValue(fmt.Sprintf("n%v", index))
		typeDB := r.FormValue(fmt.Sprintf("t%v", index))
		columns = append(columns, m.Column{Name: name, TypeDB: typeDB})
	}
	fmt.Println(tableName)
	fmt.Println(columns)
}
