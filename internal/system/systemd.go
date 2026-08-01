package system

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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
