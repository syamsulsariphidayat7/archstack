package system

import (
	"os"
	"os/exec"
	"strings"
)

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
