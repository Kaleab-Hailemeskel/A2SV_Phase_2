package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	myLibrary := services.Library{
		Books:   make(map[int]*models.Book),
		Members: make(map[int]*models.Member),
	}

	myLibrary.Books[1001] = &models.Book{
		ID:     1001,
		Title:  "ኦሮማይ",
		Author: "በዓሉ ግርማ",
		Status: models.BookAvailable,
	}
	myLibrary.Books[1002] = &models.Book{
		ID:     1002,
		Title:  "አዲስ አበባ",
		Author: "አዳም ረታ",
		Status: models.BookAvailable,
	}
	myLibrary.Books[1003] = &models.Book{
		ID:     1003,
		Title:  "ፍቅር እስከ መቃብር",
		Author: "ሐዲስ ዓለማየሁ",
		Status: models.BookAvailable,
	}
	myLibrary.Books[1004] = &models.Book{
		ID:     1004,
		Title:  "የተቀበረው ምሥጢር",
		Author: "ይስማዕከ ወርቁ",
		Status: models.BookAvailable,
	}
	myLibrary.Books[1005] = &models.Book{
		ID:     1005,
		Title:  "ደራሲው",
		Author: "ካሌብ ግዛው",
		Status: models.BookAvailable,
	}

	myLibrary.Members[2001] = &models.Member{
		ID:            2001,
		Name:          "አበበ በለጠ",
		BorrowedBooks: []*models.Book{}, // Initially empty slice
	}
	myLibrary.Members[2002] = &models.Member{
		ID:            2002,
		Name:          "ፋጡማ ዑመር",
		BorrowedBooks: []*models.Book{}, // Initially empty slice
	}
	myLibrary.Members[2003] = &models.Member{
		ID:            2003,
		Name:          "ዳዊት ተስፋዬ",
		BorrowedBooks: []*models.Book{}, // Initially empty slice
	}
	controllers.StartLibrary(&myLibrary)
}
