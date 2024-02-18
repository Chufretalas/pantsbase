package models

// response to the GET /api/tables endpoint
type TableResponse struct {
	TableName string   `json:"table_name"`
	Columns   []Column `json:"columns"`
}
