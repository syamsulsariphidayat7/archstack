# PROGRESS.md — ArchStack

## Fase Proyek
Development (fitur iterasi) — release v0.3.0 sudah rilis, v0.4.0 (tool service-only) siap dirilis

## Status Terakhir
- Release v0.1.0: initial CLI (list/install/uninstall/run + TUI + 12 tool)
- Release v0.2.0 (2 Agu 2026): + command `version`, `status [tool]`, `search <query>`, `stop <tool>`, `upgrade`; registry 28 tool; version via ldflags; fix TUI berat (snapshot state, ganti query pacman per frame). Tag terpush, release + asset `amd64`/`arm64` sudah live di GitHub.
- Release v0.3.0: cleanup dead code machinery prompt/answers — nginx kini langsung install tanpa tanya (semua tool tanpa prompt), -521 baris. Tag terpush, release + asset `amd64`/`arm64` sudah live di GitHub.
- `install.sh` one-liner divalidasi: pola nama asset (`archstack-linux-$GOARCH`) cocok dengan asset release v0.2.0.
- Sesion v0.3.0:
  - **Prompt nginx dihapus** (keputusan user: "nginx sama dengan yang lain aja") → `install nginx` dan TUI sekarang langsung install tanpa tanya.
  - **Cleanup dead code machinery prompt/answers**: hapus `internal/registry/answers.go`, `internal/tui/prompt.go`, tipe `Prompt`/`PromptType` + field `Prompts` di registry, `collectAnswers`/`printAnswers` di `install.go`, serta `statePrompts`/`answers`/`promptModel`/`updatePrompts` di TUI. Semua tool kini tanpa prompt.
  - **Rilis v0.3.0 + update binary terinstal**: tag v0.3.0 terpush, release + asset live, `/usr/local/bin/archstack` terverifikasi v0.3.0 (byte-identik dengan asset release).
  - **Unit test + CI**: tambah test untuk registry/upgrade/search/padRight (murni, tanpa sudo/pacman/systemctl/network) + workflow `.github/workflows/ci.yml` (gofmt/vet/test/build di push/PR). Refactor kecil: ekstrak `isUpToDate` (upgrade.go) dan `searchTools` (search.go) agar testable.
  - **Rekomendasi versi saat install** (dinamis): `AvailableVersion` query `pacman -Si`/`yay -Si`, CLI install & TUI menampilkan "versi rekomendasi" sebelum install. Helper murni `parseVersion` di-cover unit test.
  - **Kolom VERSI di status/list**: `pacman -Q` (versi terinstall, 1 call) + `pacman -Sl` (versi repo, 1 call) via `ToolVersionStrings`; terinstall → tampil versi terinstall, selain itu versi rekomendasi. `installedPkgs` di TUI pakai `InstalledVersions` (hemat 1 subprocess).
- Sesion ini (2026-10-19) — **tool service-only pertama: `omniroute`**:
  - Keputusan user: archstack **hanya mengurus service**, binary/npm tidak dipasang dan tidak dihapus (`SourceServiceOnly`). Uninstall = service saja; unit yang sudah ada diadopsi (backup bila isi beda).
  - `registry.Tool`: tambah `ServiceScope` (system/user), `Linger`, `PidFile`, `HealthURL`, `SourceServiceOnly`; entry `omniroute` = user service + linger + pid file `~/.omniroute/server/.pid` + health `:20128`.
  - `internal/system/userunit.go`: `ResolveBinary` (PATH + `~/.npm-global/bin`, `~/.local/bin`, bun, deno, `/usr/local/bin`), `RenderUnit` (`Type=notify`, watchdog 180s, foreground `serve --no-open`), `WriteUserUnit` (backup bila isi beda), `RemoveUserUnit` (backup unit tanpa marker), `systemctl --user` (daemon-reload/restart), linger, `WaitHealthy`.
  - `internal/cli/service.go`: `install` = cek binary → tolak bila ada instance manual → tulis/adopsi unit → daemon-reload → restart bila unit berubah → `enable --now` → linger → health check 30s. `uninstall` = `disable --now` + hapus unit, binary & `~/.omniroute/` utuh. `status <tool>` menampilkan binary/unit/pid/linger/health.
  - `systemd.go`: helper scope-aware (`ServiceStatusFor`, `StartServiceFor`, `StopServiceFor`, command untuk TUI). TUI memanggil ulang binary sendiri (`execSelf`) untuk install/uninstall service-only, jadi logika unit tidak diduplikasi.
  - Unit test murni: render unit (termasuk larangan `--daemon`), idempotensi tulis unit, backup unit asing, `parseLinger`, `ExpandHome`, `ResolveBinary`.
  - Validasi live: adopsi unit manual → restart + health OK; install kedua idempotent (tanpa restart); stop/run; uninstall (unit hilang, binary + `storage.sqlite` utuh); reinstall; penolakan saat daemon manual masih jalan.
  - **Bug yang ditemukan saat uji live**: `install` sempat short-circuit "sudah terinstal" karena binary ada (padahal service belum tentu terdaftar) → dilewati khusus service-only; dan `PidFile` berisi PID child, bukan MainPID → cek konflik hanya dilakukan saat service tidak aktif.

## Langkah Selanjutnya
- Rilis v0.4.0: bump tag, CI build amd64/arm64, update binary `/usr/local/bin`.
- Ide berikutnya (opsional):
  - Integration test end-to-end (mis. mock `exec.Command` untuk install/status)
  - Changelog otomatis dari release notes GitHub
  - Dukungan distro lain (apt/dnf) — perlu refactor `internal/system/`

## Blocker/Catatan
- `sudo` butuh password (no tty) → operasi root dijalankan user manual
- **Go toolchain tidak terpasang di mesin ini** (2026-10-19): verifikasi gofmt/vet/test/build memakai Go portable 1.26.9 di `/tmp/opencode/gotool` (unduhan resmi go.dev). Pasang `go` via pacman bila ingin build permanen.
- mongodb dari AUR (`mongodb-bin`) via yay
- deno/bun/rabbitmq sudah di repo `extra` (pakai pacman)
