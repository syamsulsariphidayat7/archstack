package web

import (
	"os"
	"testing"
	"time"
)

func TestNotifyAddr(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", ""},
		{"/run/user/1000/systemd/notify", "/run/user/1000/systemd/notify"},
		{"@/run/systemd/notify", "\x00/run/systemd/notify"},
	}
	for _, c := range cases {
		if got := notifyAddr(c.in); got != c.want {
			t.Errorf("notifyAddr(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestInitNotifyMenangkapDanMembersihkanEnv(t *testing.T) {
	t.Setenv("NOTIFY_SOCKET", "/tmp/notify-test.sock")
	t.Setenv("WATCHDOG_USEC", "4000000")
	t.Setenv("WATCHDOG_PID", "")
	defer func() { notifyTarget = "" }()

	interval, ok := initNotify()
	if !ok || interval != 2*time.Second {
		t.Errorf("initNotify interval = %v,%v — want 2s,true", interval, ok)
	}
	if notifyTarget != "/tmp/notify-test.sock" {
		t.Errorf("notifyTarget = %q", notifyTarget)
	}
	for _, key := range []string{"NOTIFY_SOCKET", "WATCHDOG_USEC", "WATCHDOG_PID"} {
		if _, ada := os.LookupEnv(key); ada {
			t.Errorf("%s harus sudah dihapus dari environment agar tidak diwarisi child process", key)
		}
	}
}

func TestParseWatchdogInterval(t *testing.T) {
	d, ok := parseWatchdogInterval("180000000") // WatchdogSec=180 bawaan RenderUnit
	if !ok || d != 90*time.Second {
		t.Errorf("parseWatchdogInterval(180s) = %v,%v — want 90s,true", d, ok)
	}
	for _, bad := range []string{"", "0", "-5", "abc"} {
		if _, ok := parseWatchdogInterval(bad); ok {
			t.Errorf("parseWatchdogInterval(%q) harus false", bad)
		}
	}
}
