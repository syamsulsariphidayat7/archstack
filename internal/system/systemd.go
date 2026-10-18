package system

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

type ServiceState int

const (
	ServiceNone ServiceState = iota
	ServiceRunning
	ServiceStopped
)

func ServiceStatus(svc string) (ServiceState, error) {
	if svc == "" {
		return ServiceNone, nil
	}
	out, err := exec.Command("systemctl", "is-active", svc).Output()
	if err != nil {
		return ServiceStopped, nil
	}
	if strings.TrimSpace(string(out)) == "active" {
		return ServiceRunning, nil
	}
	return ServiceStopped, nil
}

func StartService(svc string) error {
	cmd := exec.Command("sudo", "systemctl", "enable", "--now", svc)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StopService(svc string) error {
	cmd := exec.Command("sudo", "systemctl", "stop", svc)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func HasServiceUnit(svc string) bool {
	if svc == "" {
		return false
	}
	out, _ := exec.Command("systemctl", "list-unit-files", svc+".service", "--no-pager").Output()
	return bytes.Contains(out, []byte(svc+".service"))
}

// ServiceStatusFor membaca status service sesuai scope tool (system/user).
func ServiceStatusFor(t *registry.Tool) (ServiceState, error) {
	if t.Service == "" {
		return ServiceNone, nil
	}
	if t.ServiceScope == registry.ServiceUser {
		return UserServiceStatus(t.Service)
	}
	return ServiceStatus(t.Service)
}

// StartServiceFor menyalakan service sesuai scope tool (system/user).
func StartServiceFor(t *registry.Tool) error {
	if t.ServiceScope == registry.ServiceUser {
		return StartUserService(t.Service)
	}
	return StartService(t.Service)
}

// StopServiceFor menghentikan service sesuai scope tool (system/user).
func StopServiceFor(t *registry.Tool) error {
	if t.ServiceScope == registry.ServiceUser {
		return StopUserService(t.Service)
	}
	return StopService(t.Service)
}

// StartServiceCommand mengembalikan command start untuk dipakai TUI
// (tea.ExecProcess): `systemctl --user enable --now` atau `sudo systemctl ...`.
func StartServiceCommand(t *registry.Tool) *exec.Cmd {
	if t.ServiceScope == registry.ServiceUser {
		return exec.Command("systemctl", "--user", "enable", "--now", t.Service)
	}
	return exec.Command("sudo", "systemctl", "enable", "--now", t.Service)
}

// StopServiceCommand mengembalikan command stop untuk dipakai TUI.
func StopServiceCommand(t *registry.Tool) *exec.Cmd {
	if t.ServiceScope == registry.ServiceUser {
		return exec.Command("systemctl", "--user", "stop", t.Service)
	}
	return exec.Command("sudo", "systemctl", "stop", t.Service)
}
