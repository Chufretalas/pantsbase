package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

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

	allIds := make([]string, 0, len(colsSchemas))
	for _, schema := range colsSchemas {
		allIds = append(allIds, schema.Id)
	}

	Temps.ExecuteTemplate(w, "table_view", map[string]any{
		"Schema":         colsSchemas,
		"HiddenInputIds": strings.Join(allIds, " "),
		"TableName":      r.URL.Query().Get("name"),
	})
}
