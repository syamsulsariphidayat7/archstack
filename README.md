# ArchStack

CLI untuk mengelola development tools di Arch Linux. Mode TUI interaktif + subcommand langsung.

## Instalasi

**One-liner:**
```bash
curl -fsSL https://raw.githubusercontent.com/syamsulsariphidayat7/archstack/main/scripts/install.sh | bash
```

**Manual (Go toolchain):**
```bash
go install github.com/syamsulsariphidayat7/archstack@latest
```

## Penggunaan

**Mode interaktif (TUI):**
```bash
archstack
```
Navigasi panah ↑/↓, Enter untuk aksi, q/Esc keluar.

**Mode command langsung:**
```bash
archstack list                    # Daftar semua tool
archstack status                  # Status semua tool
archstack status nginx            # Status tool tertentu
archstack search database         # Cari tool
archstack install nginx           # Instal tool
archstack install php node redis  # Instal beberapa tool sekaligus
archstack uninstall nginx         # Hapus tool
archstack run postgres            # Jalankan service tool
archstack stop postgres           # Hentikan service tool
archstack install omniroute       # Daftarkan autostart service (binary sudah ada)
archstack status omniroute        # Detail service + linger + health
archstack upgrade                 # Perbarui archstack
archstack version                 # Tampilkan versi
```

## Menambah Tool Baru

Edit `internal/registry/registry.go`, tambah entri ke `defaultTools`:

```go
{
    Name: "namatool", Pkg: "namapackage", From: SourcePacman,
    Binary: "namabinary", Service: "namaservice",
    Desc: "Deskripsi tool",
},
```

Semua tool langsung install tanpa tanya.

### Tool service-only

Untuk tool yang binary-nya dipasang di luar archstack (npm global, dsb) dan
hanya perlu autostart service, gunakan `SourceServiceOnly` + user service:

```go
{
    Name: "omniroute", From: SourceServiceOnly, Binary: "omniroute",
    Service: "omniroute", ServiceScope: ServiceUser, Linger: true,
    PidFile:   "~/.omniroute/server/.pid",
    HealthURL: "http://127.0.0.1:20128/api/health",
    Desc:      "OmniRoute AI proxy router — autostart service (binary manual)",
},
```

- `install` → tulis `~/.config/systemd/user/<service>.service` (backup bila sudah
  ada dan isinya beda), `enable --now`, `enable-linger`, lalu health check.
- `uninstall` → `disable --now` + hapus unit; **binary dan data tidak disentuh**.
- `PidFile` mendeteksi instance manual yang masih memakai port yang sama;
  install ditolak dengan instruksi `<binary> stop`.
- Binary dicari di `PATH` lalu `~/.npm-global/bin`, `~/.local/bin`, `~/.bun/bin`,
  `~/.deno/bin`, `/usr/local/bin`.
- Unit selalu foreground (`Type=notify`) — jangan pakai mode daemon internal.

## Build Manual

```bash
git clone https://github.com/syamsulsariphidayat7/archstack.git
cd archstack
go build -ldflags "-X github.com/syamsulsariphidayat7/archstack/internal/cli.Version=v0.5.0" -o archstack .
sudo cp archstack /usr/local/bin/
```

## Lisensi

MIT
