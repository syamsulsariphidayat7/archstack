package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const repoOwner = "syamsulsariphidayat7"
const repoName = "archstack"

func upgradeCmd() error {
	fmt.Printf("[info] Versi saat ini: %s\n", Version)

	latest, err := fetchLatestTag()
	if err != nil {
		return fmt.Errorf("[error] Gagal cek rilis terbaru: %w", err)
	}
	if strings.TrimPrefix(latest, "v") == strings.TrimPrefix(Version, "v") {
		fmt.Printf("[ok] Sudah versi terbaru: %s\n", Version)
		return nil
	}
	fmt.Printf("[info] Versi terbaru: %s\n", latest)

	goarch, err := goarch()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s-linux-%s", repoOwner, repoName, latest, repoName, goarch)

	tmp, err := os.MkdirTemp("", "archstack-upgrade-")
	if err != nil {
		return fmt.Errorf("[error] Gagal buat temp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	binPath := filepath.Join(tmp, repoName)
	fmt.Printf("[download] %s\n", url)
	if err := download(url, binPath); err != nil {
		return fmt.Errorf("[error] Gagal download: %w", err)
	}
	if err := os.Chmod(binPath, 0755); err != nil {
		return fmt.Errorf("[error] Gagal set permission: %w", err)
	}

	dest, err := currentBinaryPath()
	if err != nil {
		return err
	}
	fmt.Printf("[install] Mengganti %s...\n", dest)
	if err := replaceBinary(binPath, dest); err != nil {
		return fmt.Errorf("[error] Gagal instalasi: %w", err)
	}
	fmt.Printf("[done] archstack %s terpasang. Jalankan ulang untuk pakai versi baru.\n", latest)
	return nil
}

func fetchLatestTag() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", repoName)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var rel struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	if rel.TagName == "" {
		return "", fmt.Errorf("tag_name kosong")
	}
	return rel.TagName, nil
}

func goarch() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "amd64", nil
	case "arm64":
		return "arm64", nil
	}
	return "", fmt.Errorf("[error] Arch tidak didukung: %s", runtime.GOARCH)
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func currentBinaryPath() (string, error) {
	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			return resolved, nil
		}
		return exe, nil
	}
	if _, err := exec.LookPath(repoName); err == nil {
		return "/usr/local/bin/" + repoName, nil
	}
	return "/usr/local/bin/" + repoName, nil
}

func replaceBinary(src, dest string) error {
	if isWritable(dest) {
		return os.Rename(src, dest)
	}
	cmd := exec.Command("sudo", "install", "-m", "0755", src, dest)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isWritable(path string) bool {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".archstack-wtest-")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}
