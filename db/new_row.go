package db

import (
	"fmt"
)

func NewRow(tableName string, values []any) {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
	}

	queryStr := fmt.Sprintf("INSERT INTO \"%v\" (", tableName)
	for idx, col := range cols {
		if idx != len(cols)-1 {
			queryStr += "\"" + col.ColName + "\", "
		} else {
			queryStr += "\"" + col.ColName + "\")\n VALUES ("
		}
	}

	for idx, value := range values {
		if idx != len(cols)-1 {
			queryStr += fmt.Sprintf("%v, ", value)
		} else {
			queryStr += fmt.Sprintf("%v);", value)
		}
	}

	fmt.Println(queryStr)

	_, err = DB.Exec(queryStr)

	if err != nil {
		fmt.Println(err)
	}
}
