package system

import (
	"io"
	"os"
	"path/filepath"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

// SetupDone melaporkan apakah setup satu kali tool sudah pernah dilakukan,
// dilihat dari berkas penanda di direktori datanya (mis. PG_VERSION untuk
// postgres). Tool tanpa SetupCmd selalu false.
func SetupDone(t *registry.Tool) bool {
	if !t.NeedsSetup() || t.SetupDir == "" || t.SetupMarker == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(t.SetupDir, t.SetupMarker))
	return err == nil
}

// RunSetup menjalankan perintah inisialisasi tool di terminal (untuk CLI):
// sudo polos, password ditanya lewat tty.
func RunSetup(t *registry.Tool) error {
	if !t.NeedsSetup() {
		return nil
	}
	cmd := sudoCmdOpts("", false, t.SetupUser, t.SetupCmd[0], t.SetupCmd[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunSetupW menjalankan perintah inisialisasi tool untuk Web UI: output ke
// w, password sudo via sudo -S (atau -n bila kosong).
func RunSetupW(t *registry.Tool, pass string, w io.Writer) error {
	if !t.NeedsSetup() {
		return nil
	}
	return runTo(sudoCmdOpts(pass, true, t.SetupUser, t.SetupCmd[0], t.SetupCmd[1:]...), w)
}
