package db

import (
	"fmt"
)

// TODO: make this safer with ?'s to avoid SQL injections (args...)
// TODO: return errors in more cases of an error happening here
// inserts new data into the table and returns the ids of the new rows
func NewRows(tableName string, rows []map[string]any) ([]int, error) {

	if len(rows) == 0 {
		return []int{}, nil
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
			if !ok {
				return []int{}, fmt.Errorf("missing column from a row. row index: %v, column name: %v", rowIdx, col.ColName)
			}

			value = SanitizeValue(value, col.Type)

			if value == nil {
				queryStr += "NULL"
			} else {
				if col.Type == "TEXT" {
					queryStr += fmt.Sprintf("'%v'", value)
				} else {
					queryStr += fmt.Sprintf("%v", value)
				}
			}

			if idx != len(cols)-1 {
				queryStr += ","
			}
		}
		queryStr += ")"
		if rowIdx != len(rows)-1 {
			queryStr += ","
		} else {
			queryStr += " RETURNING id;"
		}
	}

	fmt.Printf("queryStr: %v\n", queryStr)

	var ids []int
	returned_rows, err := DB.Query(queryStr)
	if err != nil {
		fmt.Println(err)
	}
	for returned_rows.Next() {
		var id int
		if err := returned_rows.Scan(&id); err != nil {
			fmt.Println(err)
		}
		ids = append(ids, id)
	}
	fmt.Println("IDs: ", ids)

	return ids, nil
}
