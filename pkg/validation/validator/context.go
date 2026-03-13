package validator

import "errors"

var ErrIndexOutOfBounds = errors.New("index out of bounds")

type FieldRef struct {
	Name  string
	Index int
}

type Context struct {
	Values []any
}

func (ctx *Context) Iter(ref *FieldRef) ([]any, error) {
	if ref.Index < 0 || ref.Index >= len(ctx.Values) {
		return nil, ErrIndexOutOfBounds
	}

	field := ctx.Values[ref.Index]
	if val, ok := field.([]any); ok {
		return val, nil
	}
	return []any{field}, nil
}
