package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

func TestSetupDonePenanda(t *testing.T) {
	dir := t.TempDir()
	tool := &registry.Tool{
		Name: "contoh", SetupDir: dir, SetupMarker: "PG_VERSION",
		SetupCmd: []string{"initdb", "-D", dir},
	}
	if SetupDone(tool) {
		t.Error("penanda belum ada — SetupDone harus false")
	}
	if err := os.WriteFile(filepath.Join(dir, "PG_VERSION"), []byte("18\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !SetupDone(tool) {
		t.Error("penanda sudah ada — SetupDone harus true")
	}

	tanpaSetup := &registry.Tool{Name: "redis"}
	if SetupDone(tanpaSetup) {
		t.Error("tool tanpa SetupCmd harus false")
	}
}

func TestSudoCmdOpts(t *testing.T) {
	// Web dengan password: -S, dan -u bila ada user.
	cmd := sudoCmdOpts("rahasia", true, "postgres", "initdb", "-D", "/data")
	got := strings.Join(cmd.Args, " ")
	want := "sudo -S -p  -u postgres initdb -D /data"
	if got != want {
		t.Errorf("argv = %q, want %q", got, want)
	}
	// Web tanpa password: -n (gagal cepat).
	cmd = sudoCmdOpts("", true, "", "pacman", "-S", "nginx")
	if got := strings.Join(cmd.Args, " "); got != "sudo -n pacman -S nginx" {
		t.Errorf("argv = %q", got)
	}
	// CLI interaktif: sudo polos tanpa -S/-n (prompt tty).
	cmd = sudoCmdOpts("", false, "postgres", "initdb")
	if got := strings.Join(cmd.Args, " "); got != "sudo -u postgres initdb" {
		t.Errorf("argv = %q", got)
	}
}
