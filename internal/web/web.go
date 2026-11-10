// Package web menyediakan Web UI lokal untuk archstack: satu server HTTP
// yang hanya di-bind ke localhost oleh pemanggil (internal/cli), dengan
// API JSON kecil di atas internal/registry + internal/system.
//
// Keamanan: siapa pun yang bisa membuka UI ini setara memegang root
// (install/hapus paket, nyala/matikan service), jadi server TIDAK BOLEH
// diekspos ke jaringan. Aksi berjalan satu per satu (antre ditolak) karena
// pacman mengunci database selama proses berjalan.
package web

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

//go:embed static/index.html
var indexHTML []byte

// ---------- View model ----------

type toolView struct {
	Name        string `json:"name"`
	Pkg         string `json:"pkg"`
	From        string `json:"from"`
	Desc        string `json:"desc"`
	Service     string `json:"service"`
	Scope       string `json:"scope"` // "user" | "system"
	ServiceOnly bool   `json:"serviceOnly"`
	Installed   bool   `json:"installed"`
	Svc         string `json:"svc"` // "" | "running" | "stopped"
	Version     string `json:"version"`
	Unit        bool   `json:"unit"` // unit user ada (khusus service-only)
}

func makeView(t registry.Tool, installed bool, svc, version string, unit bool) toolView {
	scope := "system"
	if t.ServiceScope == registry.ServiceUser {
		scope = "user"
	}
	pkg := t.Pkg
	if t.IsServiceOnly() {
		pkg = "— service only —"
	}
	return toolView{
		Name:        t.Name,
		Pkg:         pkg,
		From:        string(t.From),
		Desc:        t.Desc,
		Service:     t.Service,
		Scope:       scope,
		ServiceOnly: t.IsServiceOnly(),
		Installed:   installed,
		Svc:         svc,
		Version:     version,
		Unit:        unit,
	}
}

func collectTools() []toolView {
	tools := registry.AllTools()
	versions := system.ToolVersionStrings(tools)
	out := make([]toolView, 0, len(tools))
	for i := range tools {
		t := tools[i]
		svc := ""
		if t.Service != "" {
			state, _ := system.ServiceStatusFor(&t)
			switch state {
			case system.ServiceRunning:
				svc = "running"
			case system.ServiceStopped:
				svc = "stopped"
			}
		}
		unit := false
		if t.IsServiceOnly() {
			unit = system.UserUnitExists(t.Service)
		}
		out = append(out, makeView(t, system.IsInstalled(t.Binary, t.Pkg), svc, versions[t.Name], unit))
	}
	return out
}

// ---------- Validasi aksi (murni, untuk unit test) ----------

func validateAction(t *registry.Tool, action string) error {
	switch action {
	case "install", "uninstall":
		return nil
	case "run", "stop":
		if t.Service == "" {
			return fmt.Errorf("%s tidak punya service", t.Name)
		}
		return nil
	default:
		return fmt.Errorf("aksi tidak dikenal: %s", action)
	}
}

// needsSudo melaporkan apakah aksi memerlukan sudo. Aksi user-service dan
// install via yay berjalan sebagai user biasa.
func needsSudo(t *registry.Tool, action string) bool {
	if t.IsServiceOnly() {
		return false
	}
	switch action {
	case "install":
		return t.From != registry.SourceYay
	case "uninstall":
		return true
	case "run", "stop":
		return t.ServiceScope != registry.ServiceUser
	}
	return false
}

// ---------- Job: satu proses pada satu waktu, log terkumpul per baris ----------

type jobState struct {
	ID      int      `json:"id"`
	Running bool     `json:"running"`
	Tool    string   `json:"tool"`
	Action  string   `json:"action"`
	Lines   []string `json:"lines"`
	Error   string   `json:"error,omitempty"`
}

type jobManager struct {
	mu      sync.Mutex
	state   jobState
	partial string
}

func (m *jobManager) tryStart(tool, action string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.Running {
		return false
	}
	m.state = jobState{ID: m.state.ID + 1, Running: true, Tool: tool, Action: action, Lines: []string{}}
	m.partial = ""
	return true
}

func (m *jobManager) snapshot() jobState {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := m.state
	out.Lines = append([]string(nil), m.state.Lines...)
	return out
}

func (m *jobManager) finish(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.partial != "" {
		m.state.Lines = append(m.state.Lines, m.partial)
		m.partial = ""
	}
	m.state.Running = false
	if err != nil {
		m.state.Error = err.Error()
	}
}

type lineWriter struct{ m *jobManager }

func (w *lineWriter) Write(p []byte) (int, error) {
	w.m.mu.Lock()
	defer w.m.mu.Unlock()
	s := w.m.partial + string(p)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	parts := strings.Split(s, "\n")
	w.m.partial = parts[len(parts)-1]
	for _, line := range parts[:len(parts)-1] {
		w.m.state.Lines = append(w.m.state.Lines, line)
	}
	return len(p), nil
}

