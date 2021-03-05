package util_test

import (
	"project/app/pkg/util"
	"testing"
	"unicode"
)

func TestGenerateRandomDigits(t *testing.T) {
	width := 6
	str := util.GenerateRandomDigits(6)
	if len(str) != width {
		t.Fail()
	}
	for _, r := range []rune(str) {
		if !unicode.IsDigit(r) {
			t.Fail()
		}
	}
}
