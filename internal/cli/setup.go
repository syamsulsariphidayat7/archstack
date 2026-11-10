package cli

import (
	"fmt"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

// setupCmd menjalankan inisialisasi satu kali untuk tool yang butuh setup
// sebelum service-nya bisa start (postgres: initdb, mariadb: install-db).
// Idempotent: bila penanda setup sudah ada, tidak melakukan apa-apa.
func setupCmd(name string) error {
	tool, ok := registry.GetTool(name)
	if !ok {
		return fmt.Errorf("[error] Tool tidak dikenal: %s", name)
	}
	if !tool.NeedsSetup() {
		return fmt.Errorf("[error] %s tidak butuh setup — langsung install/run saja", name)
	}
	if !system.IsInstalled(tool.Binary, tool.Pkg) {
		return fmt.Errorf("[error] %s belum terinstal — jalankan dulu: archstack install %s", name, name)
	}
	if system.SetupDone(tool) {
		fmt.Printf("[ok] %s sudah di-setup (%s/%s ada) — tidak ada yang dilakukan\n", name, tool.SetupDir, tool.SetupMarker)
		return nil
	}
	fmt.Printf("[setup] Inisialisasi %s: %v\n", name, tool.SetupCmd)
	if err := system.RunSetup(tool); err != nil {
		return fmt.Errorf("[error] Setup %s gagal: %w", name, err)
	}
	if !system.SetupDone(tool) {
		return fmt.Errorf("[error] Setup %s selesai tapi penanda %s/%s tidak ditemukan — cek output di atas", name, tool.SetupDir, tool.SetupMarker)
	}
	fmt.Printf("[done] %s berhasil di-setup — sekarang bisa: archstack run %s\n", name, name)
	if tool.SetupNote != "" {
		fmt.Printf("[info] Berikutnya: %s\n", tool.SetupNote)
	}
	return nil
}
