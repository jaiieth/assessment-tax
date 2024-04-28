package helper

import (
	"testing"
)

func TestErrorRes(t *testing.T) {
	act := "test"
	res := ErrorRes(act)

	if res.Message != act {
		t.Errorf("expected %s, got %s", act, res.Message)
	}
}
