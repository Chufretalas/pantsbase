package main

import (
	"bytes"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"path"
	"slices"
	"testing"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/routes"
)

func TestCRUD(t *testing.T) {
	dir := t.TempDir()
	t.Log("temp dir: ", dir)
	db.ConnectDB(path.Join(dir, "data.db"))
	defer db.DB.Close()

	err := db.DB.Ping()

	if err != nil {
		t.Error("Database not initialized.", err.Error())
	}

	r := http.NewServeMux()
	routes.LoadRoutes(r)

	// Creating table
	w := httptest.NewRecorder()

	newTableBody, _ := json.Marshal(map[string]string{"text": "TEXT", "int": "INTEGER", "float": "REAL"})
	req := httptest.NewRequest(http.MethodPost, "/api/new_table/test1", bytes.NewReader(newTableBody))
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed creating the table. Status code: %v. Message: %v", w.Code, w.Body)
	}

	// Creating itens
	inputRows := []any{map[string]any{"int": 17, "float": 5.3, "text": "example text"}, []map[string]any{{"int": nil, "float": 99.589, "text": "example text1"}, {"int": -200, "float": nil, "text": "example text2"}, {"int": 2, "float": 8.2, "text": nil}}}

	for idx, b := range inputRows {
		w := httptest.NewRecorder()
		newTableBody, _ := json.Marshal(b)
		req := httptest.NewRequest(http.MethodPost, "/api/tables/test1", bytes.NewReader(newTableBody))
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			if idx == 0 {
				t.Errorf("failed posting a single row. Status code: %v. Message: %v", w.Code, w.Body)
			}
			t.Errorf("failed posting multiple rows. Status code: %v. Message: %v", w.Code, w.Body)
		}
	}

	// Querying one
	w = httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/api/tables/test1/2", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed getting a single row. Status code: %v. Message: %v", w.Code, w.Body)
	}

	var respBodySingle map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &respBodySingle)
	if err != nil {
		t.Errorf("failed parsing the body while geting a single row. err: %v", err.Error())
	}

	singleGetExpected := map[string]any{"id": float64(2), "int": nil, "float": float64(99.589), "text": "example text1"}
	if !maps.Equal(respBodySingle, singleGetExpected) {
		t.Errorf("failed while getting a single row. Expected: %v. Received: %v", singleGetExpected, respBodySingle)
	}

	//Querying many

	testsMultiGet := []struct {
		url      string
		expected []map[string]any
	}{
		{"/api/tables/test1", []map[string]any{{"id": float64(1), "int": float64(17), "float": float64(5.3), "text": "example text"}, {"id": float64(2), "int": nil, "float": float64(99.589), "text": "example text1"}, {"id": float64(3), "int": float64(-200), "float": nil, "text": "example text2"}, {"id": float64(4), "int": float64(2), "float": float64(8.2), "text": nil}}},
		{"/api/tables/test1?limit=2", []map[string]any{{"id": float64(1), "int": float64(17), "float": float64(5.3), "text": "example text"}, {"id": float64(2), "int": nil, "float": float64(99.589), "text": "example text1"}}},
		{"/api/tables/test1?order_by=int&order_direction=DESC", []map[string]any{{"id": float64(1), "int": float64(17), "float": float64(5.3), "text": "example text"}, {"id": float64(4), "int": float64(2), "float": float64(8.2), "text": nil}, {"id": float64(3), "int": float64(-200), "float": nil, "text": "example text2"}, {"id": float64(2), "int": nil, "float": float64(99.589), "text": "example text1"}}},
	}

	for _, tt := range testsMultiGet {
		var respBodyMultiple []map[string]any
		w = httptest.NewRecorder()

		req = httptest.NewRequest(http.MethodGet, tt.url, nil)
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("failed getting multiple rows. URL: %v. Status code: %v. Message: %v", tt.url, w.Code, w.Body)
		}

		err = json.Unmarshal(w.Body.Bytes(), &respBodyMultiple)
		if err != nil {
			t.Errorf("failed parsing the body while geting a multiple rows. err: %v", err.Error())
		}

		if len(tt.expected) != len(respBodyMultiple) {
			t.Errorf("failed getting multiple rows. Expected: %v. Received: %v", tt.expected, respBodyMultiple)
			continue
		}

		for idx, row := range respBodyMultiple {
			if !maps.Equal(row, tt.expected[idx]) {
				t.Errorf("failed getting multiple rows. Expected: %v. Received: %v", tt.expected, respBodyMultiple)
				break
			}
		}
	}

	// Updating
	w = httptest.NewRecorder()
	updateBody, _ := json.Marshal(map[string]any{"int": 999})
	req = httptest.NewRequest(http.MethodPut, "/api/tables/test1/2", bytes.NewReader(updateBody))
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed updating a row. Status code: %v. Message: %v", w.Code, w.Body)
	}

	var respBodyUpdate map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &respBodyUpdate)
	if err != nil {
		t.Errorf("failed parsing the body while geting a single row. err: %v", err.Error())
	}

	updateExpected := map[string]any{"id": float64(2), "int": float64(999), "float": float64(99.589), "text": "example text1"}
	if !maps.Equal(respBodyUpdate, updateExpected) {
		t.Errorf("failed while getting a single row. Expected: %v. Received: %v", updateExpected, respBodyUpdate)
	}

	// Deleting Row
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/tables/test1/2", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed deleting a row. Status code: %v. Message: %v", w.Code, w.Body)
	}

	w = httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/api/tables/test1/2", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 && w.Body.String() != "no row for id = 2 was found" {
		t.Error("failed deleting a single row. quering for GET /api/tables/test1/2 does not report the deleted row as missing")
	}

	// Deleting Table
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/delete_table/test1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed deleting the table. Status code: %v. Message: %v", w.Code, w.Body)
	}

	w = httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/api/tables?only_names", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("failed GET /api/tables?only_names. Status code: %v. Message: %v", w.Code, w.Body)
	}

	var tableNames []string

	err = json.Unmarshal(w.Body.Bytes(), &tableNames)
	if err != nil {
		t.Errorf("failed parsing the body of  GET /api/tables?only_names. err: %v", err.Error())
	}

	if slices.Contains(tableNames, "test1") {
		t.Error("failed deleting the table. quering for GET /api/tables?only_names reports the deleted table as existing")
	}
}
