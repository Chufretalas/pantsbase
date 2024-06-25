package db

import (
	"fmt"
)

// TODO: make something to add default zero values instead of returning an error on nil values
// TODO: do something when the type of the received value does not match the column type
func NewRows(tableName string, rows []map[string]any) error {

	if len(rows) == 0 {
		return nil
	}

	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
	}

	queryStr := fmt.Sprintf("INSERT INTO \"%v\" (", tableName)
	for idx, col := range cols {
		if idx != len(cols)-1 {
			queryStr += "\"" + col.ColName + "\", "
		} else {
			queryStr += "\"" + col.ColName + "\")\n VALUES "
		}
	}

	for rowIdx, row := range rows {
		queryStr += "("
		for idx, col := range cols {
			value, ok := row[col.ColName]
			if !ok || value == nil {
				return fmt.Errorf("missing column from a row of column value is null. row index: %v, column name: %v", rowIdx, col.ColName)
			}

			if col.Type == "TEXT" {
				queryStr += fmt.Sprintf("'%v'", value)
			} else {
				queryStr += fmt.Sprintf("%v", value)
			}

			if idx != len(cols)-1 {
				queryStr += ","
			}
		}
		queryStr += ")"
		if rowIdx != len(rows)-1 {
			queryStr += ","
		} else {
			queryStr += ";"
		}
	}

	fmt.Printf("queryStr: %v\n", queryStr)
	_, err = DB.Exec(queryStr)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
