package cli

import "testing"

func TestValidateLoopbackAddr(t *testing.T) {
	ok := []string{"127.0.0.1:8090", "localhost:8090", "[::1]:8090", "127.0.0.53:9000"}
	for _, a := range ok {
		if err := validateLoopbackAddr(a); err != nil {
			t.Errorf("validateLoopbackAddr(%q) = %v, harus lolos", a, err)
		}
	}
	bad := []string{"0.0.0.0:8090", "192.168.1.10:8090", ":8090", "example.com:8090", "127.0.0.1", "127.0.0.1:"}
	for _, a := range bad {
		if err := validateLoopbackAddr(a); err == nil {
			t.Errorf("validateLoopbackAddr(%q) harus ditolak", a)
		}
	}
}
