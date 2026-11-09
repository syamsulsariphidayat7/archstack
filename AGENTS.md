# AGENTS.md — ArchStack

CLI (Go) untuk kelola development tools di Arch Linux. TUI (bubbletea) + subcommand.

## Perintah
- `gofmt -l .` — cek format (harus kosong)
- `go build ./...` — build
- `go vet ./...` — lint
- `go test ./...` — unit test
- `go build -ldflags "-X github.com/syamsulsariphidayat7/archstack/internal/cli.Version=v0.5.0" -o archstack .` — build dengan versi

## Struktur
- `main.go` → entry, panggil `cli.Root`
- `internal/cli/` — subcommand: root, install, uninstall, run, stop, status, search, version, upgrade, web
  - `service.go` — alur install/uninstall tool service-only
  - `web.go` — subcommand `web`: validasi alamat WAJIB loopback, lalu `web.Run`
- `internal/registry/` — `defaultTools` (data tool)
- `internal/system/` — pacman/yay/systemd wrapper
  - `userunit.go` — render/tulis/hapus unit user, `systemctl --user`, linger, health check
  - `execw.go` — varian streaming untuk Web UI: output ke `io.Writer`, password sudo via `sudo -S` (kosong → `sudo -n`)
- `internal/web/` — Web UI lokal: API `/api/tools` `/api/action` `/api/job`, satu job pada satu waktu, UI ter-embed di `static/index.html`; log proses di sidebar kanan (bukan modal — modal hanya untuk password sudo)
- `internal/tui/` — bubbletea UI
- `scripts/install.sh` — one-liner dari GitHub release
- `.github/workflows/ci.yml` — gofmt/vet/test/build otomatis di setiap push/PR
- `.github/workflows/release.yml` — build amd64/arm64 + upload saat tag `v*`

## Konvensi
- Tambah tool: edit `defaultTools` di `internal/registry/registry.go`
- Tool dengan `Service` otomatis dapat aksi start/stop di TUI
- **Tool service-only** (`From: SourceServiceOnly`): archstack hanya mengelola
  unit service, bukan binary. Wajib isi `Service`; `ServiceScope: ServiceUser`
  untuk unit `~/.config/systemd/user` + `Linger` bila perlu hidup tanpa login.
  `PidFile` (opsional, prefix `~/`) mendeteksi instance manual; `HealthURL`
  (opsional) dicek setelah start. Binary dicari lewat `system.ResolveBinary`.
- Service system memakai `sudo systemctl`; user service memakai
  `systemctl --user` (tanpa sudo). Semua lewat helper `system.ServiceStatusFor`,
  `StartServiceFor`, `StopServiceFor`, `StartServiceCommand`, `StopServiceCommand`.
- TUI tidak menduplikasi logika service-only: aksi install/uninstall memanggil
  ulang binary archstack (`execSelf`), sedangkan start/stop memakai command
  helper di `internal/system`.
- Unit user dirender oleh `system.RenderUnit` dan selalu `Type=notify` +
  foreground. Jangan pakai mode daemon internal (keluar cgroup, memutus
  READY/WATCHDOG).
- Unit test: file `*_test.go` di package yang sama; test harus murni (tidak memanggil sudo/pacman/systemctl/network) — ekstrak helper bila perlu. Test unit user memakai `t.Setenv("XDG_CONFIG_HOME", t.TempDir())`.
- Pesan output CLI pakai prefix `[ok]`/`[install]`/`[run]`/`[done]`/`[error]`/`[info]`
- Bahasa output: Indonesia

## Versi Library Kompatibel
| Library | Versi |
|---|---|
| go | 1.26.5 |
| github.com/charmbracelet/bubbletea | v1.3.10 |
| github.com/charmbracelet/bubbles | v1.0.0 |
| github.com/charmbracelet/lipgloss | v1.1.0 |

## Catatan
- Release baru: bump tag `vX.Y.Z`, workflow otomatis build + upload
- `archstack upgrade` self-update dari GitHub latest release
