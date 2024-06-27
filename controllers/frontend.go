package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Chufretalas/pantsbase/db"
)

var Temps *template.Template

func Index(w http.ResponseWriter, r *http.Request) {
	Temps.ExecuteTemplate(w, "index", db.GetAllTableNames())
}

func TableView(w http.ResponseWriter, r *http.Request) {
	colsSchemas, err := db.GetSchema(r.URL.Query().Get("name"))

	if err != nil {
		fmt.Println(err)
	}

	encoded, err := json.Marshal(colsSchemas)

	if err != nil {
		fmt.Println(err)
	}

	Temps.ExecuteTemplate(w, "table_view", map[string]any{
		"Schema":     colsSchemas,
		"JSONSchema": string(encoded),
		"TableName":  r.URL.Query().Get("name"),
	})
}
