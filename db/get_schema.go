package db

import (
	"fmt"
	"regexp"
	"strings"
)

var SchemaRegex *regexp.Regexp

func GetSchema(tableName string) ([][]string, error) {
	fmt.Println(tableName)

	stmt, err := DB.Prepare("SELECT sql FROM sqlite_schema WHERE name = ?;")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(tableName)

	if err != nil {
		return nil, err
	}

	parsedCols := make([][]string, 0)

	for rows.Next() {
		var rawSchema string
		err := rows.Scan(&rawSchema)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		rawCols := SchemaRegex.FindAllStringSubmatch(rawSchema, -1) // the -1 means to return all substrings matched
		for _, rawRow := range rawCols {
			parsedCols = append(parsedCols, strings.Split(rawRow[0], " "))
		}
	}

	return parsedCols, nil
}
