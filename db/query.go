package db

import (
	"fmt"
)

// Leave orderBy empty if you don't want sorting | orderDirection should be "ASC" or "DESC", default is DESC
func Query(tableName string, limit int, orderBy string, orderDirection string) ([]map[string]any, error) {

	schema, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	willOrder := false

	for _, col := range schema {
		if orderBy == col.Name {
			willOrder = true
			break
		}
	}

	if willOrder && orderDirection != "DESC" && orderDirection != "ASC" {
		orderDirection = "DESC"
	}

	var queryString string

	if willOrder {
		queryString = fmt.Sprintf(`SELECT * FROM "%v" ORDER BY "%v" %v`, tableName, orderBy, orderDirection)
	} else {
		queryString = fmt.Sprintf(`SELECT * FROM "%v"`, tableName)
	}

	if limit > 0 {
		queryString += fmt.Sprintf(` LIMIT %v`, limit)
	}

	return queryAndReadData(queryString)
}

func QueryOne(tableName string, id int) ([]map[string]any, error) {
	queryString := fmt.Sprintf(`SELECT * FROM "%v" WHERE id=?;`, tableName)

	return queryAndReadData(queryString, id)
}

func queryAndReadData(queryString string, args ...any) ([]map[string]any, error) {
	rows, err := DB.Query(queryString, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	count := len(columns)
	values := make([]any, count)
	scanArgs := make([]any, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	masterData := make([]map[string]any, 0)

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		entrie := make(map[string]any)

		for i, v := range values {
			entrie[columns[i]] = v
		}
		masterData = append(masterData, entrie)
	}

	return masterData, nil
}
