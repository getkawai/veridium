package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/network"
)

const (
	versionFile = "version.json"
	localFolder = "libraries"
)

type Logger func(ctx context.Context, msg string, args ...any)

type ProgressCallback func(bytesComplete, totalBytes int64, mbps float64, done bool)

type VersionTag struct {
	Library   LibraryType `json:"library"`
	Version   string      `json:"version"`
	Arch      string      `json:"arch"`
	OS        string      `json:"os"`
	Processor string      `json:"processor"`
	Latest    string      `json:"-"`
}

type Options struct {
	BasePath     string
	Arch         download.Arch
	OS           download.OS
	Processor    download.Processor
	AllowUpgrade bool
	Version      string
	LibraryType  LibraryType
}

type Option func(*Options)

func WithBasePath(basePath string) Option {
	return func(o *Options) {
		o.BasePath = basePath
	}
}

func WithArch(arch download.Arch) Option {
	return func(o *Options) {
		o.Arch = arch
	}
}

func WithOS(opSys download.OS) Option {
	return func(o *Options) {
		o.OS = opSys
	}
}

func WithProcessor(processor download.Processor) Option {
	return func(o *Options) {
		o.Processor = processor
	}
}

func WithAllowUpgrade(allow bool) Option {
	return func(o *Options) {
		o.AllowUpgrade = allow
	}
}

func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func WithLibraryType(libType LibraryType) Option {
	return func(o *Options) {
		o.LibraryType = libType
	}
}

type Libs struct {
	path         string
	arch         download.Arch
	os           download.OS
	processor    download.Processor
	allowUpgrade bool
	version      string
	libType      LibraryType
}

func New(opts ...Option) (*Libs, error) {
	arch, err := defaults.Arch("")
	if err != nil {
		return nil, err
	}

	opSys, err := defaults.OS("")
	if err != nil {
		return nil, err
	}

	processor, err := defaults.Processor("")
	if err != nil {
		return nil, err
	}

	options := Options{
		BasePath:     "",
		Arch:         arch,
		OS:           opSys,
		Processor:    processor,
		AllowUpgrade: true,
		LibraryType:  LibraryLlama,
	}

	for _, opt := range opts {
		opt(&options)
	}

	basePath := defaults.BaseDir(options.BasePath)

	lib := Libs{
		path:         filepath.Join(basePath, localFolder, options.LibraryType.Subfolder()),
		arch:         options.Arch,
		os:           options.OS,
		processor:    options.Processor,
		allowUpgrade: options.AllowUpgrade,
		version:      options.Version,
		libType:      options.LibraryType,
	}

	return &lib, nil
}

func (lib *Libs) LibsPath() string {
	return lib.path
}

func (lib *Libs) Arch() download.Arch {
	return lib.arch
}

func (lib *Libs) OS() download.OS {
	return lib.os
}

func (lib *Libs) Processor() download.Processor {
	return lib.processor
}

func (lib *Libs) LibraryType() LibraryType {
	return lib.libType
}

func (lib *Libs) Download(ctx context.Context, log Logger) (VersionTag, error) {
	return lib.DownloadWithProgress(ctx, log, nil)
}

func (lib *Libs) DownloadWithProgress(ctx context.Context, log Logger, progressCb ProgressCallback) (VersionTag, error) {
	if !hasNetwork() {
		vt, err := lib.InstalledVersion()
		if err != nil {
			return VersionTag{}, fmt.Errorf("download: no network available: %w", err)
		}

		log(ctx, "download-libraries: no network available, using current version")
		return vt, nil
	}

	log(ctx, "download-libraries: check libraries version information", "library", lib.libType, "arch", lib.arch, "os", lib.os, "processor", lib.processor)

	tag, err := lib.VersionInformation()
	if err != nil {
		if tag.Latest == "" {
			return VersionTag{}, fmt.Errorf("download-libraries: error retrieving version info: %w", err)
		}

		log(ctx, "download-libraries: unable to check latest version, using installed version", "library", lib.libType, "arch", lib.arch, "os", lib.os, "processor", lib.processor, "latest", tag.Latest, "current", tag.Version)
		return tag, nil
	}

	if lib.version != "" {
		tag.Latest = lib.version
	}

	log(ctx, "download-libraries: check installation", "library", lib.libType.DisplayName(), "arch", lib.arch, "os", lib.os, "processor", lib.processor, "latest", tag.Latest, "current", tag.Version)

	if isTagMatch(tag, lib) {
		log(ctx, "download-libraries: already installed", "library", lib.libType.DisplayName(), "latest", tag.Latest, "current", tag.Version)
		return tag, nil
	}

	if !lib.allowUpgrade {
		log(ctx, "download-libraries: bypassing upgrade", "library", lib.libType.DisplayName(), "latest", tag.Latest, "current", tag.Version)
		return tag, nil
	}

	log(ctx, "download-libraries waiting to start download...", "library", lib.libType.DisplayName(), "tag", tag.Latest)

	newTag, err := lib.downloadVersionWithProgress(ctx, log, tag.Latest, progressCb)
	if err != nil {
		log(ctx, "download-libraries: installation error", "library", lib.libType.DisplayName(), "ERROR", err)

		vt, checkErr := lib.InstalledVersion()
		if checkErr != nil {
			return VersionTag{}, fmt.Errorf("download: failed to install %s: %w", lib.libType.DisplayName(), err)
		}

		log(ctx, "download-libraries: failed to install new version, using current version")
		return vt, nil
	}

	log(ctx, "download-libraries: updated library installed", "library", lib.libType.DisplayName(), "old-version", tag.Version, "current", newTag.Version)

	return newTag, nil
}

