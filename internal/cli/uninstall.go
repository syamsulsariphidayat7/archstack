package cli

import (
	"fmt"
	"os"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func uninstallCmd(names []string) error {
	hasError := false
	for _, name := range names {
		tool, ok := registry.GetTool(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "[error] Tool tidak dikenal: %s\n", name)
			hasError = true
			continue
		}
		if !system.IsInstalled(tool.Binary, tool.Pkg) {
			fmt.Printf("[ok] %s sudah tidak terinstal\n", name)
			continue
		}
		if tool.Service != "" {
			state, _ := system.ServiceStatus(tool.Service)
			if state == system.ServiceRunning {
				fmt.Printf("[info] Menghentikan service %s...\n", tool.Service)
				if err := system.StopService(tool.Service); err != nil {
					fmt.Fprintf(os.Stderr, "[error] Gagal stop service %s: %v\n", tool.Service, err)
					hasError = true
					continue
				}
			}
		}
		fmt.Printf("[install] Menghapus %s...\n", name)
		if err := system.Uninstall(tool.Pkg); err != nil {
			fmt.Fprintf(os.Stderr, "[error] Gagal hapus %s: %v\n", name, err)
			hasError = true
			continue
		}
		fmt.Printf("[done] %s berhasil dihapus\n", name)
	}
	if hasError {
		return fmt.Errorf("beberapa tool gagal dihapus")
	}
	return nil
}
