package cli

import (
	"strings"
	"testing"
)

func TestSearchTools(t *testing.T) {
	cases := []struct {
		query string
		min   int
	}{
		{"php", 1},      // php, php-fpm, composer
		{"database", 1}, // postgres, mariadb, sqlite, mongodb
		{"data", 1},     // database-database di atas		{"node", 1},      // node, nodejs
		{"zzz-tidak-ada", 0},
	}
	for _, c := range cases {
		found := searchTools(c.query)
		if len(found) < c.min {
			t.Errorf("searchTools(%q) = %d hasil, want >= %d", c.query, len(found), c.min)
		}
		for _, tool := range found {
			hay := strings.ToLower(tool.Name + " " + tool.Pkg + " " + tool.Desc)
			if !strings.Contains(hay, c.query) {
				t.Errorf("searchTools(%q) memuat tool yang tidak cocok: %s", c.query, tool.Name)
			}
		}
	}
}

func TestSearchMatchesUppercaseDesc(t *testing.T) {
	// Query lowercase harus menemukan tool yang deskripsinya uppercase (mis. "PHP dependency manager").
	found := searchTools("php")
	hasComposer := false
	for _, tool := range found {
		if tool.Name == "composer" {
			hasComposer = true
		}
	}
	if !hasComposer {
		t.Error("searchTools(php) harus memuat composer (deskripsi 'PHP dependency manager')")
	}
}
