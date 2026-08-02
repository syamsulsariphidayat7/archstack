package tui

import (
	"os/exec"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

type toolState struct {
	installed bool
	svc       system.ServiceState
}

func snapshotStates(tools []registry.Tool) map[string]toolState {
	installed := installedPkgs()
	states := make(map[string]toolState, len(tools))
	for _, t := range tools {
		st := toolState{installed: installed[t.Pkg]}
		if !st.installed && t.Binary != "" {
			if _, err := exec.LookPath(t.Binary); err == nil {
				st.installed = true
			}
		}
		if t.Service != "" {
			st.svc, _ = system.ServiceStatus(t.Service)
		}
		states[t.Name] = st
	}
	return states
}

func installedPkgs() map[string]bool {
	set := make(map[string]bool)
	out, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return set
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			set[line] = true
		}
	}
	return set
}
