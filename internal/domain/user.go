package domain

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
	valErrors := make(map[string]string)

	nameLen := len([]rune(u.fullName))
	if nameLen < 3 || nameLen > 100 {
		valErrors["full_name"] = "full_name must be between 3 and 100 symbols"
	}

	loginLen := len([]rune(u.login))
	if loginLen < 3 || loginLen > 25 {
		valErrors["login"] = "login must be between 3 and 25 symbols"
	}

	if len(valErrors) != 0 {
		return NewValidationError(valErrors)
	}

	return nil
}