func (lib *Libs) InstalledVersion() (VersionTag, error) {
	versionInfoPath := filepath.Join(lib.path, versionFile)

	d, err := os.ReadFile(versionInfoPath)
	if err != nil {
		return VersionTag{}, fmt.Errorf("installed-version: unable to read version info file: %w", err)
	}

	var tag VersionTag
	if err := json.Unmarshal(d, &tag); err != nil {
		return VersionTag{}, fmt.Errorf("installed-version: unable to parse version info file: %w", err)
	}

	downloader := GetDownloader(lib.libType)
	if downloader == nil {
		return VersionTag{}, ErrUnsupportedLibrary
	}

	libraryName := downloader.LibraryName(lib.os.String())
	if libraryName == "" {
		return VersionTag{}, fmt.Errorf("installed-version: empty library name for %s", lib.libType.DisplayName())
	}

	libraryPath := filepath.Join(lib.path, libraryName)
	if _, err := os.Stat(libraryPath); err != nil {
		return VersionTag{}, fmt.Errorf("installed-version: library file missing: %s: %w", libraryPath, err)
	}

	tag.Library = lib.libType
	return tag, nil
}

func (lib *Libs) VersionInformation() (VersionTag, error) {
	tag, err := lib.InstalledVersion()
	tag.Library = lib.libType

	if err != nil {
		tag = VersionTag{
			Library:   lib.libType,
			Version:   "",
			Arch:      lib.arch.String(),
			OS:        lib.os.String(),
			Processor: lib.processor.String(),
		}
	}

	version, err := lib.getLatestVersion()
	if err != nil {
		return tag, fmt.Errorf("version-information: unable to get latest version of %s: %w", lib.libType.DisplayName(), err)
	}

	tag.Latest = version

	return tag, nil
}

func (lib *Libs) DownloadVersion(ctx context.Context, log Logger, version string) (VersionTag, error) {
	return lib.downloadVersionWithProgress(ctx, log, version, nil)
}

func (lib *Libs) downloadVersionWithProgress(ctx context.Context, log Logger, version string, progressCb ProgressCallback) (VersionTag, error) {
	tempPath := filepath.Join(lib.path, "temp")

	progress := func(src string, currentSize int64, totalSize int64, mibPerSec float64, complete bool) {
		if progressCb != nil {
			progressCb(currentSize, totalSize, mibPerSec, complete)
		}
		log(ctx, fmt.Sprintf("\r\x1b[Kdownload-libraries: Downloading %s... %d MiB of %d MiB (%.2f MiB/s)", src, currentSize/(1024*1024), totalSize/(1024*1024), mibPerSec))
	}

	err := lib.downloadLibrary(ctx, version, tempPath, progress)
	if err != nil {
		os.RemoveAll(tempPath)
		return VersionTag{}, fmt.Errorf("download-version: unable to install %s: %w", lib.libType.DisplayName(), err)
	}

	if err := lib.swapTempForLib(tempPath); err != nil {
		os.RemoveAll(tempPath)
		return VersionTag{}, fmt.Errorf("download-version: unable to swap temp for lib: %w", err)
	}

	if err := lib.createVersionFile(version); err != nil {
		return VersionTag{}, fmt.Errorf("download-version: unable to create version file: %w", err)
	}

	return lib.VersionInformation()
}

func (lib *Libs) swapTempForLib(tempPath string) error {
	if err := os.MkdirAll(lib.path, 0755); err != nil {
		return fmt.Errorf("swap-temp-for-lib: unable to create lib path: %w", err)
	}

	entries, err := os.ReadDir(lib.path)
	if err != nil {
		return fmt.Errorf("swap-temp-for-lib: unable to read libPath: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == "temp" || entry.Name() == versionFile {
			continue
		}

		os.Remove(filepath.Join(lib.path, entry.Name()))
	}

	tempEntries, err := os.ReadDir(tempPath)
	if err != nil {
		return fmt.Errorf("swap-temp-for-lib: unable to read temp: %w", err)
	}

	for _, entry := range tempEntries {
		src := filepath.Join(tempPath, entry.Name())
		dst := filepath.Join(lib.path, entry.Name())
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("swap-temp-for-lib: unable to move %s: %w", entry.Name(), err)
		}
	}

	os.RemoveAll(tempPath)

	return nil
}

func (lib *Libs) createVersionFile(version string) error {
	versionInfoPath := filepath.Join(lib.path, versionFile)

	if err := os.MkdirAll(lib.path, 0755); err != nil {
		return fmt.Errorf("create-version-file: creating directory: %w", err)
	}

	f, err := os.Create(versionInfoPath)
	if err != nil {
		return fmt.Errorf("create-version-file: creating version info file: %w", err)
	}
	defer f.Close()

	t := VersionTag{
		Library:   lib.libType,
		Version:   version,
		Arch:      lib.arch.String(),
		OS:        lib.os.String(),
		Processor: lib.processor.String(),
	}

	d, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("create-version-file: marshalling version info: %w", err)
	}

	if _, err := f.Write(d); err != nil {
		return fmt.Errorf("create-version-file: writing version info: %w", err)
	}

	return nil
}

func isTagMatch(tag VersionTag, libs *Libs) bool {
	return tag.Latest == tag.Version && tag.Arch == libs.arch.String() && tag.OS == libs.os.String() && tag.Processor == libs.processor.String()
}

// hasNetwork checks network connectivity using the shared network utility.
func hasNetwork() bool {
	return network.HasNetwork()
}
