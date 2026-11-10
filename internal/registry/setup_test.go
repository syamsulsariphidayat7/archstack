package registry

import "testing"

// Hanya postgres dan mariadb yang butuh setup satu kali (inisialisasi
// direktori data sebelum service pertama kali bisa start).
func TestSetupFields(t *testing.T) {
	pg, ok := GetTool("postgres")
	if !ok || !pg.NeedsSetup() {
		t.Fatal("postgres harus NeedsSetup")
	}
	if pg.SetupDir != "/var/lib/postgres/data" || pg.SetupMarker != "PG_VERSION" || pg.SetupUser != "postgres" {
		t.Errorf("field setup postgres salah: %+v", pg)
	}
	if len(pg.SetupCmd) == 0 || pg.SetupCmd[0] != "initdb" {
		t.Errorf("SetupCmd postgres harus diawali initdb: %v", pg.SetupCmd)
	}

	my, ok := GetTool("mariadb")
	if !ok || !my.NeedsSetup() {
		t.Fatal("mariadb harus NeedsSetup")
	}
	if len(my.SetupCmd) == 0 || my.SetupCmd[0] != "mariadb-install-db" {
		t.Errorf("SetupCmd mariadb harus diawali mariadb-install-db: %v", my.SetupCmd)
	}

	for _, name := range []string{"redis", "nginx", "docker", "git", "omniroute", "archstack-web"} {
		tool, ok := GetTool(name)
		if !ok {
			t.Fatalf("tool %q tidak ada", name)
		}
		if tool.NeedsSetup() {
			t.Errorf("%s tidak boleh NeedsSetup", name)
		}
	}
}
