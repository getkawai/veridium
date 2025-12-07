package onnx

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kawai-network/veridium/pkg/grab"
)

const defaultONNXVersion = "1.22.0"

type platformSpec struct {
	archivePattern string
	checksum       string
}

var platformMatrix = map[string]map[string]platformSpec{
	"linux": {
		"amd64": {
			archivePattern: "onnxruntime-linux-x64-%s.tgz",
			checksum:       "8344d55f93d5bc5021ce342db50f62079daf39aaafb5d311a451846228be49b3",
		},
		"arm64": {
			archivePattern: "onnxruntime-linux-aarch64-%s.tgz",
			checksum:       "bb76395092d150b52c7092dc6b8f2fe4d80f0f3bf0416d2f269193e347e24702",
		},
	},
	"darwin": {
		"amd64": {
			archivePattern: "onnxruntime-osx-universal2-%s.tgz",
			checksum:       "cfa6f6584d87555ed9f6e7e8a000d3947554d589efe3723b8bfa358cd263d03c",
		},
		"arm64": {
			archivePattern: "onnxruntime-osx-universal2-%s.tgz",
			checksum:       "cfa6f6584d87555ed9f6e7e8a000d3947554d589efe3723b8bfa358cd263d03c",
		},
	},
}

var (
	once          sync.Once
	initErr       error
	isInitialized bool
	installMutex  sync.Mutex
)

// Config holds configuration for ONNX Runtime installation
type Config struct {
	Version string // ONNX Runtime version (default: 1.22.0)
	DestDir string // Destination directory (default: /usr)
	TmpDir  string // Temporary directory (default: system temp)
	Arch    string // Target architecture (default: runtime.GOARCH)
	OS      string // Target operating system (default: runtime.GOOS)
	Silent  bool   // Suppress log output
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: firstNonEmpty(os.Getenv("ONNX_VERSION"), defaultONNXVersion),
		DestDir: "/usr",
		TmpDir:  firstNonEmpty(os.Getenv("TMPDIR"), os.TempDir()),
		Arch:    normalizeArch(firstNonEmpty(os.Getenv("PHOTOPRISM_ARCH"), runtime.GOARCH)),
		OS:      normalizeOS(runtime.GOOS),
		Silent:  false,
	}
}

// EnsureInstalled checks if ONNX Runtime is installed, and installs it if not.
// This is the recommended function to call on app startup.
// It's thread-safe and will only install once even if called multiple times.
func EnsureInstalled(cfg *Config) error {
	once.Do(func() {
		if cfg == nil {
			cfg = DefaultConfig()
		}

		// Check if already installed
		if IsInstalled(cfg.DestDir) {
			if !cfg.Silent {
				log.Printf("✅ ONNX Runtime already installed in %s", cfg.DestDir)
			}
			isInitialized = true
			return
		}

		// Install if not present
		initErr = Install(cfg)
		if initErr == nil {
			isInitialized = true
		}
	})
	return initErr
}

// AutoInstall automatically installs ONNX Runtime on first call (thread-safe)
// Subsequent calls return the cached result
// Deprecated: Use EnsureInstalled instead which checks for existing installation
func AutoInstall(cfg *Config) error {
	return EnsureInstalled(cfg)
}

// Install installs ONNX Runtime with the given configuration
func Install(cfg *Config) error {
	installMutex.Lock()
	defer installMutex.Unlock()

	if cfg == nil {
		cfg = DefaultConfig()
	}

	return runInstaller(cfg)
}

