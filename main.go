// pwsh-updater installs or updates PowerShell on Linux from the official GitHub releases.
package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	installDir = "/opt/microsoft/powershell/7"
	linkPath   = "/usr/bin/pwsh"
	repo       = "PowerShell/PowerShell"
)

var version = "dev" // overridden with -ldflags "-X main.version=..."

const (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[1;33m"
	cyan   = "\033[0;36m"
	reset  = "\033[0m"
)

func info(format string, a ...any) { fmt.Printf(cyan+format+reset+"\n", a...) }
func ok(format string, a ...any)   { fmt.Printf(green+format+reset+"\n", a...) }
func warn(format string, a ...any) { fmt.Printf(yellow+format+reset+"\n", a...) }

func main() {
	var (
		showVersion = flag.Bool("version", false, "print pwsh-updater version")
		check       = flag.Bool("check", false, "only check whether a newer version is available")
		force       = flag.Bool("force", false, "reinstall even if the version is already installed")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pwsh-updater [flags] [version]\n\n")
		fmt.Fprintf(os.Stderr, "Installs or updates PowerShell (latest stable if no version is given).\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("pwsh-updater", version)
		return
	}
	if err := run(strings.TrimPrefix(flag.Arg(0), "v"), *check, *force); err != nil {
		fmt.Fprintf(os.Stderr, red+"Error: %v"+reset+"\n", err)
		os.Exit(1)
	}
}

func run(target string, check, force bool) error {
	arch, err := releaseArch()
	if err != nil {
		return err
	}
	if target == "" {
		info("Looking up the latest PowerShell release...")
		if target, err = latestVersion(); err != nil {
			return err
		}
	}
	current := installedVersion()
	if current == "" {
		warn("PowerShell is not installed.")
	} else {
		info("Installed: %s — Target: %s", current, target)
	}
	if check {
		if current == target {
			ok("Up to date.")
		} else {
			warn("Update available: %s", target)
		}
		return nil
	}
	if current == target && !force {
		ok("PowerShell %s is already installed (use --force to reinstall).", target)
		return nil
	}

	if os.Geteuid() != 0 {
		return reexecWithSudo()
	}
	return install(target, arch)
}

// releaseArch maps GOARCH to the suffix used in PowerShell release archives.
func releaseArch() (string, error) {
	if runtime.GOOS != "linux" {
		return "", errors.New("only Linux is supported")
	}
	switch runtime.GOARCH {
	case "amd64":
		return "x64", nil
	case "arm64":
		return "arm64", nil
	case "arm":
		return "arm32", nil
	}
	return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
}

func latestVersion() (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/releases/latest", nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API: %s", resp.Status)
	}
	var rel struct {
		Tag string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.Tag, "v"), nil
}

// installedVersion returns the version reported by `pwsh`, or "" if unavailable.
func installedVersion() string {
	out, err := exec.Command(linkPath, "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func reexecWithSudo() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	warn("Root privileges required, re-running with sudo...")
	cmd := exec.Command("sudo", append([]string{self}, os.Args[1:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func install(ver, arch string) error {
	archive := fmt.Sprintf("powershell-%s-linux-%s.tar.gz", ver, arch)
	base := fmt.Sprintf("https://github.com/%s/releases/download/v%s/", repo, ver)

	tmp, err := os.MkdirTemp("", "pwsh-updater-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	info("Downloading PowerShell v%s...", ver)
	archivePath := filepath.Join(tmp, archive)
	sum, err := download(base+archive, archivePath)
	if err != nil {
		return err
	}
	if err := verifyChecksum(base+"hashes.sha256", archive, sum); err != nil {
		return err
	}

	info("Extracting archive...")
	staging := filepath.Join(tmp, "extract")
	if err := extract(archivePath, staging); err != nil {
		return err
	}

	info("Installing to %s...", installDir)
	if err := os.RemoveAll(installDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(installDir), 0o755); err != nil {
		return err
	}
	if err := moveDir(staging, installDir); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Join(installDir, "pwsh"), 0o755); err != nil {
		return err
	}
	_ = os.Remove(linkPath)
	if err := os.Symlink(filepath.Join(installDir, "pwsh"), linkPath); err != nil {
		return err
	}

	ok("PowerShell v%s has been successfully installed.", ver)
	out, _ := exec.Command(linkPath, "-v").CombinedOutput()
	fmt.Print(string(out))
	return nil
}

// download saves url to dest and returns the SHA-256 of the content.
func download(url, dest string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: %s", url, resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// verifyChecksum compares sum with the entry for name in the release's hashes.sha256 file.
// A missing hashes file only produces a warning; a mismatch is an error.
func verifyChecksum(url, name, sum string) error {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		warn("Checksum file unavailable, skipping verification.")
		return nil
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == name {
			if !strings.EqualFold(fields[0], sum) {
				return errors.New("checksum mismatch, aborting")
			}
			ok("Checksum verified.")
			return nil
		}
	}
	warn("No checksum entry for %s, skipping verification.", name)
	return nil
}

func extract(archive, dest string) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	root := filepath.Clean(dest) + string(os.PathSeparator)
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(target+string(os.PathSeparator), root) {
			return fmt.Errorf("illegal path in archive: %s", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, 0o755)
		case tar.TypeReg:
			err = writeFile(target, tr, os.FileMode(hdr.Mode).Perm())
		case tar.TypeSymlink:
			if err = os.MkdirAll(filepath.Dir(target), 0o755); err == nil {
				err = os.Symlink(hdr.Linkname, target)
			}
		}
		if err != nil {
			return err
		}
	}
}

func writeFile(path string, r io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// moveDir renames src to dst, falling back to `cp -a` across filesystems.
func moveDir(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if out, err := exec.Command("cp", "-a", src, dst).CombinedOutput(); err != nil {
		return fmt.Errorf("copy failed: %v: %s", err, out)
	}
	return nil
}
