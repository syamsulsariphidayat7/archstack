# PROGRESS.md — ArchStack

## Fase Proyek
Development (fitur iterasi) — release v0.5.0 sudah rilis (tool service-only omniroute)

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
- Release v0.4.0 & v0.5.0: **tool service-only pertama: `omniroute`**:
  - Keputusan user: archstack **hanya mengurus service**, binary/npm tidak dipasang dan tidak dihapus (`SourceServiceOnly`). Uninstall = service saja; unit yang sudah ada diadopsi (backup bila isi beda).
  - `registry.Tool`: tambah `ServiceScope` (system/user), `Linger`, `PidFile`, `HealthURL`, `SourceServiceOnly`; entry `omniroute` = user service + linger + pid file `~/.omniroute/server/.pid` + health `:20128`.
  - `internal/system/userunit.go`: `ResolveBinary` (PATH + `~/.npm-global/bin`, `~/.local/bin`, bun, deno, `/usr/local/bin`), `RenderUnit` (`Type=notify`, watchdog 180s, foreground `serve --no-open`), `WriteUserUnit` (backup bila isi beda), `RemoveUserUnit` (backup unit tanpa marker), `systemctl --user` (daemon-reload/restart), linger, `WaitHealthy`.
  - `internal/cli/service.go`: `install` = cek binary → tolak bila ada instance manual → tulis/adopsi unit → daemon-reload → restart bila unit berubah → `enable --now` → linger → health check 30s. `uninstall` = `disable --now` + hapus unit, binary & `~/.omniroute/` utuh. `status <tool>` menampilkan binary/unit/pid/linger/health.
  - `systemd.go`: helper scope-aware (`ServiceStatusFor`, `StartServiceFor`, `StopServiceFor`, command untuk TUI). TUI memanggil ulang binary sendiri (`execSelf`) untuk install/uninstall service-only, jadi logika unit tidak diduplikasi.
  - Unit test murni: render unit (termasuk larangan `--daemon`), idempotensi tulis unit, backup unit asing, `parseLinger`, `ExpandHome`, `ResolveBinary`.
  - Validasi live: adopsi unit manual → restart + health OK; install kedua idempotent (tanpa restart); stop/run; uninstall (unit hilang, binary + `storage.sqlite` utuh); reinstall; penolakan saat daemon manual masih jalan.
  - **Bug yang ditemukan saat uji live**: `install` sempat short-circuit "sudah terinstal" karena binary ada (padahal service belum tentu terdaftar) → dilewati khusus service-only; dan `PidFile` berisi PID child, bukan MainPID → cek konflik hanya dilakukan saat service tidak aktif.

## Langkah Selanjutnya
- **v0.6.2 — Aksi Setup satu kali**: field registry `SetupDir`/`SetupMarker`/`SetupUser`/`SetupCmd`/`SetupNote` + helper `NeedsSetup()`; hanya diisi postgres (`initdb --locale en_US.UTF-8 -E UTF8 -D /var/lib/postgres/data` sebagai user postgres, penanda `PG_VERSION`) dan mariadb (`mariadb-install-db --user=mysql ...`, penanda subdir `mysql`). `system.SetupDone` mengecek penanda (idempotensi); runner `RunSetup` (CLI, sudo tty) + `RunSetupW` (web, `sudo -S`/`-n`) di atas builder bersama `sudoCmdOpts` (execw.go kini mendukung `sudo -u`). Subcommand `archstack setup <tool>` + tombol Setup & badge PERLU SETUP di Web UI; `run` (CLI + web) menolak start dan mengarahkan ke setup bila penanda belum ada. SetupNote memberi langkah manual berikutnya (user/database/password TIDAK dikelola archstack). Unit test murni: field registry, `SetupDone` di tempdir, argv `sudoCmdOpts`, validasi aksi + needsSudo web.
- **v0.6.1 — Web UI jadi service sendiri**: entri registry `archstack-web` (tool service-only kedua setelah omniroute): `Binary: archstack`, user service + linger, `UnitArgs: [web]`, `HealthURL: http://127.0.0.1:8090/api/tools`. `archstack install archstack-web` mendaftarkan autostart; kartu `archstack-web` muncul di dashboard/TUI/CLI. Syaratnya dipenuhi di `internal/web/notify.go`: server kini listen dulu lalu `sd_notify READY=1` + ping `WATCHDOG=1` setengah interval `WATCHDOG_USEC` (unit RenderUnit selalu `Type=notify` + WatchdogSec 180) — no-op penuh bila dijalankan dari terminal (tanpa `NOTIFY_SOCKET`). Temuan saat uji live VM: `systemctl` mengirim datagram `EXIT_STATUS` ke `NOTIFY_SOCKET` yang diwarisinya, jadi `initNotify` menangkap socket + interval lalu MENGHAPUS ketiga variabel itu dari environment proses agar child (pacman/systemctl dari handler) tidak mengotori kanal notifikasi. Unit test murni: entri registry archstack-web, `notifyAddr`, `parseWatchdogInterval`. Efek samping yang didokumentasikan: Stop dari dashboard mematikan dashboard sendiri; nyalakan lagi dari terminal (`archstack run archstack-web`). Didahului perubahan gaya: lebar halaman web full width (commit terpisah, tanpa tag).
- **Web UI (v0.6.0, sudah rilis)**: subcommand `archstack web` — server HTTP lokal (bind loopback saja, ditolak keras bila non-lokal) + UI satu halaman ter-embed (`internal/web/static/index.html`, tanpa dependensi baru). API: `GET /api/tools`, `POST /api/action`, `GET /api/job`. Layout: daftar tool di kiri, **log proses permanen di sidebar kanan** (sticky; di layar sempit jadi panel nempel di bawah) — keputusan user, bukan modal. Modal hanya untuk password sudo (dialirkan ke `sudo -S`, tidak disimpan; kosong = andalkan cache sudo via `sudo -n`). Aksi satu per satu (409 bila sibuk) karena pacman mengunci database. Fungsi streaming baru di `internal/system/execw.go` (`InstallPacmanW`/`InstallYayW`/`UninstallW`/`StartServiceForW`/`StopServiceForW` dengan `io.Writer`); aksi service-only memakai ulang CLI lewat self-exec. Unit test murni: validasi aksi, `needsSudo`, job manager + line writer, validasi alamat loopback.
- v0.6.0 (opsional): fitur baru atau perbaikan berikutnya.
- Ide berikutnya (opsional):
  - Integration test end-to-end (mis. mock `exec.Command` untuk install/status)
  - Changelog otomatis dari release notes GitHub
  - Dukungan distro lain (apt/dnf) — perlu refactor `internal/system/`

## Blocker/Catatan
- `sudo` butuh password (no tty) → operasi root dijalankan user manual
- **Go toolchain tidak terpasang di mesin ini** (2026-10-19): verifikasi gofmt/vet/test/build memakai Go portable 1.26.9 di `/tmp/opencode/gotool` (unduhan resmi go.dev). Pasang `go` via pacman bila ingin build permanen.
- mongodb dari AUR (`mongodb-bin`) via yay
- deno/bun/rabbitmq sudah di repo `extra` (pakai pacman)
