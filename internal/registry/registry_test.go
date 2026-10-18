package registry

import (
	"strings"
	"testing"
)

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

		if tool.Pkg == "" && tool.From != SourceServiceOnly {
			t.Errorf("%s: Pkg kosong", tool.Name)
		}
		if tool.Desc == "" {
			t.Errorf("%s: Desc kosong", tool.Name)
		}
		switch tool.From {
		case SourcePacman, SourceYay:
			if tool.Pkg == "" {
				t.Errorf("%s: Pkg kosong", tool.Name)
			}
		case SourceServiceOnly:
			if tool.Pkg != "" {
				t.Errorf("%s: tool service-only tidak boleh punya Pkg (dapat %q)", tool.Name, tool.Pkg)
			}
			if tool.Service == "" {
				t.Errorf("%s: tool service-only wajib punya Service", tool.Name)
			}
		default:
			t.Errorf("%s: Source tidak valid: %q", tool.Name, tool.From)
		}
		if tool.Binary == "" && tool.Service == "" {
			t.Errorf("%s: tidak punya Binary maupun Service", tool.Name)
		}
		if tool.ServiceScope == ServiceUser && tool.Service == "" {
			t.Errorf("%s: ServiceScope user tanpa Service", tool.Name)
		}
		if tool.Linger && tool.ServiceScope != ServiceUser {
			t.Errorf("%s: Linger hanya berlaku untuk user service", tool.Name)
		}
		if tool.PidFile != "" && !strings.HasPrefix(tool.PidFile, "~/") {
			t.Errorf("%s: PidFile harus diawali ~/ (dapat %q)", tool.Name, tool.PidFile)
		}
		if tool.HealthURL != "" && !strings.HasPrefix(tool.HealthURL, "http") {
			t.Errorf("%s: HealthURL tidak valid: %q", tool.Name, tool.HealthURL)
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
