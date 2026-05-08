package validation

import "testing"

func TestSlug(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"golang", true},
		{"redis-lua", true},
		{"go-1-22", true},
		{"", false},
		{"-lead", false},
		{"trail-", false},
		{"UPPER", false},
		{"spa ce", false},
		{"under_score", false},
		{"double--dash", false},
	}
	for _, c := range cases {
		got := slugRe.MatchString(c.in)
		if got != c.want {
			t.Errorf("slug %q: want %v, got %v", c.in, c.want, got)
		}
	}
}
