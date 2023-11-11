package db

import "fmt"

func DeleteOne(table_name string, id int) error {
	_, err := DB.Exec(fmt.Sprintf(`DELETE FROM [%v] WHERE id = %v ;`, table_name, id))

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Table name: " + table_name)
		fmt.Println("id: " + fmt.Sprintf("%v", id))
		return err
	}

	return nil
}
