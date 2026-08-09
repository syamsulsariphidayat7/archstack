package system

import "testing"

func TestParseVersion(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "format pacman",
			in:   "Repository      : extra\nName            : php\nVersion         : 8.5.9-1\n",
			want: "8.5.9-1",
		},
		{
			name: "format yay aur",
			in:   "Repository      : aur\nName            : mongodb-bin\nVersion         : 8.0.4-1\n",
			want: "8.0.4-1",
		},
		{
			name: "spasi satu",
			in:   "Version : 1.2.3-4\n",
			want: "1.2.3-4",
		},
		{
			name: "tanpa baris version",
			in:   "Name : php\nRepository : extra\n",
			want: "",
		},
		{
			name: "input kosong",
			in:   "",
			want: "",
		},
		{
			name: "version tanpa nilai",
			in:   "Version :\n",
			want: "",
		},
	}
	for _, c := range cases {
		if got := parseVersion(c.in); got != c.want {
			t.Errorf("%s: parseVersion = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestParsePacmanQ(t *testing.T) {
	out := "nodejs 26.7.0-1\nphp 8.5.9-1\n"
	m := parsePacmanQ(out)
	if m["php"] != "8.5.9-1" || m["nodejs"] != "26.7.0-1" {
		t.Errorf("parsePacmanQ = %v", m)
	}
	if len(m) != 2 {
		t.Errorf("parsePacmanQ jumlah = %d, want 2", len(m))
	}
}

func TestParsePacmanQEmpty(t *testing.T) {
	if m := parsePacmanQ(""); len(m) != 0 {
		t.Errorf("parsePacmanQ kosong = %v", m)
	}
}

func TestParsePacmanSl(t *testing.T) {
	out := "extra nginx 1.30.4-1\nextra php 8.5.9-1 [installed]\ncore linux 6.10-1\n"
	m := parsePacmanSl(out)
	if m["nginx"] != "1.30.4-1" {
		t.Errorf("parsePacmanSl nginx = %q", m["nginx"])
	}
	if m["php"] != "8.5.9-1" {
		t.Errorf("parsePacmanSl php = %q", m["php"])
	}
	if m["linux"] != "6.10-1" {
		t.Errorf("parsePacmanSl linux = %q", m["linux"])
	}
}

func TestParsePacmanSlEmpty(t *testing.T) {
	if m := parsePacmanSl(""); len(m) != 0 {
		t.Errorf("parsePacmanSl kosong = %v", m)
	}
}
