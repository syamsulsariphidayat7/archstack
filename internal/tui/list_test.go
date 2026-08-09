package tui

import "testing"

func TestPadRight(t *testing.T) {
	cases := []struct {
		s    string
		n    int
		want string
	}{
		{"abc", 5, "  "},
		{"abc", 3, ""},
		{"abc", 2, ""},
		{"", 3, "   "},
		{"abcdef", 4, ""},
	}
	for _, c := range cases {
		if got := padRight(c.s, c.n); got != c.want {
			t.Errorf("padRight(%q, %d) = %q, want %q", c.s, c.n, got, c.want)
		}
	}
}
