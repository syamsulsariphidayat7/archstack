package web

import (
	"errors"
	"testing"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

func mustTool(t *testing.T, name string) *registry.Tool {
	t.Helper()
	tool, ok := registry.GetTool(name)
	if !ok {
		t.Fatalf("tool %q tidak ada di registry", name)
	}
	return tool
}

func TestValidateAction(t *testing.T) {
	git := mustTool(t, "git")     // tanpa service
	nginx := mustTool(t, "nginx") // dengan service

	cases := []struct {
		tool    *registry.Tool
		action  string
		wantErr bool
	}{
		{git, "install", false},
		{git, "uninstall", false},
		{git, "run", true},
		{git, "stop", true},
		{nginx, "run", false},
		{nginx, "stop", false},
		{nginx, "restart", true},
		{nginx, "", true},
	}
	for _, c := range cases {
		if err := validateAction(c.tool, c.action); (err != nil) != c.wantErr {
			t.Errorf("validateAction(%s, %q) err=%v, wantErr=%v", c.tool.Name, c.action, err, c.wantErr)
		}
	}
}

func TestNeedsSudo(t *testing.T) {
	cases := []struct {
		tool   string
		action string
		want   bool
	}{
		{"php", "install", true},        // pacman → sudo
		{"pnpm", "install", false},      // yay → user biasa
		{"mongodb", "install", false},   // yay (AUR)
		{"php", "uninstall", true},      // pacman -Rns → sudo
		{"postgres", "run", true},       // service system → sudo
		{"postgres", "stop", true},      // service system → sudo
		{"omniroute", "install", false}, // service-only user service
		{"omniroute", "run", false},     // systemctl --user
		{"omniroute", "stop", false},
		{"omniroute", "uninstall", false},
	}
	for _, c := range cases {
		tool := mustTool(t, c.tool)
		if got := needsSudo(tool, c.action); got != c.want {
			t.Errorf("needsSudo(%s, %q) = %v, want %v", c.tool, c.action, got, c.want)
		}
	}
}

func TestMakeView(t *testing.T) {
	omni := mustTool(t, "omniroute")
	v := makeView(*omni, true, "running", "", true)
	if v.Scope != "user" || !v.ServiceOnly || v.Pkg != "— service only —" {
		t.Errorf("view omniroute salah: %+v", v)
	}
	nginx := mustTool(t, "nginx")
	v = makeView(*nginx, false, "stopped", "1.28.1", false)
	if v.Scope != "system" || v.ServiceOnly || v.Pkg != "nginx" || v.Svc != "stopped" {
		t.Errorf("view nginx salah: %+v", v)
	}
}

func TestJobManagerSatuProsesDanLogBaris(t *testing.T) {
	m := &jobManager{}
	if !m.tryStart("nginx", "install") {
		t.Fatal("tryStart pertama harus berhasil")
	}
	if m.tryStart("redis", "install") {
		t.Fatal("tryStart kedua saat proses jalan harus ditolak")
	}
	w := &lineWriter{m: m}
	if _, err := w.Write([]byte("[install] satu\r\ndua")); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("r\ntiga")); err != nil {
		t.Fatal(err)
	}
	m.finish(errors.New("gagal contoh"))

	snap := m.snapshot()
	if snap.Running {
		t.Error("job harus selesai setelah finish")
	}
	want := []string{"[install] satu", "duar", "tiga"}
	if len(snap.Lines) != len(want) {
		t.Fatalf("lines = %v, want %v", snap.Lines, want)
	}
	for i := range want {
		if snap.Lines[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, snap.Lines[i], want[i])
		}
	}
	if snap.Error != "gagal contoh" {
		t.Errorf("error = %q, want %q", snap.Error, "gagal contoh")
	}
	if !m.tryStart("redis", "run") {
		t.Error("setelah finish, job baru harus bisa mulai")
	}
	m.finish(nil)
}
