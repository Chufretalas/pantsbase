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
	schema, err := db.GetSchema(r.URL.Query().Get("name"))

	if err != nil {
		fmt.Println(err)
	}

	schemaForFrontend := make([]map[string]string, len(schema))

	for idx, col := range schema {
		schemaForFrontend[idx] = make(map[string]string)
		schemaForFrontend[idx]["Name"] = col.Name
		schemaForFrontend[idx]["TypeDB"] = col.TypeDB
		var inputName string
		switch col.TypeDB {
		case "INTEGER":
			inputName = fmt.Sprintf("i%v", idx)
		case "REAL":
			inputName = fmt.Sprintf("r%v", idx)
		case "TEXT":
			inputName = fmt.Sprintf("t%v", idx)

		}
		schemaForFrontend[idx]["InputName"] = inputName
	}

	encoded, err := json.Marshal(schemaForFrontend)

	if err != nil {
		fmt.Println(err)
	}

	Temps.ExecuteTemplate(w, "table_view", map[string]any{
		"Schema":     schemaForFrontend,
		"JSONSchema": string(encoded),
		"TableName":  r.URL.Query().Get("name"),
	})
}
