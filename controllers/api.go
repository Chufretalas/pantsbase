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

	tableName := mux.Vars(r)["table_name"]

	limitStr := r.URL.Query().Get("limit")

	orderBy := r.URL.Query().Get("order_by")
	orderDirection := r.URL.Query().Get("order_direction") // Needs to be "ASC" or "DESC"

	limit, err := strconv.Atoi(fmt.Sprintf("%v", limitStr))
	if err != nil && limitStr != "" {
		fmt.Println(err)
		http.Error(w, `value for 'limit' was not a valid number`, http.StatusBadRequest)
		return
	}

	data, err := db.Query(tableName, limit, orderBy, orderDirection)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not retrieve the requested data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func QueryOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	tableName := vars["table_name"]
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `id was not a valid number`, http.StatusBadRequest)
		return
	}

	data, err := db.QueryOne(tableName, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not retrieve the requested data", http.StatusInternalServerError)
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

func Tables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	onlyNames := r.URL.Query().Has("only_names")
	allNames := db.GetAllTableNames()

	if onlyNames {
		json.NewEncoder(w).Encode(allNames)
		return
	}

	resp := make([]models.TableResponse, 0, len(allNames))
	for _, name := range allNames {
		schema, err := db.GetSchema(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode([]any{})
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

func UpdateOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	tableName := vars["table_name"]

	schema, err := db.GetSchema(tableName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "invalid table name", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "The id is not a valid integer", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var body map[string]any
	err = decoder.Decode(&body)
	if err != nil {
		fmt.Println("error decoding the body\n", err)
		http.Error(w, "the body was not valid json", http.StatusBadRequest)
		return
	}

	// TODO: validate that number fields are actually numbers maybe in the db.UpdateRow, but it might not be necessary
	// parsing the request body for the table fields
	values := make([]any, 0, len(schema))
	for _, column := range schema {
		newValue := body[column.ColName]
		values = append(values, newValue)
	}

	err = db.UpdateRow(tableName, values, fmt.Sprint(id))

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not update the requested row", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	changedRow, err := db.QueryOne(tableName, id)

	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]any{"message": "could not retrieve the changed row, but it was updated correctly"})
	} else {
		json.NewEncoder(w).Encode(changedRow)
	}

}

// TODO: post (many and one should be the same endpoint)
// TODO: update many
