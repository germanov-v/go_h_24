package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// valid
		{
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012", // l 36
				Name:   "Alice",
				Age:    30, //  18 <x < 50
				Email:  "test@example.com",
				Role:   "admin",                                // in {admin, stuff}
				Phones: []string{"12345678901", "09876543210"}, // len(11)
			},
			expectedErr: nil,
		},
		// fail
		{
			in: User{
				ID:     "short-id", // // !len(36)
				Name:   "Bob",
				Age:    17,                                  // <18
				Email:  "invalid",                           // fail regexp
				Role:   "user",                              // no in  {admin, stuff}
				Phones: []string{"1234567890", "123456789"}, // !len(11)
			},
			expectedErr: errors.New("verror validation"),
		},
		// valid App.
		{
			in: App{
				Version: "1.234", // len(5)
			},
			expectedErr: nil,
		},
		// fail
		{
			in: App{
				Version: "12.341", // !len(5)
			},
			expectedErr: errors.New("verror validation"),
		},
		// success
		{
			in: Token{
				Header:    []byte("Header"),
				Payload:   []byte("Payload"),
				Signature: []byte("Signature"),
			},
			expectedErr: nil,
		},
		// success
		{
			in: Response{
				Code: 200, // in {200,404,500}
				Body: "OK",
			},
			expectedErr: nil,
		},
		// fail
		{
			in: Response{
				Code: 201,
				Body: "Created",
			},
			expectedErr: errors.New("verror validation"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			// Place your code here.

			err := Validate(tt.in)

			// отсутсвие ошибок
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("err unexpected: %v, got: %v", tt.expectedErr, err)
				}
				return
			}

			// ожидаем ошибку или ошибки?!
			if err == nil {
				t.Errorf("err expected: %v, got: nil", tt.expectedErr)
			}

			// базовая проверка на тип
			var vErr ValidationErrors
			if !errors.As(err, &vErr) || len(vErr) == 0 {
				t.Errorf("err expected: %v, got: %v", tt.expectedErr, err)
			}

			_ = tt
		})
	}
}
