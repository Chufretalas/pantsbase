package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Chufretalas/pantsbase/db"
	m "github.com/Chufretalas/pantsbase/models"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "index", db.GetAllTableNames())
}

func NewTable(w http.ResponseWriter, r *http.Request) {
	tableName := r.FormValue("name")

	var columnIndexes []string
	if r.FormValue("column_indexes") == "" {
		columnIndexes = make([]string, 0)
	} else {
		columnIndexes = strings.Split(r.FormValue("column_indexes"), " ")
	}

	columns := make([]m.Column, 0, len(columnIndexes))
	for _, index := range columnIndexes {
		name := r.FormValue(fmt.Sprintf("n%v", index))
		typeDB := r.FormValue(fmt.Sprintf("t%v", index))
		columns = append(columns, m.Column{Name: name, TypeDB: typeDB})
	}

	queryStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%v\" (\n id INTEGER PRIMARY KEY AUTOINCREMENT", tableName)
	if len(columnIndexes) != 0 {
		queryStr += ",\n"
		for idx, column := range columns {
			if idx == len(columns)-1 {
				queryStr += fmt.Sprintf(" \"%v\" %v\n", column.Name, column.TypeDB)
			} else {
				queryStr += fmt.Sprintf(" \"%v\" %v,\n", column.Name, column.TypeDB)
			}
		}
	}
	queryStr += ");"

	_, err := db.DB.Exec(queryStr)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func TableView(w http.ResponseWriter, r *http.Request) {
	colsSchemas, err := db.GetSchema(r.URL.Query().Get("name"))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(colsSchemas)
	}

	temp.ExecuteTemplate(w, "table_view", nil)
}
