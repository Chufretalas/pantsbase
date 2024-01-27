package db

import "fmt"

func UpdateRow(tableName string, values []interface{}, rowId string) {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
	}

	queryStr := fmt.Sprintf("UPDATE %v \n SET ", "\""+tableName+"\"")
	for idx, col := range cols {
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
	}
}
