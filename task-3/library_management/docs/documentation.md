# Library Management System (Go Console Application)

---

## Project Objective

This project aims to develop a **simple console-based library management system** using the Go programming language. Its primary goal is to showcase fundamental Go concepts such as **structs, interfaces, methods, slices, and maps** within a practical application.

---

## Requirements

### Structs

The system uses two main data structures to represent its core entities:

- **`Book` Struct**
    
    - `ID` (int): A unique number for each book.
        
    - `Title` (string): The book's title.
        
    - `Author` (string): The book's author.
        
    - `Status` (string): Shows if the book is `"Available"` or `"Borrowed"`.
        
- **`Member` Struct**
    
    - `ID` (int): A unique number for each library member.
        
    - `Name` (string): The member's name.
        
    - `BorrowedBooks` ([]Book): A list of `Book` structs the member has currently borrowed.
        

### Interfaces

A **`LibraryManager` interface** is defined. This acts like a blueprint, outlining all the operations a library manager should be able to perform, which helps keep the code organized.

- **`LibraryManager` Interface Methods**
    
    - `AddBook(book Book)`: Adds a new book to the library.
        
    - `RemoveBook(bookID int)`: Removes a book by its ID.
        
    - `BorrowBook(bookID int, memberID int) error`: Lets a member borrow a book. It will report an error if something goes wrong (e.g., the book isn't available).
        
    - `ReturnBook(bookID int, memberID int) error`: Lets a member return a book. It will report an error if something goes wrong.
        
    - `ListAvailableBooks() []Book`: Shows all books that are currently available.
        
    - `ListBorrowedBooks(memberID int) []Book`: Shows all books a specific member has borrowed.
        

---

## Implementation

The `LibraryManager` interface is brought to life by a **`Library` struct**. This `Library` struct is where all the books and members are stored.

- **`Library` Struct Storage**
    
    - A **map** stores all books, using the book's ID as a quick lookup key (e.g., `map[int]Book`).
        
    - Another **map** stores all members, using the member's ID as a quick lookup key (e.g., `map[int]Member`).
        

### Methods

The `Library` struct provides the actual code for all the methods listed in the `LibraryManager` interface:

- **`AddBook`**: Adds a new `Book` to the library's collection.
    
- **`RemoveBook`**: Deletes a book from the library using its ID.
    
- **`BorrowBook`**:
    
    - Checks if the book is free and if the member exists.
        
    - Changes the book's status to "Borrowed."
        
    - Adds the book to the member's list of borrowed books.
        
- **`ReturnBook`**:
    
    - Confirms the member actually borrowed the book.
        
    - Changes the book's status back to "Available."
        
    - Removes the book from the member's list of borrowed books.
        
- **`ListAvailableBooks`**: Goes through all the books and gives back a list of only those marked "Available."
    
- **`ListBorrowedBooks`**: Finds a specific member and returns their list of borrowed books.
    

---

## Console Interaction

The system features a **simple command-line interface** that lets users interact with the library. Users choose options from a menu to perform actions.

- **Key Console Functions**
    
    - Add a new book.
        
    - Remove an existing book.
        
    - Borrow a book.
        
    - Return a book.
        
    - List all available books.
        
    - List all books borrowed by a specific member.
        

Additionally, in the `library_controller.go` file, I added a new choice to the main menu. This new option lets people who use the system a lot (like librarians) **see all the books and members easily**. It's like a quick way for them to check everything in the library.