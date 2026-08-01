package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/syamsulsariphidayat7/archstack/internal/tui"
)

func Root(args []string) error {
	if len(args) == 0 {
		return tui.Run()
	}

	cmd := args[0]

	switch cmd {
	case "list":
		return listCmd()
	case "install":
		if len(args) < 2 {
			return errors.New("[error] Gunakan: archstack install <tool> [tool2 ...]")
		}
		return installCmd(args[1:])
	case "uninstall":
		if len(args) < 2 {
			return errors.New("[error] Gunakan: archstack uninstall <tool> [tool2 ...]")
		}
		return uninstallCmd(args[1:])
	case "run":
		if len(args) < 2 {
			return errors.New("[error] Gunakan: archstack run <tool>")
		}
		return runCmd(args[1])
	case "--help", "-h":
		printHelp()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "[error] Perintah tidak dikenal: %s\n", cmd)
		fmt.Fprintf(os.Stderr, "Gunakan 'archstack --help' untuk bantuan\n")
		return errors.New("unknown command")
	}
}

func printHelp() {
	fmt.Println(`archstack — Development Tools Manager

Penggunaan:
  archstack              Mode interaktif (TUI)
  archstack list         Tampilkan daftar tool
  archstack install <tool> [tool2 ...]   Instal tool
  archstack uninstall <tool> [tool2 ...]  Hapus tool
  archstack run <tool>   Jalankan service tool

Contoh:
  archstack install nginx
  archstack install php node postgres
  archstack uninstall nginx
  archstack run postgres`)
}
