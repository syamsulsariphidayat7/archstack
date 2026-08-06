# AGENTS.md — ArchStack

CLI (Go) untuk kelola development tools di Arch Linux. TUI (bubbletea) + subcommand.

## Perintah
- `go build ./...` — build
- `go vet ./...` — lint
- `go test ./...` — test (belum ada)
- `go build -ldflags "-X github.com/syamsulsariphidayat7/archstack/internal/cli.Version=v0.1.0" -o archstack .` — build dengan versi

## Struktur
- `main.go` → entry, panggil `cli.Root`
- `internal/cli/` — subcommand: root, install, uninstall, run, stop, status, search, version, upgrade
- `internal/registry/` — `defaultTools` (data tool)
- `internal/system/` — pacman/yay/systemd wrapper
- `internal/tui/` — bubbletea UI
- `scripts/install.sh` — one-liner dari GitHub release
- `.github/workflows/release.yml` — build amd64/arm64 + upload saat tag `v*`

## Konvensi
- Tambah tool: edit `defaultTools` di `internal/registry/registry.go`
- Tool dengan `Service` otomatis dapat aksi start/stop di TUI
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
