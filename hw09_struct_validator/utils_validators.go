package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

type FunValidator[T any] func(T) error //result,error

func (f FunValidator[T]) Validate(v T) error {
	return f(v)
}

func ContainInValidator(values []int) Validator[int] {
	return FunValidator[int](func(val int) error {
		for _, v := range values {
			if v == val {
				return nil
			}
		}
		return fmt.Errorf("value %d not in %v", val, values)
	})
}

func CreateStrRuleValidators(tag string) ([]Validator[string], error) {
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
			validators = append(validators, LenValidator(expected))
		case "regexp":
			validators = append(validators, RegValidator(parts[1]))

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
			validators = append(validators, ContainStrValidator(items))
		default:
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}
	}
	return validators, nil
}

func CreateIntRuleValidators(tag string) ([]Validator[int], error) {
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
			validators = append(validators, MinValidator(val))
		case "max":
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid tag rule: %q", rule)
			}
			validators = append(validators, MaxValidator(val))
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
			validators = append(validators, ContainInValidator(items))
		default:
			return nil, fmt.Errorf("invalid tag rule: %q", rule)
		}
	}

	return validators, nil
}
