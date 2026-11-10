package system

import (
	"io"
	"os/exec"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

// Varian fungsi sistem untuk pemanggil non-terminal (mis. Web UI): output
// perintah dialirkan ke io.Writer, bukan os.Stdout, dan password sudo
// diberikan lewat stdin (sudo -S). Password tidak pernah ditulis ke w.
//
// Bila pass kosong, sudo dipanggil dengan -n (non-interaktif): memakai
// cache kredensial sudo bila masih berlaku, dan langsung gagal — bukan
// menggantung menunggu input — bila password memang diperlukan.
func sudoCommand(pass, name string, args ...string) *exec.Cmd {
	return sudoCmdOpts(pass, true, "", name, args...)
}

// sudoCmdOpts menyusun perintah sudo: password via -S bila pass terisi;
// bila pass kosong dan nonInteractive, pakai -n (gagal cepat, tidak
// menggantung); bila tidak nonInteractive, sudo polos (prompt di tty).
// user mengisi sudo -u bila tidak kosong.
func sudoCmdOpts(pass string, nonInteractive bool, user, name string, args ...string) *exec.Cmd {
	var full []string
	if pass != "" {
		full = append(full, "-S", "-p", "")
	} else if nonInteractive {
		full = append(full, "-n")
	}
	if user != "" {
		full = append(full, "-u", user)
	}
	full = append(full, name)
	full = append(full, args...)
	cmd := exec.Command("sudo", full...)
	if pass != "" {
		cmd.Stdin = strings.NewReader(pass + "\n")
	}
	return cmd
}

func runTo(cmd *exec.Cmd, w io.Writer) error {
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

// InstallPacmanW memasang paket repo resmi, output ke w.
func InstallPacmanW(pkg, pass string, w io.Writer) error {
	return runTo(sudoCommand(pass, "pacman", "-S", "--needed", "--noconfirm", pkg), w)
}

// InstallYayW memasang paket AUR via yay (tanpa sudo), output ke w.
func InstallYayW(pkg string, w io.Writer) error {
	return runTo(exec.Command("yay", "-S", "--needed", "--noconfirm", pkg), w)
}

// UninstallW menghapus paket beserta dependensi yang tidak terpakai.
func UninstallW(pkg, pass string, w io.Writer) error {
	return runTo(sudoCommand(pass, "pacman", "-Rns", "--noconfirm", pkg), w)
}

// StartServiceForW menyalakan service sesuai scope tool, output ke w.
func StartServiceForW(t *registry.Tool, pass string, w io.Writer) error {
	if t.ServiceScope == registry.ServiceUser {
		return runTo(exec.Command("systemctl", "--user", "enable", "--now", t.Service), w)
	}
	return runTo(sudoCommand(pass, "systemctl", "enable", "--now", t.Service), w)
}

// StopServiceForW menghentikan service sesuai scope tool, output ke w.
func StopServiceForW(t *registry.Tool, pass string, w io.Writer) error {
	if t.ServiceScope == registry.ServiceUser {
		return runTo(exec.Command("systemctl", "--user", "stop", t.Service), w)
	}
	return runTo(sudoCommand(pass, "systemctl", "stop", t.Service), w)
}
