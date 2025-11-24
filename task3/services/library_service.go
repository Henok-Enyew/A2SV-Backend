package services

import (
	"errors"
	"sync"
	"task3/concurrency"
	"task3/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ReserveBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	books            map[int]models.Book
	members          map[int]models.Member
	mutex            sync.RWMutex
	reservationWorker *concurrency.ReservationWorker
}

func NewLibrary() *Library {
	library := &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
	
	library.reservationWorker = concurrency.NewReservationWorker(
		library.updateBookStatus,
		library.BorrowBook,
	)
	
	return library
}

func (l *Library) AddBook(book models.Book) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	book.Status = "Available"
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

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

func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mutex.RLock()
	book, exists := l.books[bookID]
	_, memberExists := l.members[memberID]
	l.mutex.RUnlock()

	if !exists {
		return errors.New("book not found")
	}

	if !memberExists {
		return errors.New("member not found")
	}

	if book.Status != "Available" {
		if book.Status == "Reserved" && l.reservationWorker.IsReserved(bookID) {
			return errors.New("book is already reserved")
		}
		return errors.New("book is not available")
	}

	return l.reservationWorker.ReserveBook(bookID, memberID)
}

func (l *Library) updateBookStatus(bookID int, status string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	book.Status = status
	l.books[bookID] = book
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

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
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var available []models.Book
	for _, book := range l.books {
		if book.Status == "Available" || book.Status == "Reserved" {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	member, exists := l.members[memberID]
	if !exists {
		return nil
	}
	return member.BorrowedBooks
}

func (l *Library) AddMember(member models.Member) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.members[member.ID] = member
}

func (l *Library) GetMember(memberID int) (models.Member, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	member, exists := l.members[memberID]
	return member, exists
}

func (l *Library) GetAllMembers() map[int]models.Member {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.members
}

