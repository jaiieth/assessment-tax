package helper

import (
	"testing"
)

func TestRoundTwoDigits(t *testing.T) {
	act := 100.232123
	res := RoundTwoDigits(act)

	expect := 100.23
	if res != expect {
		t.Errorf("expected %v, got %v", expect, res)
	}
}
