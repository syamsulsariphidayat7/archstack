package cli

import "testing"

func TestIsUpToDate(t *testing.T) {
	cases := []struct {
		latest  string
		current string
		want    bool
	}{
		{"v0.3.0", "v0.3.0", true},
		{"v0.3.0", "0.3.0", true},
		{"v0.3.0", "v0.2.0", false},
		{"v1.0.0", "v0.9.9", false},
		{"v0.3.0", "dev", false},
		{"v0.3.0", "v0.3.0-rc1", false},
	}
	for _, c := range cases {
		if got := isUpToDate(c.latest, c.current); got != c.want {
			t.Errorf("isUpToDate(%q, %q) = %v, want %v", c.latest, c.current, got, c.want)
		}
	}
}

func TestGoarchSupported(t *testing.T) {
	arch, err := goarch()
	if err != nil {
		t.Fatalf("goarch() error: %v", err)
	}
	if arch != "amd64" && arch != "arm64" {
		t.Errorf("goarch() = %q, want amd64 atau arm64", arch)
	}
}
