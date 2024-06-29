package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/models"
)

// receives the columns in a json where the key is the column name and the value is the type
func NewTable(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("table_name")

	decoder := json.NewDecoder(r.Body)
	var body map[string]string
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println("error decoding the body\n", err)
		http.Error(w, "the body was not valid json", http.StatusBadRequest)
		return
	}

	cols := make([]models.Column, 0, len(body))

	for k, v := range body {
		if v == "TEXT" || v == "REAL" || v == "INT" {
			cols = append(cols, models.Column{Name: k, TypeDB: v})
		} else {
			http.Error(w, `column type must be "TEXT", "REAL" or "INT" and nothing else`, http.StatusBadRequest)
			return
		}
	}

	err = db.NewTable(tableName, cols)
	if err != nil {
		fmt.Println("error decoding the body\n", err)
		http.Error(w, "something went wrong when creating th enew table", http.StatusInternalServerError)
		return
	}
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {

	tableName := r.PathValue("table_name")

	_, err := db.DB.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%v";`, tableName))

	if err != nil {
		log.Println(err)
	}
}

func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tableName := r.PathValue("table_name")

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

	fmt.Printf("data: %v\n", data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func QueryOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tableName := r.PathValue("table_name")

	id, err := strconv.Atoi(r.PathValue("id"))
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

	table_name := r.PathValue("table_name")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "The id is not a valid integer", http.StatusBadRequest)
		return
	}

	if db.DeleteOne(table_name, id) != nil {
		http.Error(w, "Could not delete the especified item.", http.StatusInternalServerError)
	}

}

func GetTables(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "failed to fetch the tables schemas", http.StatusInternalServerError)
			return
		}
		resp = append(resp, models.TableResponse{TableName: name, Columns: schema})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UpdateOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tableName := r.PathValue("table_name")
	id, err := strconv.Atoi(r.PathValue("id"))

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

	err = db.UpdateRow(tableName, body, fmt.Sprint(id))

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not update the requested row", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	changedRow, err := db.QueryOne(tableName, id)

	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]any{"message": "could not retrieve the updated row, but it was updated correctly"})
	} else {
		json.NewEncoder(w).Encode(changedRow)
	}

}

func PostRows(w http.ResponseWriter, r *http.Request) {

	tableName := r.PathValue("table_name")

	// duplicating the request body so it can be decoded twice
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	decoder1 := json.NewDecoder(tee)
	decoder2 := json.NewDecoder(&buf)

	// this handler can receive both only one row as well as multiple
	var simple_body map[string]any
	var multi_body []map[string]any
	simple_body_err := decoder1.Decode(&simple_body)
	multi_body_err := decoder2.Decode(&multi_body)
	if simple_body_err != nil && multi_body_err != nil {
		fmt.Println("error decoding the body\n", simple_body_err, multi_body_err)
		http.Error(w, "the body was not valid json", http.StatusBadRequest)
		return
	}

	var dbErr error
	var ids []int
	if simple_body_err != nil {
		ids, dbErr = db.NewRows(tableName, multi_body)
	} else {
		ids, dbErr = db.NewRows(tableName, []map[string]any{simple_body})
	}
	if dbErr != nil {
		fmt.Println(dbErr)
		http.Error(w, "somenthing went wrong while saving the new data", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]any{"ids": ids})
}

// TODO: update many?