// ---------- Eksekusi aksi ----------

func runSelf(w io.Writer, args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, args...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func runAction(t *registry.Tool, action, pass string, w io.Writer) error {
	switch action {
	case "install":
		if t.IsServiceOnly() {
			// Logika service-only (unit user, linger, health check) sudah
			// lengkap di CLI dan tanpa sudo — pakai ulang lewat self-exec.
			return runSelf(w, "install", t.Name)
		}
		if system.IsInstalled(t.Binary, t.Pkg) {
			fmt.Fprintf(w, "[ok] %s sudah terinstal\n", t.Name)
			return nil
		}
		if ver := system.AvailableVersion(t.Pkg, t.From == registry.SourceYay); ver != "" {
			fmt.Fprintf(w, "[install] Menginstal %s — versi rekomendasi: %s...\n", t.Name, ver)
		} else {
			fmt.Fprintf(w, "[install] Menginstal %s...\n", t.Name)
		}
		var err error
		if t.From == registry.SourceYay {
			err = system.InstallYayW(t.Pkg, w)
		} else {
			err = system.InstallPacmanW(t.Pkg, pass, w)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "[done] %s berhasil diinstal\n", t.Name)
		return nil

	case "uninstall":
		if t.IsServiceOnly() {
			return runSelf(w, "uninstall", t.Name)
		}
		if !system.IsInstalled(t.Binary, t.Pkg) {
			fmt.Fprintf(w, "[ok] %s sudah tidak terinstal\n", t.Name)
			return nil
		}
		if t.Service != "" {
			if state, _ := system.ServiceStatusFor(t); state == system.ServiceRunning {
				fmt.Fprintf(w, "[info] Menghentikan service %s...\n", t.Service)
				if err := system.StopServiceForW(t, pass, w); err != nil {
					return err
				}
			}
		}
		fmt.Fprintf(w, "[install] Menghapus %s...\n", t.Name)
		if err := system.UninstallW(t.Pkg, pass, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "[done] %s berhasil dihapus\n", t.Name)
		return nil

	case "run":
		fmt.Fprintf(w, "[run] Menjalankan service %s...\n", t.Name)
		if err := system.StartServiceForW(t, pass, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "[ok] Service %s aktif\n", t.Name)
		return nil

	case "stop":
		fmt.Fprintf(w, "[stop] Menghentikan service %s...\n", t.Name)
		if err := system.StopServiceForW(t, pass, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "[ok] Service %s berhenti\n", t.Name)
		return nil
	}
	return fmt.Errorf("aksi tidak dikenal: %s", action)
}

// ---------- HTTP server ----------

type Server struct {
	jobs *jobManager
	mux  *http.ServeMux
}

func NewServer() *Server {
	s := &Server{jobs: &jobManager{}, mux: http.NewServeMux()}
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/tools", s.handleTools)
	s.mux.HandleFunc("/api/action", s.handleAction)
	s.mux.HandleFunc("/api/job", s.handleJob)
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

// Run menjalankan server di addr (pemanggil wajib memastikan addr loopback).
// Setelah listener siap, server melapor READY ke systemd (sd_notify) dan
// menjalankan ping watchdog — perlu bila dijalankan sebagai user service
// `archstack-web` yang unitnya bertipe notify.
func Run(addr string) error {
	fmt.Printf("[info] ArchStack Web berjalan di http://%s\n", addr)
	fmt.Printf("[info] Hanya untuk lokal — jangan diekspos ke jaringan. Ctrl+C untuk berhenti.\n")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	interval, watchdogOK := initNotify()
	sdNotify("READY=1")
	if watchdogOK {
		stop := make(chan struct{})
		defer close(stop)
		go runWatchdog(interval, stop)
	}
	return http.Serve(ln, NewServer().Handler())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func (s *Server) handleTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, collectTools())
}

type actionRequest struct {
	Tool     string `json:"tool"`
	Action   string `json:"action"`
	Password string `json:"password"`
}

func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req actionRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body bukan JSON valid"})
		return
	}
	tool, ok := registry.GetTool(req.Tool)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Tool tidak dikenal: " + req.Tool})
		return
	}
	if err := validateAction(tool, req.Action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !s.jobs.tryStart(tool.Name, req.Action) {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "Proses lain sedang berjalan — tunggu sampai selesai (pacman mengunci database).",
		})
		return
	}
	pass := req.Password // hanya hidup di goroutine ini, tidak disimpan/log
	go func() {
		err := runAction(tool, req.Action, pass, &lineWriter{m: s.jobs})
		s.jobs.finish(err)
	}()
	writeJSON(w, http.StatusAccepted, map[string]bool{"ok": true})
}

func (s *Server) handleJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, s.jobs.snapshot())
}
