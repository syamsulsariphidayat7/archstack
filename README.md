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
    Prompts: []Prompt{
        {Key: "version", Question: "Pilih versi", Type: PromptSelect,
         Options: []string{"opsi1", "opsi2"}, Recommended: "opsi1"},
    },
},
```

Tool tanpa Prompts (kosong) langsung install tanpa tanya.

## Build Manual

```bash
git clone https://github.com/syamsulsariphidayat7/archstack.git
cd archstack
go build -ldflags "-X github.com/syamsulsariphidayat7/archstack/internal/cli.Version=0.1.0" -o archstack .
sudo cp archstack /usr/local/bin/
```

## Lisensi

MIT
