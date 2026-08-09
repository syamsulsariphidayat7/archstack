package registry

import "testing"

func TestAllToolsNotEmpty(t *testing.T) {
	tools := AllTools()
	if len(tools) == 0 {
		t.Fatal("AllTools() mengembalikan daftar kosong")
	}
}

func TestAllToolsDataValid(t *testing.T) {
	seen := make(map[string]bool)
	for _, tool := range AllTools() {
		if tool.Name == "" {
			t.Error("ada tool dengan Name kosong")
		}
		if seen[tool.Name] {
			t.Errorf("nama tool duplikat: %s", tool.Name)
		}
		seen[tool.Name] = true

		if tool.Pkg == "" {
			t.Errorf("%s: Pkg kosong", tool.Name)
		}
		if tool.Desc == "" {
			t.Errorf("%s: Desc kosong", tool.Name)
		}
		if tool.From != SourcePacman && tool.From != SourceYay {
			t.Errorf("%s: Source tidak valid: %q", tool.Name, tool.From)
		}
		if tool.Binary == "" && tool.Service == "" {
			t.Errorf("%s: tidak punya Binary maupun Service", tool.Name)
		}
	}
}

func TestGetTool(t *testing.T) {
	tool, ok := GetTool("php")
	if !ok {
		t.Fatal("GetTool(php) harus ditemukan")
	}
	if tool.Name != "php" || tool.Pkg != "php" || tool.From != SourcePacman {
		t.Errorf("GetTool(php) data salah: %+v", tool)
	}

	if _, ok := GetTool("tidak-ada"); ok {
		t.Error("GetTool untuk tool tak dikenal harus mengembalikan ok=false")
	}
}
