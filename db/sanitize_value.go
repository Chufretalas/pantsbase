package db

import (
	"fmt"
	"strconv"
)

// TODO: use this on update endpoints as well
// avoids TEXTs on number fields putting nils in it's place
func SanitizeValue(value any, col_type string) any {

	if col_type == "TEXT" || value == nil {
		return value
	}

	str_value := fmt.Sprintf("%v", value)

	if col_type == "INTEGER" {
		value_int, err := strconv.Atoi(str_value)
		if err != nil {
			return nil
		}
		return value_int
	}

	if col_type == "REAL" {
		value_float, err := strconv.ParseFloat(str_value, 64)
		if err != nil {
			return nil
		}
		return value_float
	}

	return nil
}
