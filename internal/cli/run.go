package cli

import (
	"fmt"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func runCmd(name string) error {
	tool, ok := registry.GetTool(name)
	if !ok {
		return fmt.Errorf("[error] Tool tidak dikenal: %s", name)
	}

	if tool.Service != "" {
		state, _ := system.ServiceStatus(tool.Service)
		if state == system.ServiceRunning {
			fmt.Printf("[ok] Service %s sudah berjalan\n", tool.Service)
			return nil
		}
		fmt.Printf("[run] Menjalankan service %s...\n", tool.Service)
		if err := system.StartService(tool.Service); err != nil {
			return fmt.Errorf("[error] Gagal start service %s: %w", tool.Service, err)
		}
		fmt.Printf("[done] Service %s berjalan\n", tool.Service)
		return nil
	}

	if tool.Binary != "" {
		fmt.Printf("[info] %s bukan service. Jalankan manual: %s\n", tool.Name, tool.Binary)
	} else {
		fmt.Printf("[info] %s tidak punya binary atau service untuk dijalankan\n", tool.Name)
	}
	return nil
}
