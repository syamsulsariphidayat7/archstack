# PROGRESS.md — ArchStack

## Fase Proyek
Development (fitur iterasi) — menuju release v0.2.0

## Status Terakhir
- Release v0.1.0 sudah rilis (initial CLI: list/install/uninstall/run + TUI + 12 tool)
- Sesion ini selesai nambah:
  - Command baru: `version`, `status [tool]`, `search <query>`, `stop <tool>`, `upgrade`
  - Registry: +16 tool (python, python-pip, go, rust, deno, bun, java, sqlite, certbot, memcached, mongodb, rabbitmq, fail2ban, ffmpeg, imagemagick, tmux) → total 28
  - Version plumbing via ldflags (internal/cli.Version)
  - README diperbarui

## Langkah Selanjutnya
- Commit fitur + tag `v0.2.0` + push (trigger release workflow)
- Validasi `install.sh` one-liner untuk v0.2.0
- Backup (opsional): fix bug answers nginx (version select belum dipakai, project_root belum generate config)
- Backup (opsional): unit test + test jalan di CI

## Blocker/Catatan
- `sudo` butuh password (no tty) → operasi root dijalankan user manual
- mongodb dari AUR (`mongodb-bin`) via yay
- deno/bun/rabbitmq sudah di repo `extra` (pakai pacman)
