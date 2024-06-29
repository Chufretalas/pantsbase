package db

import "log"

func GetAllTableNames() []string {
	selectAllTables, err := DB.Query("select name from sqlite_schema where name not like 'sqlite%'")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer selectAllTables.Close()

	tables := make([]string, 0)

	for selectAllTables.Next() {
		var table string
		err := selectAllTables.Scan(&table)
		if err != nil {
			log.Fatal(err.Error())
		}
		tables = append(tables, table)
	}
	return tables
}
