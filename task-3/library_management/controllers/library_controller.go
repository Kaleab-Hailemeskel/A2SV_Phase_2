package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func returnToMainMenu() bool {
	fmt.Print("Do you want to return to MainMenu (1. No   <Anything>. Yes): ")
	choice, _ := reader.ReadString('\n')
	choiceInt, err := strconv.Atoi(choice)
	if err != nil || choiceInt != 1 {
		return true
	}
	return false
}
func showListOfBooks(listOfBooks []models.Book) {
	fmt.Println("ID\tTitle\t\t\t\tAuthor\t\tStatus")
	for _, times := range []int{1, 4, 2} {
		fmt.Printf("%s%s", strings.Repeat("-", (times*8)-1), "|")
	}
	fmt.Println()
	for _, book := range listOfBooks {
		fmt.Println(book)
	}
}

func StartLibrary(library *services.Library) {
	for {
		fmt.Println("\n--- Library Management System ---")
		fmt.Println("1. Add a new book")
		fmt.Println("2. Remove an existing book")
		fmt.Println("3. Borrow a book")
		fmt.Println("4. Return a book")
		fmt.Println("5. List all available books")
		fmt.Println("6. List all borrowed books by a member")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
			case "10":
				fmt.Printf("Total Number of Books:  %d\n", len(library.Books))
				showListOfBooks(slices.Collect(maps.Values(library.Books)))
			case "11":
				fmt.Println("Total Number of Members", len(library.Members))
				for _, val := range library.Members {
					fmt.Println(val.ID, "\t\t", val.Name, "\t\t\t", val.BorrowedBooks)
				}
			case "1":
				for {
					fmt.Print("Enter Book ID: ")
					bookIDStr, _ := reader.ReadString('\n')
					bookID, err := strconv.Atoi(strings.TrimSpace(bookIDStr)) // Implement conversion
					if err != nil {
						fmt.Println("\t<Error> Invalid input for book ID")
						if returnToMainMenu() {
							break
						}
						continue
					}
					fmt.Print("Enter Book Title: ")
					bookTitle, _ := reader.ReadString('\n')
					bookTitle = strings.TrimSpace(bookTitle) // Implement reading
					fmt.Print("Enter Book Author: ")
					bookAuthor, _ := reader.ReadString('\n')
					bookAuthor = strings.TrimSpace(bookAuthor) // Implement reading
					addError := library.AddBook(models.Book{ID: bookID, Author: bookAuthor, Status: models.BookAvailable, Title: bookTitle})
					if addError == nil {
						fmt.Println(">> Book added successfully. <<") // Placeholder message
					} else {
						fmt.Println("\t", addError)
					}
					break
				}

			case "2":
				for {
					fmt.Print("Enter Book ID to remove: ")
					bookIDStr, _ := reader.ReadString('\n')
					bookID, err := strconv.Atoi(strings.TrimSpace(bookIDStr)) // Implement conversion
					if err != nil {
						fmt.Println("\t<Error> Invalid book ID")
						if returnToMainMenu() {
							break
						}
						continue
					}

					removingErr := library.RemoveBook(bookID)
					if removingErr == nil {
						fmt.Println("Book removed successfully.") // Placeholder message
					} else {
						fmt.Println("\t", removingErr)
					}
					break
				}

			case "3":
				for {
					fmt.Print("Enter Book ID to borrow: ")
					bookIDStr, _ := reader.ReadString('\n')
					bookID, convertionErr := strconv.Atoi(strings.TrimSpace(bookIDStr)) // Implement conversion
					if convertionErr != nil {
						fmt.Println("\t<Error> Invalid book ID")
						if returnToMainMenu() {
							break
						}
						continue

					}
					fmt.Print("Enter Member ID: ")
					memberIDStr, _ := reader.ReadString('\n')
					memberID, convertionErr := strconv.Atoi(strings.TrimSpace(memberIDStr)) // Implement conversion
					if convertionErr != nil {
						fmt.Println("\t<Error> Invalid member ID")
						if returnToMainMenu() {
							break
						}
						continue
					}
					borrowError := library.BorrowBook(bookID, memberID)
					if borrowError == nil {
						fmt.Printf("\t>> %s borrowed %s successfully <<\n", library.Members[memberID].Name, library.Books[bookID].Title)
					} else {
						fmt.Println("\t", borrowError)
					}
					break
				}

			case "4":
				for {
					fmt.Print("Enter Book ID to return: ")
					bookIDStr, _ := reader.ReadString('\n')
					bookID, convertionErr := strconv.Atoi(strings.TrimSpace(bookIDStr)) // Implement conversion
					if convertionErr != nil {
						if returnToMainMenu() {
							break
						}
						continue
					}
					fmt.Print("Enter Member ID: ")
					memberIDStr, _ := reader.ReadString('\n')
					memberID, convertionErr := strconv.Atoi(strings.TrimSpace(memberIDStr)) // Implement conversion
					if convertionErr != nil {
						fmt.Println("\t<Error> Invalid member ID")
						if returnToMainMenu() {
							break
						}
						continue
					}
					returnErr := library.ReturnBook(bookID, memberID)
					if returnErr != nil {
						fmt.Println("\t", returnErr)
					} else {

						fmt.Println("\t > Book returned successfully. <") // Placeholder message
					}
					break
				}

			case "5":
				fmt.Println("\n--- Available Books ---")
				availableBooks := library.ListAvailableBooks()
				if availableBooks == nil {
					fmt.Println("\tNo available books listed yet") // Placeholder message
				} else {
					showListOfBooks(availableBooks)
				}

			case "6":
				for {
					fmt.Print("Enter Member ID: ")
					memberIDStr, _ := reader.ReadString('\n')
					memberID, convertionErr := strconv.Atoi(strings.TrimSpace(memberIDStr)) // Implement conversion
					if convertionErr != nil {
						if returnToMainMenu() {
							break
						}
						continue
					}
					borrowedBooks := library.ListBorrowedBooks(memberID)
					fmt.Println("\n--- Borrowed Books by Member ---")
					if len(borrowedBooks) ==  0 {
						fmt.Println("\tNo borrowed books listed yet") // Placeholder message
					}
					showListOfBooks(borrowedBooks)
					break
				}

			case "7":
				fmt.Println("Exiting Library Management System. Goodbye!")
				return

			default:
				fmt.Println("\tInvalid choice. Please try again.")

		}
		fmt.Println()
		fmt.Println()
	}
}
