package concurrency

import (
	"errors"
	"sync"
	"time"
)

type ReservationRequest struct {
	BookID   int
	MemberID int
	Result   chan error
}

type ReservationWorker struct {
	requestQueue chan ReservationRequest
	reservations map[int]*ReservationInfo
	mutex        sync.RWMutex
	updateBook   func(int, string) error
	borrowBook   func(int, int) error
}

type ReservationInfo struct {
	BookID    int
	MemberID  int
	ReservedAt time.Time
	Timer     *time.Timer
}

func NewReservationWorker(updateBook func(int, string) error, borrowBook func(int, int) error) *ReservationWorker {
	worker := &ReservationWorker{
		requestQueue: make(chan ReservationRequest, 100),
		reservations: make(map[int]*ReservationInfo),
		updateBook:   updateBook,
		borrowBook:   borrowBook,
	}
	go worker.processReservations()
	return worker
}

func (rw *ReservationWorker) ReserveBook(bookID int, memberID int) error {
	rw.mutex.RLock()
	if existing, exists := rw.reservations[bookID]; exists {
		if time.Since(existing.ReservedAt) < 5*time.Second {
			rw.mutex.RUnlock()
			return errors.New("book is already reserved")
		}
	}
	rw.mutex.RUnlock()

	resultChan := make(chan error, 1)
	request := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Result:   resultChan,
	}

	select {
	case rw.requestQueue <- request:
		return <-resultChan
	case <-time.After(5 * time.Second):
		return errors.New("reservation request timeout")
	}
}

func (rw *ReservationWorker) processReservations() {
	for request := range rw.requestQueue {
		err := rw.handleReservation(request)
		request.Result <- err
	}
}

func (rw *ReservationWorker) handleReservation(request ReservationRequest) error {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	if existing, exists := rw.reservations[request.BookID]; exists {
		if time.Since(existing.ReservedAt) < 5*time.Second {
			return errors.New("book is already reserved")
		}
		existing.Timer.Stop()
		delete(rw.reservations, request.BookID)
	}

	err := rw.updateBook(request.BookID, "Reserved")
	if err != nil {
		return err
	}

	timer := time.AfterFunc(5*time.Second, func() {
		rw.cancelReservation(request.BookID)
	})

	rw.reservations[request.BookID] = &ReservationInfo{
		BookID:     request.BookID,
		MemberID:   request.MemberID,
		ReservedAt: time.Now(),
		Timer:      timer,
	}

	go rw.processBorrowing(request.BookID, request.MemberID)

	return nil
}

func (rw *ReservationWorker) processBorrowing(bookID int, memberID int) {
	time.Sleep(100 * time.Millisecond)

	rw.mutex.Lock()
	reservation, exists := rw.reservations[bookID]
	rw.mutex.Unlock()

	if !exists {
		return
	}

	if time.Since(reservation.ReservedAt) < 5*time.Second {
		err := rw.borrowBook(bookID, memberID)
		if err == nil {
			rw.mutex.Lock()
			if res, ok := rw.reservations[bookID]; ok {
				res.Timer.Stop()
				delete(rw.reservations, bookID)
			}
			rw.mutex.Unlock()
		}
	}
}

func (rw *ReservationWorker) cancelReservation(bookID int) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	if reservation, exists := rw.reservations[bookID]; exists {
		reservation.Timer.Stop()
		delete(rw.reservations, bookID)
		rw.updateBook(bookID, "Available")
	}
}

func (rw *ReservationWorker) IsReserved(bookID int) bool {
	rw.mutex.RLock()
	defer rw.mutex.RUnlock()
	_, exists := rw.reservations[bookID]
	return exists
}

