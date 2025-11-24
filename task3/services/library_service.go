package services

import (
	"errors"
	"task3/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	book.Status = "Available"
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	member, exists := l.members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := l.members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("member has not borrowed this book")
	}

	book.Status = "Available"
	l.books[bookID] = book
	l.members[memberID] = member

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, book := range l.books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.members[memberID]
	if !exists {
		return nil
	}
	return member.BorrowedBooks
}

func (l *Library) AddMember(member models.Member) {
	l.members[member.ID] = member
}

func (l *Library) GetMember(memberID int) (models.Member, bool) {
	member, exists := l.members[memberID]
	return member, exists
}

func (l *Library) GetAllMembers() map[int]models.Member {
	return l.members
}

