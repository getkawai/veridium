package libs

import (
	"context"
	"runtime"

	"github.com/kawai-network/stablediffusion"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	sddownload "github.com/kawai-network/veridium/pkg/stablediffusion/download"
)

func init() {
	registerDownloader(LibraryStableDiffusion, &sdDownloader{})
}

type sdDownloader struct{}

func (d *sdDownloader) LatestVersion() (string, error) {
	return sddownload.SDLatestVersion()
}

func (d *sdDownloader) Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error {
	if version == "" {
		version = sddownload.DefaultVersion
	}

	parsedArch, err := sddownload.ParseArch(arch)
	if err != nil {
		parsedArch = sddownload.AMD64
		if runtime.GOARCH == "arm64" {
			parsedArch = sddownload.ARM64
		}
	}

	parsedOS, err := sddownload.ParseOS(os)
	if err != nil {
		parsedOS = sddownload.Linux
		switch runtime.GOOS {
		case "darwin":
			parsedOS = sddownload.Darwin
		case "windows":
			parsedOS = sddownload.Windows
		}
	}

	_ = parsedArch
	_ = parsedOS

	return sddownload.GetWithContext(ctx, version, dest, func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
		if progress != nil {
			progress(url, bytesComplete, totalBytes, mbps, done)
		}
	})
}

func (d *sdDownloader) LibraryName(os string) string {
	return stablediffusion.LibraryName()
}
