package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

//======|| FlatStatus ||========================================

type FlatStatus string

const (
	StatusCreated      FlatStatus = "created"
	StatusApproved     FlatStatus = "approved"
	StatusDeclined     FlatStatus = "declined"
	StatusOnModeration FlatStatus = "on_moderation"
)

func ParseFlatStatus(str string) (FlatStatus, error) {
	var st FlatStatus

	switch str {
	case string(StatusCreated):
		st = StatusCreated
	case string(StatusApproved):
		st = StatusApproved
	case string(StatusDeclined):
		st = StatusDeclined
	case string(StatusOnModeration):
		st = StatusOnModeration
	default:
		return "", errors.New(fmt.Sprintf("unknown enum value %s", str))
	}

	return st, nil
}

func (st *FlatStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("faile type assertion")
	}

	status, err := ParseFlatStatus(str)
	if err != nil {
		return err
	}

	*st = status
	return nil
}

func (st FlatStatus) Value() (driver.Value, error) {
	return string(st), nil
}

//======|| UsertType ||========================================

type UserType string

const (
	Moderator UserType = "moderator"
	Client    UserType = "client"
)

func ParseUserType(str string) (UserType, error) {
	var ut UserType

	switch str {
	case string(Moderator):
		ut = Moderator
	case string(Client):
		ut = Client
	default:
		return "", errors.New(fmt.Sprintf("unknown enum value %s", str))
	}

	return ut, nil
}

func (ut *UserType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("faile type assertion")
	}

	status, err := ParseUserType(str)
	if err != nil {
		return err
	}

	*ut = status
	return nil
}

func (ut UserType) Value() (driver.Value, error) {
	return string(ut), nil
}
