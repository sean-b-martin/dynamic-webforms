package validation

import (
	"errors"
)

var (
	ErrDatatypeDuplicate = errors.New("datatype already exists")
	ErrDatatypeNotFound  = errors.New("datatype not found")
	ErrElementNotFound   = errors.New("element not found")
	ErrWrongElementType  = errors.New("wrong element type")
)
