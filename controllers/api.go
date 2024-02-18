package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/models"
	"github.com/gorilla/mux"
)

func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println("error decoding the body\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "the body was not valid json"})
		return
	}

	tableName := mux.Vars(r)["table_name"]

	limitStr := fmt.Sprintf("%v", t["limit"])

	if limitStr == "<nil>" {
		limitStr = "50"
	}

	limit, err := strconv.Atoi(fmt.Sprintf("%v", limitStr))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": `value for 'limit' was not a valid number`})
		return
	}

	orderBy := fmt.Sprintf("%v", t["orderBy"])
	orderDirec := fmt.Sprintf("%v", t["orderDirec"])

	data, err := db.Query(tableName, limit, orderBy, orderDirec)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func DeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	table_name := vars["table_name"]
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "The id is not a valid integer", http.StatusBadRequest)
		return
	}

	if db.DeleteOne(table_name, id) != nil {
		http.Error(w, "Could not delete the especified item.", http.StatusInternalServerError)
	}

}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tableName := vars["table_name"]

	_, err := db.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS \"%v\";", tableName))

	if err != nil {
		log.Println(err)
	}
}

// TODO: add parameters to just return the table names
func Tables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allNames := db.GetAllTableNames()
	resp := make([]models.TableResponse, 0, len(allNames))
	for _, name := range allNames {
		schema, err := db.GetSchema(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode([]interface{}{})
			return
		}
		columns := make([]models.Column, 0, len(schema))
		columns = append(columns, models.Column{Name: "id", TypeDB: "INTEGER"})
		for _, c := range schema {
			columns = append(columns, models.Column{Name: c.ColName, TypeDB: c.Type})
		}
		resp = append(resp, models.TableResponse{TableName: name, Columns: columns})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// TODO: update row
