package utils

import "testing"

func TestContainsString(t *testing.T) {
	list := []string{"foo", "bar", "foobar"}

	testCases := []struct {
		match    string
		expected bool
	}{
		{"foo", true},
		{"bar", true},
		{"xxx", false},
		{"foobar", true},
		{"foob", false},
	}

	for _, tc := range testCases {
		t.Run(tc.match, func(t *testing.T) {
			actual := ContainsString(list, tc.match)
			if actual != tc.expected {
				t.Errorf("Failed test %s\n", tc.match)
			}
		})
	}
}
