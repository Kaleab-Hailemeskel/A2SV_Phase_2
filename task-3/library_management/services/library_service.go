package services

import (
	"fmt"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book) error
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]*models.Book
	Members map[int]*models.Member
}

// Errors
type BookNotFoundError int
type MemberNotFoundError int
type LibraryNotFound Library
type BookNotBorrowedError int
type MemberNotHaveBookError int
type BookAlreadyExistError int
type BookAlreadyBorrowedError int

// Error methods
func (val BookNotFoundError) Error() string {
	return fmt.Sprintf("Error: Book with an ID %v doesn't exist", int(val))
}
func (val BookAlreadyBorrowedError) Error() string {
	return fmt.Sprintf("Error: Book with an ID %v is already borrowed", int(val))
}
func (val MemberNotFoundError) Error() string {
	return fmt.Sprintf("Error: Member with an ID %v doesn't exist", int(val))
}
func (val BookAlreadyExistError) Error() string {
	return fmt.Sprintf("Error: Book with an ID %v already exist", int(val))
}
func (li *LibraryNotFound) Error() string {
	return fmt.Sprintf("Uninitialized Library was found %T", li)
}
func (val BookNotBorrowedError) Error() string {
	return fmt.Sprintf("Error: Book with ID %v wasn't Borrowed", int(val))
}
func (val MemberNotHaveBookError) Error() string {
	return fmt.Sprintf("Error: Member with ID %v try to return not borrowed book", int(val))
}

var bookBorrowedByMember = func(listBook []*models.Book, bookID int) int {
	for index, eachBook := range listBook {
		if eachBook.ID == bookID {
			return index
		}
	}
	return -1
}

// Error Ends Here

func (li *Library) RemoveBook(bookID int) error {
	if li == nil {
		return &LibraryNotFound{}
	}
	if book, exist := li.Books[bookID]; !exist {
		return BookNotFoundError(bookID)
	} else if book.Status == models.BookBorrowed {
		return BookAlreadyBorrowedError(bookID)
	}
	delete(li.Books, bookID)
	return nil
}
func (li *Library) BorrowBook(bookID int, memberID int) error {
	if li == nil {
		return &LibraryNotFound{}
	}
	if book, bookExists := li.Books[bookID]; !bookExists {
		return BookNotFoundError(bookID)
	} else if book.Status == models.BookBorrowed {
		return BookAlreadyBorrowedError(bookID)
	}
	if _, memberExists := li.Members[memberID]; !memberExists {
		return MemberNotFoundError(memberID)
	}

	// directly change the status of the book
	bookToModify := li.Books[bookID]
	bookToModify.Status = models.BookBorrowed

	// reassigning the Member after borrowing a book
	borrowMember := li.Members[memberID]
	borrowMember.BorrowedBooks = append(borrowMember.BorrowedBooks, bookToModify)

	return nil
}
func (li *Library) ReturnBook(bookID int, memberID int) error {
	if li == nil {
		return &LibraryNotFound{}
	}
	if book, bookExists := li.Books[bookID]; !bookExists {
		return BookNotFoundError(bookID)
	} else if book.Status == models.BookAvailable {
		return BookNotBorrowedError(bookID)
	}

	if member, memberExists := li.Members[memberID]; !memberExists {
		return MemberNotFoundError(memberID)
	} else if bookIndex := bookBorrowedByMember(member.BorrowedBooks, bookID); bookIndex == -1 { //-1 == not found
		return MemberNotHaveBookError(member.ID)
	} else {
		// todo: the else is needed so that to access the valud of bookIndex
		borrowList := li.Members[memberID].BorrowedBooks
		borrowList[bookIndex].Status = models.BookAvailable
		li.Members[memberID].BorrowedBooks = append(borrowList[:bookIndex], borrowList[bookIndex+1:]...)
	}

	return nil
}
func (li *Library) ListAvailableBooks() []models.Book {
	if li == nil || li.Books == nil || len(li.Books) == 0 {
		return nil
	}
	res := []models.Book{}
	for _, book := range li.Books {
		if book.Status == models.BookAvailable {
			res = append(res, *book)
		}
	}
	return res
}
func (li *Library) ListBorrowedBooks(memberID int) []models.Book {
	if li == nil || li.Members == nil || len(li.Members) == 0 {
		return nil
	}
	if _, memberExist := li.Members[memberID]; !memberExist {
		return nil
	}
	res := make([]models.Book, len(li.Members[memberID].BorrowedBooks))
	for _, borrowedBook := range li.Members[memberID].BorrowedBooks {
		res = append(res, *borrowedBook)
	}
	return res
}
func (li *Library) AddBook(book *models.Book) error {
	if li == nil {
		return &LibraryNotFound{}
	}
	if _, bookExist := li.Books[book.ID]; bookExist {
		return BookAlreadyExistError(book.ID)
	}
	li.Books[book.ID] = book
	return nil
}