// IsInstalled checks if ONNX Runtime is already installed
func IsInstalled(destDir string) bool {
	if destDir == "" {
		destDir = "/usr"
	}
	libDir := filepath.Join(destDir, "lib")

	// Check for common ONNX Runtime library files
	patterns := []string{
		"libonnxruntime.so*",
		"libonnxruntime.*.dylib",
		"libonnxruntime.dylib",
	}

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(libDir, pattern))
		if len(matches) > 0 {
			return true
		}
	}

	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func firstArg(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func normalizeOS(osName string) string {
	osName = strings.ToLower(osName)
	switch osName {
	case "linux", "darwin":
		return osName
	case "mac", "macos", "osx":
		return "darwin"
	default:
		return osName
	}
}

func normalizeArch(arch string) string {
	switch strings.ToLower(arch) {
	case "x86_64", "x86-64", "amd64":
		return "amd64"
	case "arm64", "aarch64":
		return "arm64"
	}
	return strings.ToLower(arch)
}

func runInstaller(cfg *Config) error {
	if cfg.OS == "" || cfg.Arch == "" {
		return errors.New("missing system or architecture information")
	}

	spec, err := resolvePlatformSpec(cfg.OS, cfg.Arch)
	if err != nil {
		return err
	}

	destAbs, err := filepath.Abs(cfg.DestDir)
	if err != nil {
		return fmt.Errorf("resolve destination: %w", err)
	}

	if err := ensureWritable(destAbs); err != nil {
		return err
	}

	if err := os.MkdirAll(destAbs, 0o755); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	if err := os.MkdirAll(cfg.TmpDir, 0o755); err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	archiveName := fmt.Sprintf(spec.archivePattern, cfg.Version)
	archivePath := filepath.Join(cfg.TmpDir, archiveName)

	logf(cfg.Silent, "📦 Installing ONNX Runtime %s (%s/%s) into %s", cfg.Version, cfg.OS, cfg.Arch, destAbs)

	if err := maybeDownloadArchive(archivePath, spec.checksum, cfg.Version, archiveName, cfg.Silent); err != nil {
		return err
	}

	if err := verifyChecksum(archivePath, spec.checksum); err != nil {
		return err
	}

	logf(cfg.Silent, "📂 Extracting to %s...", destAbs)
	if err := extractArchive(archivePath, destAbs); err != nil {
		return err
	}

	if err := normalizeLibraryLayout(destAbs, cfg.Version, cfg.Silent); err != nil {
		return err
	}

	if cfg.OS == "linux" {
		if err := runLdconfig(destAbs); err != nil {
			logf(cfg.Silent, "⚠️  ldconfig failed: %v", err)
		}
	}

	logf(cfg.Silent, "✅ ONNX Runtime %s installed in %s", cfg.Version, destAbs)
	return nil
}

func logf(silent bool, format string, args ...interface{}) {
	if !silent {
		log.Printf(format, args...)
	}
}

func resolvePlatformSpec(system, arch string) (platformSpec, error) {
	platforms, ok := platformMatrix[system]
	if !ok {
		return platformSpec{}, fmt.Errorf("unsupported operating system %q", system)
	}
	spec, ok := platforms[arch]
	if !ok {
		return platformSpec{}, fmt.Errorf("unsupported architecture %q for %s", arch, system)
	}
	return spec, nil
}

func ensureWritable(dest string) error {
	info, err := os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat destination: %w", err)
	}
	if os.IsNotExist(err) {
		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("destination %s is not a directory", dest)
	}

	if (dest == "/usr" || dest == "/usr/local") && os.Geteuid() != 0 {
		return fmt.Errorf("run as root to install into %s", dest)
	}

	return nil
}

func maybeDownloadArchive(path, checksum, version, archiveName string, silent bool) error {
	if fileExists(path) {
		if err := verifyChecksum(path, checksum); err == nil {
			logf(silent, "ℹ️  Using cached archive %s", path)
			return nil
		}
		logf(silent, "⚠️  Cached archive %s failed checksum, re-downloading", path)
		_ = os.Remove(path)
	}

	today := timeStamp()

	primaryURL := fmt.Sprintf("https://dl.photoprism.app/onnx/runtime/v%s/%s", version, archiveName)
	if today != "" {
		primaryURL = fmt.Sprintf("%s?%s", primaryURL, today)
	}
	fallbackURL := fmt.Sprintf("https://github.com/microsoft/onnxruntime/releases/download/v%s/%s", version, archiveName)

	logf(silent, "⬇️  Downloading %s", primaryURL)
	if err := downloadWithGrab(path, primaryURL); err != nil {
		logf(silent, "⚠️  Primary download failed: %v", err)
		logf(silent, "⬇️  Trying fallback %s", fallbackURL)
		if err := downloadWithGrab(path, fallbackURL); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
	}
	return nil
}

