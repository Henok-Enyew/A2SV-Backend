# Library Management System Documentation

## Overview

This is a console-based library management system implemented in Go. The system allows users to manage books and members, handle book borrowing and returning operations, and view available and borrowed books.

## Architecture

The system follows a layered architecture with clear separation of concerns:

- **Models**: Define data structures (Book, Member)
- **Services**: Contain business logic and implement the LibraryManager interface
- **Controllers**: Handle user input and coordinate with services
- **Main**: Entry point of the application

## Components

### Models

#### Book
Represents a book in the library with the following fields:
- `ID` (int): Unique identifier for the book
- `Title` (string): Title of the book
- `Author` (string): Author of the book
- `Status` (string): Current status - "Available" or "Borrowed"

#### Member
Represents a library member with the following fields:
- `ID` (int): Unique identifier for the member
- `Name` (string): Name of the member
- `BorrowedBooks` ([]Book): Slice containing all books currently borrowed by the member

### Services

#### LibraryManager Interface
Defines the contract for library management operations:
- `AddBook(book Book)`: Adds a new book to the library
- `RemoveBook(bookID int)`: Removes a book from the library
- `BorrowBook(bookID int, memberID int) error`: Allows a member to borrow a book
- `ReturnBook(bookID int, memberID int) error`: Allows a member to return a borrowed book
- `ListAvailableBooks() []Book`: Returns all available books
- `ListBorrowedBooks(memberID int) []Book`: Returns all books borrowed by a specific member

#### Library
Implements the LibraryManager interface and manages:
- Books stored in a map with book ID as key
- Members stored in a map with member ID as key

### Controllers

#### LibraryController
Handles console interactions and provides methods for:
- Adding books and members
- Removing books
- Borrowing and returning books
- Listing available and borrowed books
- Displaying menu and processing user choices

## Features

1. **Book Management**
   - Add new books to the library
   - Remove books from the library
   - View all available books

2. **Member Management**
   - Add new members to the system
   - Track borrowed books per member

3. **Borrowing System**
   - Borrow books (with validation)
   - Return borrowed books
   - View borrowed books by member

4. **Error Handling**
   - Validates book existence before operations
   - Validates member existence before operations
   - Prevents borrowing already borrowed books
   - Validates return operations

## Usage

1. Run the application: `go run main.go`
2. Follow the menu prompts to perform operations
3. Add members before allowing them to borrow books
4. Add books before they can be borrowed

## Error Handling

The system includes comprehensive error handling for:
- Book not found scenarios
- Member not found scenarios
- Attempting to borrow already borrowed books
- Attempting to return books not borrowed by the member

## Data Storage

Currently, the system uses in-memory storage (maps). All data is lost when the program terminates. For persistent storage, a database integration would be required.

