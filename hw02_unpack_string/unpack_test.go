package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

//func TestUnpackInvalidString(t *testing.T) {
//	invalidStrings := []string{"3abc", "45", "aaa10b"}
//	for _, tc := range invalidStrings {
//
//		t.Run(tc, func(t *testing.T) {
//			_, err := Unpack(tc)
//		//	require.Falsef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
//
//		})
//	}
//}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {

		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Errorf(t, err, "actual error %q")

		})
	}
}

func TestIsValidArrRunesFailedByFirstSymbols(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "3aaa10b"}
	for _, str := range invalidStrings {
		var runes = []rune(str)
		t.Run(str, func(t *testing.T) {
			_, err := isValidArrRunes(runes, 0)
			require.Truef(t, errors.Is(err, ErrFirstSymbolIsNumber), "Must contain error %q", err)
		})
	}
}

func TestIsValidArrRunesFailedByPreviousSymbols(t *testing.T) {
	var str = "aa34bc"
	var runes = []rune(str)
	t.Run(str, func(t *testing.T) {
		_, err := isValidArrRunes(runes, 3)
		require.Truef(t, errors.Is(err, ErrPreviousInvalidItem), "Must contain error %q", err)
	})
}

func TestIsSymbolFailed(t *testing.T) {
	var runes []rune = []rune{'a', '\n'}

	for _, item := range runes {
		t.Run(string(item), func(t *testing.T) {
			result := isNumber(item)
			require.Falsef(t, result, "Must be false")
		})
	}

}

func TestIsSymbolSuccess(t *testing.T) {
	var runes []rune = []rune{'a', '\n'}

	for _, item := range runes {
		t.Run(string(item), func(t *testing.T) {
			result := isNumber(item)
			require.Falsef(t, result, "Must be success")
		})
	}

}
