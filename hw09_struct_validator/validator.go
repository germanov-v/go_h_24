package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationErrors []ValidationError

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationErrors) Error() string {
	//panic("implement me")
	var parts []string
	for _, err := range v {
		parts = append(parts, err.Field)
	}
	return strings.Join(parts, "|||")
}

func Validate(v interface{}) error {
	var errorsResult ValidationErrors

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

	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// embed поле
		//if field.Anonymous {
		//
		//}

		// public
		if !field.IsExported() {
			continue
		}

		//tag := strings.Split(field.Tag.Get("json"), ",")

		tag := field.Tag.Get("validate") // empty or nil ?
		if tag == "" {
			continue
		}

		valField := rv.Field(i)
		nameField := field.Name

		if tag == "nested" {
			internErr := Validate(valField.Interface())
			if internErr != nil {
				var internErrs ValidationErrors
				// вложенность
				if errors.As(internErr, &internErrs) {
					for _, err := range internErrs {
						errorsResult = append(errorsResult, ValidationError{
							Field: fmt.Sprintf("%s.%s", field.Name, nameField),
							Err:   err.Err,
						})
					}
				} else {
					errorsResult = append(errorsResult, ValidationError{
						Field: fmt.Sprintf("%s.%s", field.Name, nameField),
						Err:   internErr,
					})
				}
			}
			continue
			//return nil
		}

		switch valField.Kind() {
		case reflect.Struct:
			// парсинг
			validators, err := createStrRuleValidators(tag)
			if err != nil {
				return fmt.Errorf("%s: %w", nameField, err)
			}
			s := valField.String()
			for _, valid := range validators {
				if err := valid.Validate(s); err != nil {
					errorsResult = append(errorsResult, ValidationError{
						Field: fmt.Sprintf("%s.%s", field.Name, nameField),
						Err:   err,
					})
				}
			}
		case reflect.Slice:
			for i := 0; i < valField.Len(); i++ {
				item := valField.Index(i)
				//fmt.Print
				//itemName := fmt.Sprintf("%s.%s", nameField, i)
				itemName := fmt.Sprintf("%s.%d", nameField, i)
				switch item.Kind() {
				//case reflect.Struct:
				case reflect.String:
					// todo: в отдельную - копипаст пошел
					validators, err := createStrRuleValidators(tag)
					if err != nil {
						return fmt.Errorf("%s: %w", itemName, err)
					}
					s := valField.String()
					for _, valid := range validators {
						if err := valid.Validate(s); err != nil {
							errorsResult = append(errorsResult, ValidationError{
								Field: fmt.Sprintf("%s.%s", field.Name, nameField),
								Err:   err,
							})
						}
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					validators, err := createIntRuleValidators(tag)
					if err != nil {
						return fmt.Errorf("%s: %w", itemName, err)
					}
					j := int(item.Int()) // обрезать?
					for _, valid := range validators {
						if err := valid.Validate(j); err != nil {
							errorsResult = append(errorsResult, ValidationError{
								Field: fmt.Sprintf("%s.%s", itemName, nameField),
								Err:   err,
							})
						}
					}
				default:
					//break;
					continue
				}
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			validators, err := createIntRuleValidators(tag)
			if err != nil {
				return fmt.Errorf("%s: %w", nameField, err)
			}
			j := int(valField.Int())
			for _, valid := range validators {
				if err := valid.Validate(j); err != nil {
					errorsResult = append(errorsResult, ValidationError{
						Field: fmt.Sprintf("%s.%s", nameField, nameField),
						Err:   err,
					})
				}
			}

		default:
			//return fmt.Errorf("%s is not a compatible type", valField.Kind())
			continue
		}

		//continue
		//return nil
	}

	if len(errorsResult) > 0 {
		return errorsResult
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
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}

		switch parts[0] {
		case "len":
			expected, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid tag rule: %q", rule)
			}
			validators = append(validators, lenValidator(expected))
		case "regexp":
			validators = append(validators, regValidator(parts[1]))

		case "in":
			options := strings.Split(parts[1], ",")
			var items = make([]string, len(options))
			for _, option := range options {
				//temp, err := strconv.Atoi(option)
				//
				//if err != nil {
				//	return nil, fmt.Errorf("invalid tag rule: %q - %w ", rule, err)
				//}
				items = append(items, strings.TrimSpace(option))
			}
			validators = append(validators, containStrValidator(items))
		default:
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}
	}
	return validators, nil
}

func createIntRuleValidators(tag string) ([]Validator[int], error) {
	var validators []Validator[int]

	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}

		switch parts[0] {
		case "min":
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid tag rule: %q", rule)
			}
			validators = append(validators, minValidator(val))
		case "max":
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid tag rule: %q", rule)
			}
			validators = append(validators, minValidator(val))
		case "in":
			options := strings.Split(parts[1], ",")
			items := make([]int, len(options))
			for _, option := range options {
				s := strings.TrimSpace(option)
				val, err := strconv.Atoi(s)
				if err != nil {
					return nil, fmt.Errorf("invalid tag rule: %q - %w", rule, err)
				}
				items = append(items, val)
			}
			validators = append(validators, containInValidator(items))
		default:
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}
	}

	return validators, nil
}
