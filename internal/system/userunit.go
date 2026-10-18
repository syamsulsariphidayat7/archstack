package system

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// unitMarker menandai unit yang ditulis archstack. Dipakai untuk membedakan
// unit yang boleh ditimpa/dihapus dari unit yang dibuat manual oleh user.
const unitMarker = "# Managed by archstack"

// UnitOptions adalah bahan untuk merender unit systemd user.
type UnitOptions struct {
	Description     string
	ExecStart       string
	Args            []string
	PathEnv         string
	WatchdogSec     int
	StartTimeoutSec int
	RestartSec      int
}

// RenderUnit menghasilkan isi file unit. Service sengaja dijalankan di
// foreground (Type=notify): mode daemon internal akan keluar dari cgroup dan
// memutus handshake READY/WATCHDOG dengan systemd.
func RenderUnit(opts UnitOptions) string {
	if opts.WatchdogSec <= 0 {
		opts.WatchdogSec = 180
	}
	if opts.StartTimeoutSec <= 0 {
		opts.StartTimeoutSec = 300
	}
	if opts.RestartSec <= 0 {
		opts.RestartSec = 5
	}

	execParts := []string{quoteUnitArg(opts.ExecStart)}
	for _, arg := range opts.Args {
		execParts = append(execParts, quoteUnitArg(arg))
	}

	var b strings.Builder
	b.WriteString(unitMarker + "\n")
	b.WriteString("[Unit]\n")
	b.WriteString("Description=" + opts.Description + "\n")
	b.WriteString("After=network-online.target\n")
	b.WriteString("Wants=network-online.target\n")
	b.WriteString("\n[Service]\n")
	b.WriteString("Type=notify\n")
	b.WriteString("NotifyAccess=all\n")
	b.WriteString(fmt.Sprintf("WatchdogSec=%d\n", opts.WatchdogSec))
	b.WriteString(fmt.Sprintf("TimeoutStartSec=%d\n", opts.StartTimeoutSec))
	b.WriteString("ExecStart=" + strings.Join(execParts, " ") + "\n")
	b.WriteString("Restart=on-failure\n")
	b.WriteString(fmt.Sprintf("RestartSec=%d\n", opts.RestartSec))
	if opts.PathEnv != "" {
		b.WriteString(`Environment="PATH=` + opts.PathEnv + "\"\n")
	}
	b.WriteString("\n[Install]\n")
	b.WriteString("WantedBy=default.target\n")
	return b.String()
}

// quoteUnitArg mengutip argumen yang mengandung spasi/quote sesuai aturan
// quoting systemd.
func quoteUnitArg(s string) string {
	if s == "" || strings.ContainsAny(s, " \t\"'\\") {
		r := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
		return `"` + r.Replace(s) + `"`
	}
	return s
}

// SameUnitContent membandingkan dua isi unit dengan mengabaikan baris marker,
// baris kosong, dan trailing whitespace. Dipakai agar install idempotent.
func SameUnitContent(a, b string) bool {
	if a == b {
		return true
	}
	normalize := func(s string) []string {
		var out []string
		for _, line := range strings.Split(s, "\n") {
			line = strings.TrimRight(line, " \t\r")
			if line == "" || strings.HasPrefix(line, unitMarker) {
				continue
			}
			out = append(out, line)
		}
		return out
	}
	x, y := normalize(a), normalize(b)
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// UserUnitDir mengembalikan ~/.config/systemd/user (XDG-aware).
func UserUnitDir() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".config", "systemd", "user")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "systemd", "user")
}

// UserUnitPath mengembalikan path unit untuk nama service.
func UserUnitPath(service string) string {
	return filepath.Join(UserUnitDir(), service+".service")
}

// ExpandHome mengubah prefix "~/" menjadi direktori home user.
func ExpandHome(path string) string {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~"))
}

// ResolveBinary mencari binary di PATH lalu di direktori umum (npm global,
// ~/.local/bin, bun, deno, /usr/local/bin).
func ResolveBinary(name string) (string, error) {
	var candidates []string
	if strings.ContainsRune(name, filepath.Separator) {
		candidates = append(candidates, name)
	} else {
		if found, err := exec.LookPath(name); err == nil {
			candidates = append(candidates, found)
		}
		if home, err := os.UserHomeDir(); err == nil {
			for _, dir := range []string{
				filepath.Join(home, ".npm-global", "bin"),
				filepath.Join(home, ".local", "bin"),
				filepath.Join(home, ".bun", "bin"),
				filepath.Join(home, ".deno", "bin"),
			} {
				candidates = append(candidates, filepath.Join(dir, name))
			}
		}
		candidates = append(candidates, filepath.Join("/usr/local/bin", name))
	}

	for _, candidate := range candidates {
		if !isExecutable(candidate) {
			continue
		}
		if abs, err := filepath.Abs(candidate); err == nil {
			return abs, nil
		}
		return candidate, nil
	}
	return "", fmt.Errorf("%s tidak ditemukan di PATH maupun direktori umum", name)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode().Perm()&0o111 != 0
}

