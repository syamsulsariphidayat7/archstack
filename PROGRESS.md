# PROGRESS.md — ArchStack

## Fase Proyek
Development (fitur iterasi) — release v0.2.0 sudah rilis, menuju v0.3.0

## Status Terakhir
- Release v0.1.0: initial CLI (list/install/uninstall/run + TUI + 12 tool)
- Release v0.2.0 (2 Agu 2026): + command `version`, `status [tool]`, `search <query>`, `stop <tool>`, `upgrade`; registry 28 tool; version via ldflags; fix TUI berat (snapshot state, ganti query pacman per frame). Tag terpush, release + asset `amd64`/`arm64` sudah live di GitHub.
- `install.sh` one-liner divalidasi: pola nama asset (`archstack-linux-$GOARCH`) cocok dengan asset release v0.2.0.
- Sesion ini:
  - **Prompt nginx dihapus** (keputusan user: "nginx sama dengan yang lain aja") → `install nginx` dan TUI sekarang langsung install tanpa tanya.
  - **Cleanup dead code machinery prompt/answers**: hapus `internal/registry/answers.go`, `internal/tui/prompt.go`, tipe `Prompt`/`PromptType` + field `Prompts` di registry, `collectAnswers`/`printAnswers` di `install.go`, serta `statePrompts`/`answers`/`promptModel`/`updatePrompts` di TUI. Semua tool kini tanpa prompt. File `~/.config/archstack/answers.json` di mesin user jadi tidak terpakai (aman dihapus manual).

## Langkah Selanjutnya
- Backup (opsional): unit test + test jalan di CI (belum ada test sama sekali; `go test ./...` saat ini kosong)

## Blocker/Catatan
- `sudo` butuh password (no tty) → operasi root dijalankan user manual
- mongodb dari AUR (`mongodb-bin`) via yay
- deno/bun/rabbitmq sudah di repo `extra` (pakai pacman)
