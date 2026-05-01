package db_gorm

import "errors"

var (
	ErrCheckConstraintViolated = errors.New("check constraint violated")
	ErrDuplicatedKey           = errors.New("duplicated key")
	ErrRecordNotFound          = errors.New("record not found")
)
