package mathutil

import "testing"

func TestSum(t *testing.T) {
	got := Sum([]int{1, 2, 3, 4})
	want := 10
	if got != want {
		// Report a clear failure message for CI logs.
		t.Errorf("Sum() = %d; want %d", got, want)
	}
}

func TestIsPrime(t *testing.T) {
	cases := map[int]bool{
		1:  false,
		2:  true,
		7:  true,
		9:  false,
		13: true,
	}
	for input, want := range cases {
		if got := IsPrime(input); got != want {
			t.Errorf("IsPrime(%d) = %v; want %v", input, got, want)
		}
	}
}
