package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

// installServiceOnly memasang service untuk binary yang sudah ada. archstack
// tidak memasang/menghapus binary pada mode ini.
func installServiceOnly(tool *registry.Tool) error {
	if tool.Service == "" {
		return fmt.Errorf("[error] %s tidak punya nama service", tool.Name)
	}
	if tool.ServiceScope != registry.ServiceUser {
		return fmt.Errorf("[error] %s: service-only saat ini hanya mendukung user service", tool.Name)
	}

	fmt.Printf("[install] Menyiapkan service %s (binary tidak dipasang archstack)\n", tool.Name)

	binPath, err := system.ResolveBinary(tool.Binary)
	if err != nil {
		return fmt.Errorf("[error] %v\n        Pasang %s dulu secara manual (mis. npm i -g %s), lalu ulangi",
			err, tool.Binary, tool.Name)
	}
	fmt.Printf("[info] Binary    : %s\n", binPath)

	if err := ensureNoManualInstance(tool); err != nil {
		return err
	}

	unitDesc := tool.UnitDescription
	if unitDesc == "" {
		unitDesc = tool.Desc
	}
	content := system.RenderUnit(system.UnitOptions{
		Description: unitDesc,
		ExecStart:   binPath,
		Args:        tool.UnitArgs,
		PathEnv:     servicePathEnv(binPath),
	})
	unitPath := system.UserUnitPath(tool.Service)
	changed, backup, err := system.WriteUserUnit(tool.Service, content)
	if err != nil {
		return fmt.Errorf("[error] Gagal menulis unit %s: %w", unitPath, err)
	}
	switch {
	case changed && backup != "":
		fmt.Printf("[info] Unit lama dibackup: %s\n", backup)
		fmt.Printf("[done] Unit ditulis  : %s\n", unitPath)
	case changed:
		fmt.Printf("[done] Unit ditulis  : %s\n", unitPath)
	default:
		fmt.Printf("[ok] Unit sudah sesuai: %s\n", unitPath)
	}

	if err := system.UserDaemonReload(); err != nil {
		return fmt.Errorf("[error] systemctl --user daemon-reload gagal: %w", err)
	}

	// Unit yang berubah tidak otomatis dipakai proses yang sedang jalan.
	if state, _ := system.UserServiceStatus(tool.Service); changed && state == system.ServiceRunning {
		fmt.Printf("[run] Unit berubah — restart %s...\n", tool.Service)
		if err := system.RestartUserService(tool.Service); err != nil {
			return fmt.Errorf("[error] Gagal restart %s: %w\n        Cek: journalctl --user -u %s -n 50",
				tool.Service, err, tool.Service)
		}
	}

	fmt.Printf("[run] Mengaktifkan service %s...\n", tool.Service)
	if err := system.StartUserService(tool.Service); err != nil {
		return fmt.Errorf("[error] Gagal menyalakan %s: %w\n        Cek: journalctl --user -u %s -n 50",
			tool.Service, err, tool.Service)
	}

	state, _ := system.UserServiceStatus(tool.Service)
	if state != system.ServiceRunning {
		return fmt.Errorf("[error] Service %s tidak aktif setelah start\n        Cek: journalctl --user -u %s -n 50",
			tool.Service, tool.Service)
	}
	fmt.Printf("[ok] Service %s aktif\n", tool.Service)

	if tool.Linger {
		if enabled, err := system.LingerEnabled(); err == nil && enabled {
			fmt.Printf("[ok] Linger aktif — service tetap jalan setelah boot tanpa login\n")
		} else if err := system.EnableLinger(); err != nil {
			fmt.Printf("[info] Linger belum aktif; jalankan manual: loginctl enable-linger %s\n",
				system.CurrentUsername())
		} else {
			fmt.Printf("[done] Linger diaktifkan (autostart saat boot)\n")
		}
	}

	if tool.HealthURL != "" {
		fmt.Printf("[info] Menunggu health check %s...\n", tool.HealthURL)
		if err := system.WaitHealthy(tool.HealthURL, 30*time.Second); err != nil {
			fmt.Printf("[info] Health check belum OK (%v) — service tetap terpasang\n", err)
		} else {
			fmt.Printf("[ok] Health check OK\n")
		}
	}

	return nil
}

