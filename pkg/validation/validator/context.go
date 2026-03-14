package validator

import (
	"errors"
)

var ErrIndexOutOfBounds = errors.New("index out of bounds")

type FieldRef struct {
	Name  string
	Index int
}

type Context struct {
	Values []any
}

func (ctx *Context) Get(ref *FieldRef) (any, error) {
	if ref.Index < 0 || ref.Index >= len(ctx.Values) {
		return nil, ErrIndexOutOfBounds
	}
	return ctx.Values[ref.Index], nil
}

func (ctx *Context) ForEachValue(ref *FieldRef, fn func(val any) error) (*ValidationError, error) {
	field, err := ctx.Get(ref)
	if err != nil {
		return nil, err
	}

	valErr := ValidationError{}
	if value, ok := field.([]any); ok {
		for i, v := range value {
			if err := fn(v); err != nil {
				valErr.AddFieldError(FieldError{
					Name:  ref.Name,
					Index: i,
					Msg:   err.Error(),
				})
			}
		}
	} else {
		// scalar value
		if err := fn(field); err != nil {
			valErr.AddFieldError(FieldError{
				Name:  ref.Name,
				Index: -1,
				Msg:   err.Error(),
			})
		}
	}

	if !valErr.Empty() {
		return &valErr, nil
	}

	return nil, nil
}
