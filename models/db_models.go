package models

type Column struct {
	Name   string `json:"name"`
	TypeDB string `json:"type"`
}

type Table struct {
	Name    string
	Columns []Column
}

// TODO: shcmea should actually be a []Column, but it maight be too hard to fix now, maybe one day
type Schema struct {
	ColName string
	Type    string // INT, REAL or TEXT
	Id      string // for the forms when passing to a template
}
