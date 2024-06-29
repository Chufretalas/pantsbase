package db

import (
	"fmt"

	"github.com/Chufretalas/pantsbase/models"
)

func NewTable(tableName string, columns []models.Column) error {
	queryStr := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%v" (\n id INTEGER PRIMARY KEY AUTOINCREMENT`, tableName)
	if len(columns) != 0 {
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

	_, err := DB.Exec(queryStr)

	return err
}
