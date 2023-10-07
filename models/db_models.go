package models

type Column struct {
	Name   string
	TypeDB string
}

type Table struct {
	Name    string
	Columns []Column
}

type Schema struct {
	ColName string
	Type    string
	Id      string // for the forms when passing to a template
}
