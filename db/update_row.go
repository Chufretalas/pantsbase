package db

import "fmt"

func UpdateRow(tableName string, values map[string]any, rowId string) error {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	queryInputValues := make([]any, 0, len(values))

	queryStr := fmt.Sprintf("UPDATE %v \n SET ", "\""+tableName+"\"")
	for idx, col := range cols {

		value := SanitizeValue(values[col.Name], col.TypeDB)

		queryStr += fmt.Sprintf(`"%v" = ?`, col.Name)
		queryInputValues = append(queryInputValues, value)

		if idx != len(cols)-1 {
			queryStr += ",\n"
		} else {
			queryStr += "\n"
		}
	}

	queryStr += "WHERE id = ?;"
	queryInputValues = append(queryInputValues, rowId)

	fmt.Printf("queryStr: %v\n", queryStr)
	fmt.Printf("queryInputValues: %v\n", queryInputValues)

	_, err = DB.Exec(queryStr, queryInputValues...)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
