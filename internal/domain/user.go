package domain

import (
	"fmt"
)

type User struct {
	id       int
	version  int
	fullName string
	login    string
}

func NewUser(
	id int,
	version int,
	fullName string,
	login string,
) User {
	return User{
		id:       id,
		version:  version,
		fullName: fullName,
		login:    login,
	}
}

func NewUninitializedUser(
	fullName string,
	login string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		login,
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

func (u User) Login() string {
	return u.login
}

func (u User) Validate() error {
	nameLen := len([]rune(u.fullName))
	if nameLen < 3 || nameLen > 100 {
		return fmt.Errorf(
			"invalid 'FullName' length %d: %w",
			nameLen,
			ErrInvalidArgument,
		)
	}

	loginLen := len([]rune(u.login))
	if loginLen < 3 || loginLen > 25 {
		return fmt.Errorf(
			"invalid 'Login' length %d: %w",
			loginLen,
			ErrInvalidArgument,
		)
	}

	return nil
}