func downloadWithGrab(dst, url string) error {
	req, err := grab.NewRequest(dst, url)
	if err != nil {
		return err
	}

	resp := grab.DefaultClient.Do(req)
	<-resp.Done

	if err := resp.Err(); err != nil {
		return err
	}

	return nil
}

func verifyChecksum(path, expected string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("hash archive: %w", err)
	}

	actual := hex.EncodeToString(hash.Sum(nil))
	expected = strings.ToLower(expected)

	if actual != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}

	return nil
}

func extractArchive(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("decompress archive: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read archive: %w", err)
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, fsMode(header.FileInfo().Mode())); err != nil {
				return fmt.Errorf("create dir %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := writeFileFromTar(tr, target, header.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("create symlink dir: %w", err)
			}
			if err := os.RemoveAll(target); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove existing symlink %s: %w", target, err)
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("create symlink %s: %w", target, err)
			}
		default:
			// Ignore other types.
		}
	}

	return nil
}

func fsMode(mode os.FileMode) os.FileMode {
	if mode == 0 {
		return 0o755
	}
	return mode
}

func writeFileFromTar(r io.Reader, target string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create file dir: %w", err)
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fsMode(mode))
	if err != nil {
		return fmt.Errorf("create file %s: %w", target, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("write file %s: %w", target, err)
	}

	return nil
}

func normalizeLibraryLayout(dest, version string, silent bool) error {
	outputLibDir := filepath.Join(dest, "lib")
	if err := os.MkdirAll(outputLibDir, 0o755); err != nil {
		return fmt.Errorf("create lib dir: %w", err)
	}

	candidates := []string{
		filepath.Join(dest, fmt.Sprintf("onnxruntime-linux-x64-%s", version)),
		filepath.Join(dest, fmt.Sprintf("onnxruntime-linux-aarch64-%s", version)),
		filepath.Join(dest, fmt.Sprintf("onnxruntime-osx-universal2-%s", version)),
	}

	for _, base := range candidates {
		libDir := filepath.Join(base, "lib")
		if fi, err := os.Stat(libDir); err == nil && fi.IsDir() {
			if err := copyLibraries(libDir, outputLibDir); err != nil {
				return err
			}
			if err := os.RemoveAll(base); err != nil {
				logf(silent, "⚠️  Failed to remove %s: %v", base, err)
			}
		}
	}

	return nil
}

func copyLibraries(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "libonnxruntime") {
			continue
		}

		srcPath := filepath.Join(srcDir, name)
		destPath := filepath.Join(destDir, name)

		info, err := os.Lstat(srcPath)
		if err != nil {
			return fmt.Errorf("stat %s: %w", srcPath, err)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(srcPath)
			if err != nil {
				return fmt.Errorf("readlink %s: %w", srcPath, err)
			}
			if err := os.RemoveAll(destPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove %s: %w", destPath, err)
			}
			if err := os.Symlink(target, destPath); err != nil {
				return fmt.Errorf("create symlink %s: %w", destPath, err)
			}
			continue
		}

		if err := copyFile(srcPath, destPath, info.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fsMode(mode))
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s: %w", dst, err)
	}

	return nil
}

func runLdconfig(dest string) error {
	if dest == "/usr" || dest == "/usr/local" {
		return exec.Command("ldconfig").Run()
	}
	return exec.Command("ldconfig", "-n", filepath.Join(dest, "lib")).Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func timeStamp() string {
	if v := strings.TrimSpace(os.Getenv("TODAY")); v != "" {
		return strings.ReplaceAll(v, "\n", "")
	}
	return time.Now().UTC().Format("20060102")
}
