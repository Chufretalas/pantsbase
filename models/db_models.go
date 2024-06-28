package models

// TODO: eliminate Column or Schema, it's confusing having both
type Column struct {
	Name   string `json:"name"`
	TypeDB string `json:"type"`
}

type Table struct {
	Name    string
	Columns []Column
}

// type Schema struct {
// 	ColName   string
// 	Type      string // INTEGER, REAL or TEXT
// 	InputName string // for the forms when passing to a template
// }
