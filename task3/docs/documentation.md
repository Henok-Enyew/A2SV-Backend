# Library Management System Documentation

## Overview

This is a console-based library management system implemented in Go with concurrent book reservation capabilities. The system allows users to manage books and members, handle book borrowing and returning operations, and supports concurrent book reservations using Goroutines, Channels, and Mutexes.

## Architecture

The system follows a layered architecture with clear separation of concerns:

- **Models**: Define data structures (Book, Member)
- **Services**: Contain business logic and implement the LibraryManager interface
- **Controllers**: Handle user input and coordinate with services
- **Concurrency**: Handles concurrent reservation requests using Goroutines and Channels
- **Main**: Entry point of the application

## Components

### Models

#### Book
Represents a book in the library with the following fields:
- `ID` (int): Unique identifier for the book
- `Title` (string): Title of the book
- `Author` (string): Author of the book
- `Status` (string): Current status - "Available", "Reserved", or "Borrowed"

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
- `ReserveBook(bookID int, memberID int) error`: Reserves a book for a member (new)
- `ReturnBook(bookID int, memberID int) error`: Allows a member to return a borrowed book
- `ListAvailableBooks() []Book`: Returns all available books
- `ListBorrowedBooks(memberID int) []Book`: Returns all books borrowed by a specific member

#### Library
Implements the LibraryManager interface and manages:
- Books stored in a map with book ID as key
- Members stored in a map with member ID as key
- Mutex (sync.RWMutex) for thread-safe concurrent access
- ReservationWorker for handling concurrent reservation requests

### Controllers

#### LibraryController
Handles console interactions and provides methods for:
- Adding books and members
- Removing books
- Borrowing and returning books
- Reserving books (with concurrent support)
- Listing available and borrowed books
- Displaying menu and processing user choices

### Concurrency

#### ReservationWorker
Handles concurrent book reservations using:
- **Goroutines**: Processes multiple reservation requests simultaneously
- **Channels**: Queues incoming reservation requests (buffered channel with capacity 100)
- **Mutex (sync.RWMutex)**: Prevents race conditions when updating reservation state
- **Timers**: Automatically cancels reservations after 5 seconds if not borrowed

The worker implements the following flow:
1. Reservation requests are queued through a channel
2. A dedicated Goroutine processes requests from the queue
3. Each reservation is tracked with a timer
4. If a book is not borrowed within 5 seconds, the reservation is automatically cancelled
5. After reservation, a separate Goroutine attempts to borrow the book asynchronously

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

4. **Concurrent Reservation System**
   - Reserve books with automatic borrowing
   - Handle multiple concurrent reservation requests
   - Auto-cancellation of reservations after 5 seconds
   - Thread-safe operations using Mutexes
   - Prevents double reservations

5. **Error Handling**
   - Validates book existence before operations
   - Validates member existence before operations
   - Prevents borrowing already borrowed books
   - Prevents reserving already reserved books
   - Validates return operations
   - Handles reservation timeouts

## Concurrency Approach

The system implements concurrent book reservations using the following Go concurrency primitives:

### 1. Goroutines
- **Reservation Processing**: A dedicated Goroutine (`processReservations`) continuously processes reservation requests from the queue
- **Borrowing Processing**: Each reservation spawns a Goroutine to handle the asynchronous borrowing process
- **Timer Management**: Each reservation has its own timer Goroutine for auto-cancellation

### 2. Channels
- **Request Queue**: A buffered channel (`requestQueue`) with capacity 100 queues incoming reservation requests
- **Result Communication**: Each reservation request includes a result channel to communicate success/failure back to the caller
- **Non-blocking Operations**: The channel-based design allows multiple requests to be queued without blocking

### 3. Mutexes (sync.RWMutex)
- **Library Operations**: All book and member map operations are protected by RWMutex for thread-safe concurrent access
- **Reservation State**: The reservation map is protected by RWMutex to prevent race conditions
- **Read-Write Optimization**: Uses RLock for read operations and Lock for write operations to maximize concurrency

### 4. Timer-based Auto-Cancellation
- Each reservation creates a timer that fires after 5 seconds
- If the timer fires before the book is borrowed, the reservation is automatically cancelled
- The book status is reverted to "Available" if the reservation expires

### 5. Concurrent Request Handling
- Multiple members can attempt to reserve books simultaneously
- The system ensures only one reservation per book at a time
- Requests are processed in order through the channel queue
- Race conditions are prevented through proper mutex usage

## Usage

1. Run the application: `go run main.go`
2. Follow the menu prompts to perform operations
3. Add members before allowing them to borrow or reserve books
4. Add books before they can be borrowed or reserved
5. Use "Reserve Book" option to reserve a book (automatically borrowed after reservation, expires in 5 seconds)

## Error Handling

The system includes comprehensive error handling for:
- Book not found scenarios
- Member not found scenarios
- Attempting to borrow already borrowed books
- Attempting to reserve already reserved books
- Attempting to return books not borrowed by the member
- Reservation request timeouts
- Concurrent access conflicts

## Thread Safety

All critical sections are protected using Mutexes:
- Book map operations (Add, Remove, Update)
- Member map operations (Add, Get)
- Reservation state management
- Status updates

This ensures that multiple Goroutines can safely access shared resources without data races or corruption.

## Data Storage

Currently, the system uses in-memory storage (maps). All data is lost when the program terminates. For persistent storage, a database integration would be required.

