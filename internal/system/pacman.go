package system

import (
	"os"
	"os/exec"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

// AvailableVersion mengembalikan versi terbaru yang tersedia untuk pkg
// dari repo resmi (pacman) atau AUR (yay). Kosong jika query gagal.
func AvailableVersion(pkg string, aur bool) string {
	if pkg == "" {
		return ""
	}
	cmdName := "pacman"
	if aur {
		cmdName = "yay"
	}
	out, err := exec.Command(cmdName, "-Si", pkg).Output()
	if err != nil {
		return ""
	}
	return parseVersion(string(out))
}

// parseVersion mengekstrak nilai baris "Version : x.y.z" dari output
// pacman -Si / yay -Si. Mengembalikan "" jika tidak ketemu.
func parseVersion(out string) string {
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Version" {
			continue
		}
		v := strings.TrimSpace(parts[1])
		if v != "" {
			return v
		}
	}
	return ""
}

// InstalledVersions mengembalikan versi terinstall per package dari `pacman -Q`.
// Kosong/nil jika query gagal.
func InstalledVersions() map[string]string {
	out, err := exec.Command("pacman", "-Q").Output()
	if err != nil {
		return nil
	}
	return parsePacmanQ(string(out))
}

// AvailableVersions mengembalikan versi yang tersedia di repo resmi dari `pacman -Sl`.
// Kosong/nil jika query gagal.
func AvailableVersions() map[string]string {
	out, err := exec.Command("pacman", "-Sl").Output()
	if err != nil {
		return nil
	}
	return parsePacmanSl(string(out))
}

// ToolVersionStrings mengembalikan versi tampilan per tool (key: tool.Name):
// versi terinstall jika ada, selain itu versi yang tersedia (repo resmi / AUR).
func ToolVersionStrings(tools []registry.Tool) map[string]string {
	installed := InstalledVersions()
	available := AvailableVersions()
	out := make(map[string]string, len(tools))
	for _, t := range tools {
		if t.Pkg == "" {
			continue
		}
		if v := installed[t.Pkg]; v != "" {
			out[t.Name] = v
			continue
		}
		if v := available[t.Pkg]; v != "" {
			out[t.Name] = v
			continue
		}
		if t.From == registry.SourceYay {
			out[t.Name] = AvailableVersion(t.Pkg, true)
		}
	}
	return out
}

// parsePacmanQ mengekstrak "nama versi" per baris dari output `pacman -Q`.
func parsePacmanQ(out string) map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) >= 2 {
			m[f[0]] = f[1]
		}
	}
	return m
}

// parsePacmanSl mengekstrak "repo nama versi [marker]" per baris dari output `pacman -Sl`.
func parsePacmanSl(out string) map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) >= 3 {
			m[f[1]] = f[2]
		}
	}
	return m
}

func IsInstalled(binary, pkg string) bool {
	if binary != "" {
		if _, err := exec.LookPath(binary); err == nil {
			return true
		}
	}
	if pkg != "" {
		out, _ := exec.Command("pacman", "-Qi", pkg).Output()
		return strings.Contains(string(out), "Name")
	}
	return false
}

func InstallPacman(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-S", "--needed", "--noconfirm", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InstallYay(pkg string) error {
	cmd := exec.Command("yay", "-S", "--needed", "--noconfirm", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Uninstall(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
