# PROGRESS.md — ArchStack

## Fase Proyek
Development (fitur iterasi) — release v0.3.0 sudah rilis, menuju v0.4.0

## Status Terakhir
- Release v0.1.0: initial CLI (list/install/uninstall/run + TUI + 12 tool)
- Release v0.2.0 (2 Agu 2026): + command `version`, `status [tool]`, `search <query>`, `stop <tool>`, `upgrade`; registry 28 tool; version via ldflags; fix TUI berat (snapshot state, ganti query pacman per frame). Tag terpush, release + asset `amd64`/`arm64` sudah live di GitHub.
- Release v0.3.0: cleanup dead code machinery prompt/answers — nginx kini langsung install tanpa tanya (semua tool tanpa prompt), -521 baris. Tag terpush, release + asset `amd64`/`arm64` sudah live di GitHub.
- `install.sh` one-liner divalidasi: pola nama asset (`archstack-linux-$GOARCH`) cocok dengan asset release v0.2.0.
- Sesion ini:
  - **Prompt nginx dihapus** (keputusan user: "nginx sama dengan yang lain aja") → `install nginx` dan TUI sekarang langsung install tanpa tanya.
  - **Cleanup dead code machinery prompt/answers**: hapus `internal/registry/answers.go`, `internal/tui/prompt.go`, tipe `Prompt`/`PromptType` + field `Prompts` di registry, `collectAnswers`/`printAnswers` di `install.go`, serta `statePrompts`/`answers`/`promptModel`/`updatePrompts` di TUI. Semua tool kini tanpa prompt.
  - **Rilis v0.3.0 + update binary terinstal**: tag v0.3.0 terpush, release + asset live, `/usr/local/bin/archstack` terverifikasi v0.3.0 (byte-identik dengan asset release).
  - **Unit test + CI**: tambah test untuk registry/upgrade/search/padRight (murni, tanpa sudo/pacman/systemctl/network) + workflow `.github/workflows/ci.yml` (gofmt/vet/test/build di push/PR). Refactor kecil: ekstrak `isUpToDate` (upgrade.go) dan `searchTools` (search.go) agar testable.
  - **Rekomendasi versi saat install** (dinamis): `AvailableVersion` query `pacman -Si`/`yay -Si`, CLI install & TUI menampilkan "versi rekomendasi" sebelum install. Helper murni `parseVersion` di-cover unit test.
  - **Kolom VERSI di status/list**: `pacman -Q` (versi terinstall, 1 call) + `pacman -Sl` (versi repo, 1 call) via `ToolVersionStrings`; terinstall → tampil versi terinstall, selain itu versi rekomendasi. `installedPkgs` di TUI pakai `InstalledVersions` (hemat 1 subprocess).

## Langkah Selanjutnya
- Semua item prioritas selesai. Ide berikutnya (opsional):
  - Integration test end-to-end (mis. mock `exec.Command` untuk install/status)
  - Changelog otomatis dari release notes GitHub
  - Dukungan distro lain (apt/dnf) — perlu refactor `internal/system/`

## Blocker/Catatan
- `sudo` butuh password (no tty) → operasi root dijalankan user manual
- mongodb dari AUR (`mongodb-bin`) via yay
- deno/bun/rabbitmq sudah di repo `extra` (pakai pacman)
