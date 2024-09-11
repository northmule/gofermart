package util

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"password123", "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"},
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}

	for _, tc := range testCases {
		actual := PasswordHash(tc.input)
		if actual != tc.expected {
			t.Errorf("PasswordHash(%q) = %q; expected %q", tc.input, actual, tc.expected)
		}
	}
}
