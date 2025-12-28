package storage

import "errors"

var (
	ErrTaskNotFound  = errors.New("task is not found")
	ErrWrongArgument = errors.New("wrong argument")
)
