package db

import (
	"fmt"
)

// Leave orderBy empty if you don't want sorting | orderDirection should be "ASC" or "DESC", default is DESC
func Query(tableName string, limit int, orderBy string, orderDirection string) ([]map[string]interface{}, error) {

	schema, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	willOrder := false

	for _, col := range schema {
		if orderBy == col.ColName {
			willOrder = true
			break
		}
	}

	if willOrder && orderDirection != "DESC" && orderDirection != "ASC" {
		orderDirection = "DESC"
	}

	if limit <= 0 {
		limit = 1
	}

	var queryText string

	if willOrder {
		queryText = fmt.Sprintf(`SELECT * FROM "%v" ORDER BY "%v" %v LIMIT %v`, tableName, orderBy, orderDirection, limit)
	} else {
		queryText = fmt.Sprintf(`SELECT * FROM "%v" LIMIT %v`, tableName, limit)
	}

	rows, err := DB.Query(queryText)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	masterData := make([]map[string]interface{}, 0)

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		entrie := make(map[string]interface{})

		for i, v := range values {
			entrie[columns[i]] = v
		}
		masterData = append(masterData, entrie)
	}

	return masterData, nil
}
