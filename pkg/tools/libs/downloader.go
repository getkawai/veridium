package libs

import (
	"context"
	"errors"
	"sync"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
)

var (
	ErrUnsupportedLibrary = errors.New("unsupported library type")
)

type Downloader interface {
	LatestVersion() (string, error)
	Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error
	LibraryName(os string) string
}

var (
	downloadersMu sync.RWMutex
	downloaders   = make(map[LibraryType]Downloader)
)

func registerDownloader(libType LibraryType, d Downloader) {
	downloadersMu.Lock()
	defer downloadersMu.Unlock()
	downloaders[libType] = d
}

// GetDownloader returns the registered downloader for the given library type.
// Returns nil if no downloader is registered for the type.
func GetDownloader(libType LibraryType) Downloader {
	downloadersMu.RLock()
	defer downloadersMu.RUnlock()
	return downloaders[libType]
}

func (lib *Libs) getLatestVersion() (string, error) {
	downloader := GetDownloader(lib.libType)
	if downloader == nil {
		return "", ErrUnsupportedLibrary
	}
	return downloader.LatestVersion()
}

func (lib *Libs) downloadLibrary(ctx context.Context, version, dest string, progress download.ProgressCallback) error {
	downloader := GetDownloader(lib.libType)
	if downloader == nil {
		return ErrUnsupportedLibrary
	}
	return downloader.Download(ctx, lib.arch.String(), lib.os.String(), lib.processor.String(), version, dest, progress)
}
