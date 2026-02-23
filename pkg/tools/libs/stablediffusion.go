package libs

import (
	"context"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	sddownload "github.com/kawai-network/veridium/pkg/stablediffusion/download"
)

func init() {
	registerDownloader(LibraryStableDiffusion, &sdDownloader{})
}

type sdDownloader struct{}

func (d *sdDownloader) LatestVersion() (string, error) {
	// Return the default tested version
	return sddownload.DefaultVersion, nil
}

func (d *sdDownloader) Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error {
	if version == "" {
		version = sddownload.DefaultVersion
	}

	// Note: sddownload.GetWithContext handles architecture/OS detection internally
	// The arch, os, and processor parameters are not currently used by the SD download
	// but are kept for interface compatibility with the Downloader interface

	return sddownload.GetWithContext(ctx, version, dest, func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
		if progress != nil {
			progress(url, bytesComplete, totalBytes, mbps, done)
		}
	})
}

func (d *sdDownloader) LibraryName(os string) string {
	return sddownload.LibraryName()
}