// uninstallServiceOnly menghapus unit + mematikan service. Binary dan data
// milik tool tidak disentuh.
func uninstallServiceOnly(tool *registry.Tool) error {
	if tool.Service == "" {
		return fmt.Errorf("[error] %s tidak punya nama service", tool.Name)
	}
	if tool.ServiceScope != registry.ServiceUser {
		return fmt.Errorf("[error] %s: service-only hanya mendukung user service", tool.Name)
	}

	unitExists := system.UserUnitExists(tool.Service)
	state, _ := system.UserServiceStatus(tool.Service)
	if !unitExists && state != system.ServiceRunning {
		fmt.Printf("[ok] Service %s belum terpasang\n", tool.Service)
		fmt.Printf("[info] Binary dan data %s dibiarkan\n", tool.Name)
		return nil
	}

	fmt.Printf("[run] Menonaktifkan service %s...\n", tool.Service)
	if err := system.DisableUserService(tool.Service); err != nil {
		fmt.Printf("[info] disable --now gagal (lanjut hapus unit): %v\n", err)
	}

	if backup, err := system.RemoveUserUnit(tool.Service); err != nil {
		return fmt.Errorf("[error] Gagal menghapus unit: %w", err)
	} else if backup != "" {
		fmt.Printf("[info] Unit tidak bertanda archstack, dibackup: %s\n", backup)
	}

	if err := system.UserDaemonReload(); err != nil {
		fmt.Printf("[info] daemon-reload gagal: %v\n", err)
	}

	fmt.Printf("[done] Service %s dihapus\n", tool.Service)
	fmt.Printf("[info] Binary dan data %s tidak disentuh\n", tool.Name)
	return nil
}

// ensureNoManualInstance menolak install bila tool sedang jalan di luar systemd
// (mis. daemon yang dijalankan manual), karena port-nya akan bentrok.
func ensureNoManualInstance(tool *registry.Tool) error {
	if tool.PidFile == "" {
		return nil
	}
	// Service yang sedang aktif sudah dikelola systemd (PID di file bisa saja
	// child process, bukan MainPID), jadi bukan konflik.
	if state, _ := system.UserServiceStatus(tool.Service); state == system.ServiceRunning {
		return nil
	}
	pid, ok := system.ReadPidFile(tool.PidFile)
	if !ok || !system.ProcessAlive(pid) {
		return nil
	}
	if system.UserServiceMainPID(tool.Service) == pid {
		return nil
	}
	return fmt.Errorf("[error] %s sedang berjalan manual (PID %d) di luar systemd\n"+
		"        Hentikan dulu: %s stop — lalu ulangi perintah ini", tool.Name, pid, tool.Binary)
}

// servicePathEnv menyusun PATH untuk unit agar shebang `#!/usr/bin/env node`
// dan helper CLI tetap ketemu saat dijalankan systemd.
func servicePathEnv(binPath string) string {
	dirs := []string{filepath.Dir(binPath)}
	if nodePath, err := exec.LookPath("node"); err == nil {
		dirs = append(dirs, filepath.Dir(nodePath))
	}
	dirs = append(dirs, "/usr/local/sbin", "/usr/local/bin", "/usr/bin", "/bin")

	seen := make(map[string]bool, len(dirs))
	var out []string
	for _, dir := range dirs {
		if dir == "" || seen[dir] {
			continue
		}
		seen[dir] = true
		out = append(out, dir)
	}
	return strings.Join(out, ":")
}

// printServiceDetail menampilkan detail tambahan untuk tool service-only.
func printServiceDetail(tool *registry.Tool) {
	fmt.Printf("\nDetail %s:\n", tool.Name)
	if path, err := system.ResolveBinary(tool.Binary); err == nil {
		fmt.Printf("  binary : %s\n", path)
	} else {
		fmt.Printf("  binary : tidak ditemukan\n")
	}
	unitState := "tidak ada"
	if system.UserUnitExists(tool.Service) {
		unitState = "ada"
	}
	fmt.Printf("  unit   : %s (%s)\n", system.UserUnitPath(tool.Service), unitState)

	state, _ := system.UserServiceStatus(tool.Service)
	if state == system.ServiceRunning {
		if pid := system.UserServiceMainPID(tool.Service); pid > 0 {
			fmt.Printf("  pid    : %d\n", pid)
		}
	}
	if tool.Linger {
		linger := "tidak aktif"
		if enabled, err := system.LingerEnabled(); err == nil && enabled {
			linger = "aktif"
		}
		fmt.Printf("  linger : %s\n", linger)
	}
	if tool.HealthURL != "" && state == system.ServiceRunning {
		if err := system.WaitHealthy(tool.HealthURL, 3*time.Second); err != nil {
			fmt.Printf("  health : gagal (%v)\n", err)
		} else {
			fmt.Printf("  health : OK\n")
		}
	}
}
