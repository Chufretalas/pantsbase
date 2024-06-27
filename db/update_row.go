package db

import "fmt"

// TODO: make this safer with ?'s to avoid SQL injections (args...)
func UpdateRow(tableName string, values map[string]any, rowId string) error {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	queryStr := fmt.Sprintf("UPDATE %v \n SET ", "\""+tableName+"\"")
	for idx, col := range cols {
		value := SanitizeValue(values[col.ColName], col.Type)
		if value == nil {
			queryStr += fmt.Sprintf("\"%v\" = NULL", col.ColName)
		} else {
			if col.Type == "TEXT" {
				queryStr += fmt.Sprintf("\"%v\" = '%v'", col.ColName, value)
			} else {
				queryStr += fmt.Sprintf("\"%v\" = %v", col.ColName, value)
			}
		}
		if idx != len(cols)-1 {
			queryStr += ",\n"
		} else {
			queryStr += "\n"
		}
	}

	queryStr += fmt.Sprintf("WHERE id = %v;", rowId)

	fmt.Println(queryStr)

	_, err = DB.Exec(queryStr)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
