package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	//panic("implement me")
	var parts []string
	for _, err := range v {
		parts = append(parts, err.Field)
	}
	return strings.Join(parts, "|||")
}

func Validate(v interface{}) error {
	var errs error

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		if rv.IsNil() {
			return errors.New("pointer is nil")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("value is not a struct")
	}

	return nil
}

type Validator[T any] interface {
	Validate(T) error
}

type funValidator[T any] func(T) error

func (f funValidator[T]) Validate(v T) error {
	return f(v)
}

func lenValidator(data int) Validator[string] {
	return funValidator[string](func(s string) error {
		if len(s) != data {
			return fmt.Errorf("length of %s is %d, want %d", s, len(s), data)
		}
		return nil
	})
}
