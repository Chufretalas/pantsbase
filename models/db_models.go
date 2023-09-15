package models

type Column struct {
	Name   string
	TypeDB string
}

type Table struct {
	Name    string
	Columns []Column
}
