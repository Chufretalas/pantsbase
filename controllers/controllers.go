package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Chufretalas/pantsbase/db"
	m "github.com/Chufretalas/pantsbase/models"
	"github.com/gorilla/mux"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "index", db.GetAllTableNames())
}

// TODO: disallow non-ASCII characters on table columns
func NewTable(w http.ResponseWriter, r *http.Request) {
	tableName := strings.Trim(r.FormValue("name"), " ")

	var columnIndexes []string
	if r.FormValue("column_indexes") == "" {
		columnIndexes = make([]string, 0)
	} else {
		columnIndexes = strings.Split(r.FormValue("column_indexes"), " ")
	}

	columns := make([]m.Column, 0, len(columnIndexes))
	for _, index := range columnIndexes {
		name := strings.Trim(r.FormValue(fmt.Sprintf("n%v", index)), " ")
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

	tableName := r.FormValue("table_name")
	fieldIds := strings.Split(r.FormValue("new_row_ids"), " ") // it's actually the field names not ids, I'm just dense and I am not changing it now
	fmt.Println(fieldIds)

	values, err := extractRowFormValues(r, fieldIds)

	if err == nil {
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

	tableName := mux.Vars(r)["table_name"]

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
	vars := mux.Vars(r)

	table_name := vars["table_name"]
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "The id is not a valid integer", http.StatusBadRequest)
		return
	}

	if db.DeleteOne(table_name, id) != nil {
		http.Error(w, "Could not delete the especified item.", http.StatusInternalServerError)
	}

}

func UpdateRow(w http.ResponseWriter, r *http.Request) {

	tableName := r.FormValue("update_table_name")
	rowId := r.FormValue("update_row_id")
	fieldIds := strings.Split(r.FormValue("update_row_ids"), " ")

	values, err := extractRowFormValues(r, fieldIds)

	if err == nil {
		db.UpdateRow(tableName, values, rowId)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusMovedPermanently)
}

func extractRowFormValues(r *http.Request, fieldIds []string) ([]interface{}, error) {
	failed := false
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

	if failed {
		return nil, errors.New("error parsing the form")
	}

	return values, nil
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tableName := vars["table_name"]

	_, err := db.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS \"%v\";", tableName))

	if err != nil {
		log.Println(err)
	}
}
