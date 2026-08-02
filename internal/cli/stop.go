package cli

import (
	"fmt"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func stopCmd(name string) error {
	tool, ok := registry.GetTool(name)
	if !ok {
		return fmt.Errorf("[error] Tool tidak dikenal: %s", name)
	}

	if tool.Service == "" {
		return fmt.Errorf("[info] %s tidak punya service untuk distop\n", name)
	}

	state, _ := system.ServiceStatus(tool.Service)
	if state != system.ServiceRunning {
		fmt.Printf("[ok] Service %s tidak sedang berjalan\n", tool.Service)
		return nil
	}

	fmt.Printf("[run] Menghentikan service %s...\n", tool.Service)
	if err := system.StopService(tool.Service); err != nil {
		return fmt.Errorf("[error] Gagal stop service %s: %w", tool.Service, err)
	}
	fmt.Printf("[done] Service %s berhenti\n", tool.Service)
	return nil
}
