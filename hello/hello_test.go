package hello

import "testing"

func TestYourFunction(t *testing.T) {
	got := YourFunction()
	want := "expected result"

	if got != want {
		t.Errorf("YourFunction() = %q, want %q", got, want)
	}
}
