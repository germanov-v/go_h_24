package hw09structvalidator

import (
	"fmt"
	"regexp"
)

func LenValidator(data int) Validator[string] {
	return FunValidator[string](func(s string) error {
		if len(s) != data {
			//return errors.New("len fail")
			return fmt.Errorf("length of %s is %d, want %d", s, len(s), data)
		}
		return nil
	})
}

func RegValidator(pat string) Validator[string] {
	r, err := regexp.Compile(pat)
	if err != nil {
		return FunValidator[string](func(s string) error {
			return fmt.Errorf("invalid regular expression: %s - %w", s, err)
		})
	}
	return FunValidator[string](func(s string) error {

		if !r.MatchString(s) {
			// с кавычками пробуем формат %q
			return fmt.Errorf("invalid regular expression: %q  - %w - %q", s, err, pat)
		}

		return nil
	})
}

func ContainStrValidator(values []string) Validator[string] {
	return FunValidator[string](func(s string) error {
		for _, v := range values {
			if v == s {
				return nil
			}
		}
		return fmt.Errorf("%q not in %v", s, values)
	})
}

func MinValidator(min int) Validator[int] {
	return FunValidator[int](func(v int) error {
		if v < min {
			return fmt.Errorf("value %d is less than minimum %d", v, min)
		}
		return nil
	})
}

func MaxValidator(max int) Validator[int] {
	return FunValidator[int](func(v int) error {
		if v > max {
			return fmt.Errorf("value %d is greater than maximum %d", v, max)
		}
		return nil
	})
}
