package domain

import (
	"fmt"
)

type User struct {
	id       int
	version  int
	fullName string
}

const (
	UninitializedID      = -1
	UninitializedVersion = -1
)

func NewUser(
	id int,
	version int,
	fullName string,
) User {
	return User{
		id:       id,
		version:  version,
		fullName: fullName,
	}
}

func NewUninitializedUser(
	fullName string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
	)
}

func (u User) ID() int {
	return u.id
}

func (u User) Version() int {
	return u.version
}

func (u User) FullName() string {
	return u.fullName
}

func (u User) Validate() error {
	nameLen := len([]byte(u.fullName))
	if nameLen < 3 || nameLen > 20 {
		return fmt.Errorf(
			"validating user: %s: %w",
			"invalid 'FullName' length",
			ErrInvalidArgument,
		)
	}

	return nil
}
