package main

import (
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	const s, sep, want = "chicken", "ken", 4
	got := strings.Index(s, sep)
	if got != want {
		t.Errorf("Index(%q,%q) = %v; want %v", s, sep, got, want) //注意原slide中的got和want写反了
	}
}

func TestIndex1(t *testing.T) {
	var tests = []struct {
		s   string
		sep string
		out int
	}{
		{"", "", 0},
		{"", "a", -1},
		{"fo", "foo", -1},
		{"foo", "foo", 0},
		{"oofofoofooo", "f", 2},
		// etc
	}
	for _, test := range tests {
		actual := strings.Index(test.s, test.sep)
		if actual != test.out {
			t.Errorf("Index(%q,%q) = %v; want %v",
				test.s, test.sep, actual, test.out)
		}
	}
}
