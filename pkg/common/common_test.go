package common

import "testing"

func TestDummy(t *testing.T) {
	out := DummyFunction()
	if out != "OK" {
		t.Errorf("Expected 'OK', but returned '%s'", out)
	}
}
