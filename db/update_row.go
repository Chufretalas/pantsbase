package db

import "fmt"

// the values need to be in the order as they appear in the db.GetSchema function
// values with nil will be ignored
func UpdateRow(tableName string, values []any, rowId string) error {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	queryStr := fmt.Sprintf("UPDATE %v \n SET ", "\""+tableName+"\"")
	for idx, col := range cols {
		if values[idx] == nil {
			continue
		}
		queryStr += fmt.Sprintf("\"%v\" = %v", col.ColName, values[idx])
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
