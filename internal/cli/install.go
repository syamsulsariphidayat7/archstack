package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func listCmd() error {
	return printStatusTable(registry.AllTools())
}

func installCmd(names []string) error {
	hasError := false
	for _, name := range names {
		tool, ok := registry.GetTool(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "[error] Tool tidak dikenal: %s\n", name)
			hasError = true
			continue
		}
		if !tool.IsServiceOnly() && system.IsInstalled(tool.Binary, tool.Pkg) {
			fmt.Printf("[ok] %s sudah terinstal\n", name)
			continue
		}
		if err := installOne(tool); err != nil {
			printInstallError(name, err)
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("beberapa tool gagal diinstal")
	}
	return nil
}

// printInstallError menghindari prefix ganda untuk error yang sudah memakai
// format "[error] ..." (mis. petunjuk multi-baris dari modul service).
func printInstallError(name string, err error) {
	msg := err.Error()
	if strings.HasPrefix(msg, "[error]") {
		fmt.Fprintln(os.Stderr, msg)
		return
	}
	fmt.Fprintf(os.Stderr, "[error] %s: %v\n", name, err)
}

func installOne(tool *registry.Tool) error {
	if tool.IsServiceOnly() {
		return installServiceOnly(tool)
	}
	ver := system.AvailableVersion(tool.Pkg, tool.From == registry.SourceYay)
	if ver != "" {
		fmt.Printf("[install] Menginstal %s — versi rekomendasi: %s...\n", tool.Name, ver)
	} else {
		fmt.Printf("[install] Menginstal %s...\n", tool.Name)
	}
	switch tool.From {
	case registry.SourcePacman:
		return system.InstallPacman(tool.Pkg)
	case registry.SourceYay:
		return system.InstallYay(tool.Pkg)
	}
	return nil
}
