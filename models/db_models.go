package models

type Column struct {
	Name   string `json:"name"`
	TypeDB string `json:"type"`
}

type Table struct {
	Name    string
	Columns []Column
}
