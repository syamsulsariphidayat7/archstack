package registry

import "testing"

// archstack-web adalah tool service-only kedua (setelah omniroute): Web UI
// archstack sendiri, dijalankan sebagai user service agar bisa autostart
// dan dikelola dari dashboard/TUI/CLI archstack.
func TestArchstackWebEntry(t *testing.T) {
	tool, ok := GetTool("archstack-web")
	if !ok {
		t.Fatal("archstack-web harus ada di registry")
	}
	if !tool.IsServiceOnly() {
		t.Error("archstack-web harus service-only (binary tidak dipasang/dihapus archstack)")
	}
	if tool.Binary != "archstack" || tool.Service != "archstack-web" {
		t.Errorf("binary/service salah: %+v", tool)
	}
	if tool.ServiceScope != ServiceUser {
		t.Error("archstack-web harus user service (tanpa sudo)")
	}
	if !tool.Linger {
		t.Error("archstack-web harus linger agar autostart setelah boot")
	}
	if len(tool.UnitArgs) != 1 || tool.UnitArgs[0] != "web" {
		t.Errorf("UnitArgs harus [web], dapat %v", tool.UnitArgs)
	}
	if tool.HealthURL != "http://127.0.0.1:8090/api/tools" {
		t.Errorf("HealthURL salah: %q", tool.HealthURL)
	}
}
