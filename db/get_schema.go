package db

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Chufretalas/pantsbase/models"
)

var SchemaRegex *regexp.Regexp

func GetSchema(tableName string) ([]models.Schema, error) {
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

	parsedCols := make([]models.Schema, 0)

	for rows.Next() {
		var rawSchema string
		err := rows.Scan(&rawSchema)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		rawCols := SchemaRegex.FindAllStringSubmatch(rawSchema, -1) // the -1 means to return all substrings matched
		for i, rawRow := range rawCols {
			sepIndex := strings.LastIndex(rawRow[0], " ")
			name := strings.Trim(rawRow[0][:sepIndex], "\"")
			colType := rawRow[0][sepIndex+1:]
			var id string
			switch colType {
			case "INTEGER":
				id = fmt.Sprintf("i%v", i)
			case "REAL":
				id = fmt.Sprintf("r%v", i)
			case "TEXT":
				id = fmt.Sprintf("t%v", i)

			}
			parsedCols = append(parsedCols, models.Schema{ColName: name, Type: colType, Id: id})
		}
	}

	return parsedCols, nil
}
