package model

import "errors"

// Квартира
type Flat struct {
	ID      int64      `json:"id" db:"id"`
	HouseID int64      `json:"house_id" db:"house_id"`
	Price   int64      `json:"price" db:"price"`
	Rooms   int64      `json:"rooms" db:"rooms"`
	Status  FlatStatus `json:"status" db:"status"`
}

var (
	ErrImpossibleTransition = errors.New("impossible status transition")
)

func (f *Flat) Approve() error {
	if f.Status != StatusOnModeration {
		return ErrImpossibleTransition
	}
	f.Status = StatusApproved
	return nil
}

func (f *Flat) Decline() error {
	if f.Status != StatusOnModeration {
		return ErrImpossibleTransition
	}
	f.Status = StatusDeclined
	return nil
}

func (f *Flat) StartModeration() error {
	if f.Status != StatusCreated {
		return ErrImpossibleTransition
	}
	f.Status = StatusOnModeration
	return nil
}
