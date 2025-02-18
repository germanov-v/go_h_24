package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

type funValidator[T any] func(T) error //result,error

func (f funValidator[T]) Validate(v T) error {
	return f(v)
}

func lenValidator(data int) Validator[string] {
	return funValidator[string](func(s string) error {
		if len(s) != data {
			//return errors.New("len fail")
			return fmt.Errorf("length of %s is %d, want %d", s, len(s), data)
		}
		return nil
	})
}

func regValidator(pat string) Validator[string] {
	r, err := regexp.Compile(pat)
	if err != nil {
		return funValidator[string](func(s string) error {
			return fmt.Errorf("invalid regular expression: %s - %w", s, err)
		})
	}
	return funValidator[string](func(s string) error {

		if !r.MatchString(s) {
			// с кавычками пробуем формат %q
			return fmt.Errorf("invalid regular expression: %q  - %w - %q", s, err, pat)
		}

		return nil
	})
}

func containStrValidator(values []string) Validator[string] {
	return funValidator[string](func(s string) error {
		for _, v := range values {
			if v == s {
				return nil
			}
		}
		return fmt.Errorf("%q not in %v", s, values)
	})
}

func minValidator(min int) Validator[int] {
	return funValidator[int](func(v int) error {
		if v < min {
			return fmt.Errorf("value %d is less than minimum %d", v, min)
		}
		return nil
	})
}

func maxValidator(max int) Validator[int] {
	return funValidator[int](func(v int) error {
		if v > max {
			return fmt.Errorf("value %d is greater than maximum %d", v, max)
		}
		return nil
	})
}

func containInValidator(values []int) Validator[int] {
	return funValidator[int](func(val int) error {
		for _, v := range values {
			if v == val {
				return nil
			}
		}
		return fmt.Errorf("value %d not in %v", val, values)
	})
}

func createStrRuleValidators(tag string) ([]Validator[string], error) {
	var validators []Validator[string]
	rules := strings.Split(tag, "|")
	for index, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule==""{
			continue
		}
		parts:=strings.SplitN(rule, ":",2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}

		switch parts[0] {
		case "len":
          expected, err := strconv.Atoi(parts[1])
		case "regexp":

		case "in"
		}
	}
	return validators, nil
}

func createIntRuleValidators(tag string) ([]Validator[int], error) {
	var validators []Validator[int]

	return validators, nil
}
