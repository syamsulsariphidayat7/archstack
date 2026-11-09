package cli

import (
	"fmt"
	"net"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/web"
)

// webCmd menjalankan Web UI lokal. Bind SELALU loopback: UI ini setara
// akses root (install/hapus paket, nyala/matikan service), jadi alamat
// non-lokal ditolak keras, bukan sekadar diperingatkan.
func webCmd(args []string) error {
	addr := "127.0.0.1:8090"
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--addr" && i+1 < len(args):
			i++
			addr = args[i]
		case strings.HasPrefix(a, "--addr="):
			addr = strings.TrimPrefix(a, "--addr=")
		case a == "-h" || a == "--help":
			fmt.Println("Penggunaan: archstack web [--addr 127.0.0.1:8090]")
			fmt.Println("Web UI lokal untuk mengelola tool. Hanya alamat loopback yang diterima.")
			return nil
		default:
			return fmt.Errorf("[error] Opsi tidak dikenal untuk web: %s", a)
		}
	}
	if err := validateLoopbackAddr(addr); err != nil {
		return err
	}
	return web.Run(addr)
}

// validateLoopbackAddr memastikan addr berbentuk host:port dengan host
// loopback (127.0.0.0/8, ::1, atau "localhost").
func validateLoopbackAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("[error] Alamat web tidak valid %q (bentuk: host:port): %v", addr, err)
	}
	if port == "" {
		return fmt.Errorf("[error] Port kosong pada alamat %q", addr)
	}
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("[error] archstack web hanya boleh bind ke localhost (127.0.0.1 / ::1), bukan %q — UI ini setara akses root", host)
	}
	return nil
}