// WriteUserUnit menulis unit user. Jika file sudah ada dan isinya berbeda,
// file lama dibackup dulu. Mengembalikan changed=false bila isi identik.
func WriteUserUnit(service, content string) (changed bool, backup string, err error) {
	path := UserUnitPath(service)
	existing, readErr := os.ReadFile(path)
	if readErr == nil {
		if SameUnitContent(string(existing), content) {
			return false, "", nil
		}
		backup = path + ".bak." + strconv.FormatInt(time.Now().Unix(), 10)
		if err = os.WriteFile(backup, existing, 0o644); err != nil {
			return false, "", err
		}
	} else if !errors.Is(readErr, os.ErrNotExist) {
		return false, "", readErr
	}

	if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, backup, err
	}
	if err = os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, backup, err
	}
	return true, backup, nil
}

// RemoveUserUnit menghapus unit user. Unit yang tidak bertanda archstack
// dibackup dulu agar tidak hilang tanpa jejak.
func RemoveUserUnit(service string) (backup string, err error) {
	path := UserUnitPath(service)
	content, readErr := os.ReadFile(path)
	if errors.Is(readErr, os.ErrNotExist) {
		return "", nil
	}
	if readErr != nil {
		return "", readErr
	}
	if !strings.Contains(string(content), unitMarker) {
		backup = path + ".bak." + strconv.FormatInt(time.Now().Unix(), 10)
		if err := os.WriteFile(backup, content, 0o644); err != nil {
			return "", err
		}
	}
	if err := os.Remove(path); err != nil {
		return backup, err
	}
	return backup, nil
}

// UserDaemonReload menjalankan `systemctl --user daemon-reload`.
func UserDaemonReload() error {
	return exec.Command("systemctl", "--user", "daemon-reload").Run()
}

// UserServiceStatus membaca status service user.
func UserServiceStatus(service string) (ServiceState, error) {
	if service == "" {
		return ServiceNone, nil
	}
	out, err := exec.Command("systemctl", "--user", "is-active", service).Output()
	if err != nil {
		return ServiceStopped, nil
	}
	if strings.TrimSpace(string(out)) == "active" {
		return ServiceRunning, nil
	}
	return ServiceStopped, nil
}

// StartUserService mengaktifkan (enable) sekaligus menyalakan service user.
func StartUserService(service string) error {
	cmd := exec.Command("systemctl", "--user", "enable", "--now", service)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StopUserService menghentikan service user tanpa menonaktifkan autostart.
func StopUserService(service string) error {
	cmd := exec.Command("systemctl", "--user", "stop", service)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RestartUserService me-restart service user (dipakai saat unit berubah).
func RestartUserService(service string) error {
	cmd := exec.Command("systemctl", "--user", "restart", service)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DisableUserService mematikan sekaligus menonaktifkan autostart service user.
func DisableUserService(service string) error {
	cmd := exec.Command("systemctl", "--user", "disable", "--now", service)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UserUnitExists melaporkan apakah file unit ada.
func UserUnitExists(service string) bool {
	_, err := os.Stat(UserUnitPath(service))
	return err == nil
}

// UserServiceMainPID mengembalikan MainPID unit user (0 jika tidak ada/bukan
// dikelola systemd).
func UserServiceMainPID(service string) int {
	out, err := exec.Command("systemctl", "--user", "show", "-p", "MainPID", "--value", service).Output()
	if err != nil {
		return 0
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil || pid <= 0 {
		return 0
	}
	return pid
}

// CurrentUsername mengembalikan nama user untuk loginctl.
func CurrentUsername() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	if u := os.Getenv("LOGNAME"); u != "" {
		return u
	}
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return ""
}

// LingerEnabled membaca status linger user saat ini.
func LingerEnabled() (bool, error) {
	name := CurrentUsername()
	if name == "" {
		return false, errors.New("nama user tidak terdeteksi")
	}
	out, err := exec.Command("loginctl", "show-user", name, "-p", "Linger").Output()
	if err != nil {
		return false, err
	}
	return parseLinger(string(out)), nil
}

func parseLinger(out string) bool {
	return strings.Contains(out, "Linger=yes")
}

// EnableLinger mengaktifkan linger agar user service tetap jalan setelah boot
// tanpa perlu login. Tidak butuh sudo untuk user sendiri di sesi lokal.
func EnableLinger() error {
	name := CurrentUsername()
	if name == "" {
		return errors.New("nama user tidak terdeteksi")
	}
	return exec.Command("loginctl", "enable-linger", name).Run()
}

// ReadPidFile membaca PID dari file pid (prefix ~/ didukung).
func ReadPidFile(path string) (int, bool) {
	data, err := os.ReadFile(ExpandHome(path))
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

// ProcessAlive melaporkan apakah PID masih hidup.
func ProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// WaitHealthy mem-poll URL sampai status 2xx atau timeout habis.
func WaitHealthy(url string, timeout time.Duration) error {
	if url == "" {
		return nil
	}
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(500 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = errors.New("timeout")
	}
	return lastErr
}
