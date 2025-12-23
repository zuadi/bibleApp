package models

import "strings"

type DBBook struct {
	Name   string
	Number int
}

func (b *DBBook) GetBookNameTillIndex(index int) string {
	return b.Name[:index]
}

func (b *DBBook) GetTrimmedBookName() string {
	return strings.TrimSpace(b.Name)
}
