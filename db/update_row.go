package db

import (
	"fmt"
	"slices"

	"github.com/Chufretalas/pantsbase/models"
)

func UpdateRow(tableName string, values map[string]any, rowId string) error {
	cols, err := GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	relevantKeys := make([]string, 0, len(values))
	for k, _ := range values {
		if slices.ContainsFunc(cols, func(col models.Column) bool {
			return col.Name == k
		}) {
			relevantKeys = append(relevantKeys, k)
		}
	}

	if len(relevantKeys) == 0 {
		return nil
	}

	queryInputValues := make([]any, 0, len(values))

	queryStr := fmt.Sprintf("UPDATE %v \n SET ", "\""+tableName+"\"")
	for idx, key := range relevantKeys {

		col := cols[slices.IndexFunc(cols, func(col models.Column) bool {
			return col.Name == key
		})]

		value := SanitizeValue(values[col.Name], col.TypeDB)

		queryStr += fmt.Sprintf(`"%v" = ?`, col.Name)
		queryInputValues = append(queryInputValues, value)

		if idx != len(relevantKeys)-1 {
			queryStr += ",\n"
		} else {
			queryStr += "\n"
		}
	}

	queryStr += "WHERE id = ?;"
	queryInputValues = append(queryInputValues, rowId)

	// fmt.Printf("queryStr: %v\n", queryStr)
	// fmt.Printf("queryInputValues: %v\n", queryInputValues)

	_, err = DB.Exec(queryStr, queryInputValues...)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
