package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task3/models"
	"task3/services"
)

type LibraryController struct {
	library *services.Library
	scanner *bufio.Scanner
}

func NewLibraryController(library *services.Library) *LibraryController {
	return &LibraryController{
		library: library,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (lc *LibraryController) AddBook() {
	fmt.Print("Enter book ID: ")
	lc.scanner.Scan()
	bookID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	fmt.Print("Enter book title: ")
	lc.scanner.Scan()
	title := strings.TrimSpace(lc.scanner.Text())

	fmt.Print("Enter book author: ")
	lc.scanner.Scan()
	author := strings.TrimSpace(lc.scanner.Text())

	book := models.Book{
		ID:     bookID,
		Title:  title,
		Author: author,
	}

	lc.library.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (lc *LibraryController) RemoveBook() {
	fmt.Print("Enter book ID to remove: ")
	lc.scanner.Scan()
	bookID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	lc.library.RemoveBook(bookID)
	fmt.Println("Book removed successfully!")
}

func (lc *LibraryController) BorrowBook() {
	fmt.Print("Enter book ID: ")
	lc.scanner.Scan()
	bookID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	fmt.Print("Enter member ID: ")
	lc.scanner.Scan()
	memberID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	err := lc.library.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Book borrowed successfully!")
	}
}

func (lc *LibraryController) ReturnBook() {
	fmt.Print("Enter book ID: ")
	lc.scanner.Scan()
	bookID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	fmt.Print("Enter member ID: ")
	lc.scanner.Scan()
	memberID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	err := lc.library.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Book returned successfully!")
	}
}

func (lc *LibraryController) ListAvailableBooks() {
	books := lc.library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}

	fmt.Println("\nAvailable Books:")
	fmt.Println("ID\tTitle\t\tAuthor\t\tStatus")
	fmt.Println(strings.Repeat("-", 60))
	for _, book := range books {
		fmt.Printf("%d\t%s\t\t%s\t\t%s\n", book.ID, book.Title, book.Author, book.Status)
	}
}

func (lc *LibraryController) ListBorrowedBooks() {
	fmt.Print("Enter member ID: ")
	lc.scanner.Scan()
	memberID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	books := lc.library.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("No borrowed books for this member.")
		return
	}

	fmt.Printf("\nBorrowed Books for Member ID %d:\n", memberID)
	fmt.Println("ID\tTitle\t\tAuthor\t\tStatus")
	fmt.Println(strings.Repeat("-", 60))
	for _, book := range books {
		fmt.Printf("%d\t%s\t\t%s\t\t%s\n", book.ID, book.Title, book.Author, book.Status)
	}
}

func (lc *LibraryController) AddMember() {
	fmt.Print("Enter member ID: ")
	lc.scanner.Scan()
	memberID, _ := strconv.Atoi(strings.TrimSpace(lc.scanner.Text()))

	fmt.Print("Enter member name: ")
	lc.scanner.Scan()
	name := strings.TrimSpace(lc.scanner.Text())

	member := models.Member{
		ID:            memberID,
		Name:          name,
		BorrowedBooks: []models.Book{},
	}

	lc.library.AddMember(member)
	fmt.Println("Member added successfully!")
}

func (lc *LibraryController) ShowMenu() {
	fmt.Println("\n=== Library Management System ===")
	fmt.Println("1. Add Book")
	fmt.Println("2. Remove Book")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. List Available Books")
	fmt.Println("6. List Borrowed Books")
	fmt.Println("7. Add Member")
	fmt.Println("8. Exit")
	fmt.Print("Choose an option: ")
}

func (lc *LibraryController) Run() {
	for {
		lc.ShowMenu()
		lc.scanner.Scan()
		choice := strings.TrimSpace(lc.scanner.Text())

		switch choice {
		case "1":
			lc.AddBook()
		case "2":
			lc.RemoveBook()
		case "3":
			lc.BorrowBook()
		case "4":
			lc.ReturnBook()
		case "5":
			lc.ListAvailableBooks()
		case "6":
			lc.ListBorrowedBooks()
		case "7":
			lc.AddMember()
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

