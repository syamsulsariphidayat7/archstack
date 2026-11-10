package web

import (
	"net"
	"os"
	"strconv"
	"time"
)

// Dukungan sd_notify: bila dijalankan sebagai service systemd bertipe
// notify (unit hasil RenderUnit di internal/system selalu Type=notify +
// WatchdogSec), server wajib melapor READY setelah listen dan mengirim
// ping watchdog berkala — bila tidak, systemd menganggap start gagal dan
// watchdog membunuh prosesnya. Bila NOTIFY_SOCKET tidak ada (server
// dijalankan manual dari terminal), semua fungsi di sini no-op.
//
// NOTIFY_SOCKET ditangkap sekali saat start lalu DIHAPUS dari environment
// proses: child process (pacman/systemctl yang dipanggil handler) tidak
// boleh mewarisinya — systemctl mengirim datagram EXIT_STATUS ke socket
// warisan itu, yang mengotori kanal notifikasi service ini.

// notifyTarget adalah alamat socket hasil tangkapan initNotify.
var notifyTarget string

// notifyAddr menerjemahkan NOTIFY_SOCKET menjadi alamat dial: socket
// abstrak diawali "@" yang diganti byte NUL.
func notifyAddr(raw string) string {
	if raw == "" {
		return ""
	}
	if raw[0] == '@' {
		return "\x00" + raw[1:]
	}
	return raw
}

// initNotify menangkap konfigurasi notifikasi dari environment, lalu
// menghapus variabelnya dari environment proses (lihat catatan di atas).
// Mengembalikan interval ping watchdog dan apakah watchdog aktif.
func initNotify() (time.Duration, bool) {
	notifyTarget = notifyAddr(os.Getenv("NOTIFY_SOCKET"))
	interval, ok := watchdogInterval()
	os.Unsetenv("NOTIFY_SOCKET")
	os.Unsetenv("WATCHDOG_USEC")
	os.Unsetenv("WATCHDOG_PID")
	return interval, ok
}

func sdNotify(state string) {
	if notifyTarget == "" {
		return
	}
	conn, err := net.Dial("unixgram", notifyTarget)
	if err != nil {
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte(state))
}

// parseWatchdogInterval mengubah WATCHDOG_USEC (mikrodetik) menjadi
// interval ping setengahnya, sesuai anjuran systemd.
func parseWatchdogInterval(usecStr string) (time.Duration, bool) {
	usec, err := strconv.ParseInt(usecStr, 10, 64)
	if err != nil || usec <= 0 {
		return 0, false
	}
	return time.Duration(usec) * time.Microsecond / 2, true
}

// watchdogInterval membaca konfigurasi watchdog dari environment, dan
// menghormati WATCHDOG_PID (ping hanya bila watchdog memang untuk proses ini).
func watchdogInterval() (time.Duration, bool) {
	if pidStr := os.Getenv("WATCHDOG_PID"); pidStr != "" {
		pid, err := strconv.Atoi(pidStr)
		if err != nil || pid != os.Getpid() {
			return 0, false
		}
	}
	return parseWatchdogInterval(os.Getenv("WATCHDOG_USEC"))
}

func runWatchdog(interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sdNotify("WATCHDOG=1")
		case <-stop:
			return
		}
	}
}
