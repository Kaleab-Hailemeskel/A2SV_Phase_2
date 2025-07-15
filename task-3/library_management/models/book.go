package models

import "fmt"

const (
	BookAvailable = "Available"
	BookBorrowed  = "Borrowed"
)

func (book Book) String() string {
	return fmt.Sprintf("%d | %s | %s | %s", book.ID, book.Title, book.Author, book.Status)
}

type Book struct {
	ID     int
	Title  string
	Author string
	Status string // can be 'Available' OR 'Borrowed'
}
