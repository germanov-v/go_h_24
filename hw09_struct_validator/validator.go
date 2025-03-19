package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
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
	if rv.Kind() == reflect.Ptr {
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

		case reflect.String:
			validators, err := CreateStrRuleValidators(tag)
			if err != nil {
				return fmt.Errorf("%s: %w", nameField, err)
			}
			s := valField.String() //valField.String()
			for _, valid := range validators {
				if err := valid.Validate(s); err != nil {
					errorsResult = append(errorsResult, ValidationError{
						Field: fmt.Sprintf("%s.%s", field.Name, nameField),
						Err:   err,
					})
				}
			}
		case reflect.Struct:
			// парсинг
			validators, err := CreateStrRuleValidators(tag)
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
					validators, err := CreateStrRuleValidators(tag)
					if err != nil {
						return fmt.Errorf("%s: %w", itemName, err)
					}
					s := item.String() //valField.String()
					for _, valid := range validators {
						if err := valid.Validate(s); err != nil {
							errorsResult = append(errorsResult, ValidationError{
								Field: fmt.Sprintf("%s.%s", field.Name, nameField),
								Err:   err,
							})
						}
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					validators, err := CreateIntRuleValidators(tag)
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
			validators, err := CreateIntRuleValidators(tag)
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
