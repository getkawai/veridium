package libs

import (
	"context"
	"errors"

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

var downloaders = make(map[LibraryType]Downloader)

func registerDownloader(libType LibraryType, d Downloader) {
	downloaders[libType] = d
}

func getDownloader(libType LibraryType) Downloader {
	return downloaders[libType]
}

func GetDownloader(libType LibraryType) Downloader {
	return downloaders[libType]
}

func (lib *Libs) getLatestVersion() (string, error) {
	downloader := getDownloader(lib.libType)
	if downloader == nil {
		return "", ErrUnsupportedLibrary
	}
	return downloader.LatestVersion()
}

func (lib *Libs) downloadLibrary(ctx context.Context, version, dest string, progress download.ProgressCallback) error {
	downloader := getDownloader(lib.libType)
	if downloader == nil {
		return ErrUnsupportedLibrary
	}
	return downloader.Download(ctx, lib.arch.String(), lib.os.String(), lib.processor.String(), version, dest, progress)
}
