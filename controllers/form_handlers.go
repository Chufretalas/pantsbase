package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Chufretalas/pantsbase/db"
	m "github.com/Chufretalas/pantsbase/models"
)

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
