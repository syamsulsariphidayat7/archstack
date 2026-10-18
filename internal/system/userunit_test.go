package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderUnitKeyDirectives(t *testing.T) {
	content := RenderUnit(UnitOptions{
		Description: "OmniRoute AI proxy router",
		ExecStart:   "/home/user/.npm-global/bin/omniroute",
		Args:        []string{"serve", "--no-open"},
		PathEnv:     "/usr/bin:/bin",
	})

	wants := []string{
		unitMarker,
		"[Unit]",
		"Description=OmniRoute AI proxy router",
		"After=network-online.target",
		"Type=notify",
		"NotifyAccess=all",
		"WatchdogSec=180",
		"TimeoutStartSec=300",
		"ExecStart=/home/user/.npm-global/bin/omniroute serve --no-open",
		"Restart=on-failure",
		"RestartSec=5",
		`Environment="PATH=/usr/bin:/bin"`,
		"[Install]",
		"WantedBy=default.target",
	}
	for _, want := range wants {
		if !strings.Contains(content, want) {
			t.Errorf("unit tidak memuat %q:\n%s", want, content)
		}
	}
	if strings.Contains(content, "--daemon") {
		t.Error("unit tidak boleh memakai --daemon: mode itu keluar dari cgroup dan memutus notify systemd")
	}
}

func TestRenderUnitQuotesArgsWithSpaces(t *testing.T) {
	content := RenderUnit(UnitOptions{
		Description: "Test",
		ExecStart:   "/opt/my app/bin/run",
		Args:        []string{"serve", "--no-open"},
	})
	if !strings.Contains(content, `ExecStart="/opt/my app/bin/run" serve --no-open`) {
		t.Errorf("argumen dengan spasi tidak dikutip benar:\n%s", content)
	}
}

func TestSameUnitContent(t *testing.T) {
	base := RenderUnit(UnitOptions{Description: "Test", ExecStart: "/bin/true", Args: []string{"serve"}})

	// Marker + trailing whitespace + baris kosong diabaikan.
	variant := strings.Replace(base, unitMarker+"\n", "", 1)
	variant = strings.ReplaceAll(variant, "\n", "   \n")
	if !SameUnitContent(base, variant) {
		t.Error("isi unit yang hanya beda marker/whitespace harus dianggap sama")
	}

	changed := strings.Replace(base, "WatchdogSec=180", "WatchdogSec=60", 1)
	if SameUnitContent(base, changed) {
		t.Error("perubahan nilai direktif harus terdeteksi berbeda")
	}
}

func TestParseLinger(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"Linger=yes\n", true},
		{"Linger=no\n", false},
		{"", false},
		{"Linger=maybe", false},
	}
	for _, c := range cases {
		if got := parseLinger(c.in); got != c.want {
			t.Errorf("parseLinger(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestExpandHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if got, want := ExpandHome("~/.omniroute/server/.pid"), filepath.Join(home, ".omniroute/server/.pid"); got != want {
		t.Errorf("ExpandHome(~/...) = %q, want %q", got, want)
	}
	if got := ExpandHome("/abs/path"); got != "/abs/path" {
		t.Errorf("path absolut berubah: %q", got)
	}
	if got := ExpandHome("relative/path"); got != "relative/path" {
		t.Errorf("path relatif berubah: %q", got)
	}
}

func TestResolveBinaryFromPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	t.Setenv("HOME", t.TempDir())

	name := "archstack-test-binary"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveBinary(name)
	if err != nil {
		t.Fatalf("ResolveBinary(%q) error: %v", name, err)
	}
	if got != path {
		t.Errorf("ResolveBinary(%q) = %q, want %q", name, got, path)
	}
}

func TestResolveBinaryMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("HOME", t.TempDir())

	if _, err := ResolveBinary("archstack-tidak-ada-xyz"); err == nil {
		t.Error("binary yang tidak ada harus mengembalikan error")
	}
}

func TestWriteUserUnitIdempotent(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	content := RenderUnit(UnitOptions{Description: "Test", ExecStart: "/bin/true"})

	changed, backup, err := WriteUserUnit("archstack-test", content)
	if err != nil {
		t.Fatalf("WriteUserUnit pertama: %v", err)
	}
	if !changed || backup != "" {
		t.Errorf("penulisan pertama harus changed=true tanpa backup, dapat changed=%v backup=%q", changed, backup)
	}

	changed, backup, err = WriteUserUnit("archstack-test", content)
	if err != nil {
		t.Fatalf("WriteUserUnit kedua: %v", err)
	}
	if changed || backup != "" {
		t.Errorf("isi identik harus changed=false tanpa backup, dapat changed=%v backup=%q", changed, backup)
	}

	updated := strings.Replace(content, "WatchdogSec=180", "WatchdogSec=90", 1)
	changed, backup, err = WriteUserUnit("archstack-test", updated)
	if err != nil {
		t.Fatalf("WriteUserUnit update: %v", err)
	}
	if !changed || backup == "" {
		t.Errorf("isi berubah harus changed=true dengan backup, dapat changed=%v backup=%q", changed, backup)
	}
	if _, err := os.Stat(backup); err != nil {
		t.Errorf("file backup tidak ada: %v", err)
	}

	data, err := os.ReadFile(UserUnitPath("archstack-test"))
	if err != nil {
		t.Fatalf("baca unit: %v", err)
	}
	if !strings.Contains(string(data), "WatchdogSec=90") {
		t.Error("unit tidak berisi versi terbaru")
	}
}

func TestRemoveUserUnitBacksUpForeignUnit(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	path := UserUnitPath("archstack-foreign")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("[Unit]\nDescription=manual\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	backup, err := RemoveUserUnit("archstack-foreign")
	if err != nil {
		t.Fatalf("RemoveUserUnit: %v", err)
	}
	if backup == "" {
		t.Error("unit tanpa marker harus dibackup")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("unit seharusnya sudah dihapus")
	}
	if _, err := os.Stat(backup); err != nil {
		t.Errorf("backup tidak ada: %v", err)
	}
}

func TestRemoveUserUnitMissingIsNoop(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	backup, err := RemoveUserUnit("archstack-none")
	if err != nil {
		t.Errorf("RemoveUserUnit untuk unit yang tidak ada harus nil, dapat %v", err)
	}
	if backup != "" {
		t.Errorf("tidak seharusnya ada backup, dapat %q", backup)
	}
}
