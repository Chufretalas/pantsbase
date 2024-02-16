package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Chufretalas/pantsbase/db"
	m "github.com/Chufretalas/pantsbase/models"
)

// utils

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

// end utils

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

	err := db.NewTable(tableName, columns)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
