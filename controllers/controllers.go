package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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

	temp.ExecuteTemplate(w, "table_view", map[string]interface{}{
		"Schema":         colsSchemas,
		"HiddenInputIds": strings.Join(allIds, " "),
		"TableName":      r.URL.Query().Get("name"),
	})
}

func NewRow(w http.ResponseWriter, r *http.Request) {

	failed := false

	tableName := r.FormValue("table_name")
	fieldIds := strings.Split(r.FormValue("new_row_ids"), " ") // it's actually the field names not ids, I'm just dense and I am not changing it now
	fmt.Println(fieldIds)
	values := make([]interface{}, 0, len(fieldIds))

	for _, ids := range fieldIds {
		valueRaw := r.FormValue(ids)
		switch string(ids[0]) {
		case "i":
			convertedValue, err := strconv.Atoi(valueRaw)
			if err != nil {
				fmt.Println(err)
				failed = true
				break
			}
			values = append(values, convertedValue)

		case "f":
			convertedValue, err := strconv.ParseFloat(valueRaw, 64)
			if err != nil {
				fmt.Println(err)
				failed = true
				break
			}
			values = append(values, convertedValue)

		default:
			values = append(values, "'"+r.FormValue(ids)+"'")
		}
	}

	if !failed {
		db.NewRow(tableName, values)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusMovedPermanently)
}

func Query(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	tableName := fmt.Sprintf("%v", t["tableName"])

	limit, err := strconv.Atoi(fmt.Sprintf("%v", t["limit"]))
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	orderBy := fmt.Sprintf("%v", t["orderBy"])
	orderDirec := fmt.Sprintf("%v", t["orderDirec"])

	data, err := db.Query(tableName, limit, orderBy, orderDirec)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func DeleteOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Could not get the id from the request url params", http.StatusBadRequest)
		return
	}

	table_name := r.URL.Query().Get("table_name")

	if table_name == "" {
		http.Error(w, "Could not get the table_name from the request url params", http.StatusBadRequest)
		return
	}

	if db.DeleteOne(table_name, id) != nil {
		http.Error(w, "Could not delete the especified item.", http.StatusInternalServerError)
	}

}
