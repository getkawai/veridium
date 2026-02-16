package libs

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	"github.com/kawai-network/whisper"
)

func init() {
	registerDownloader(LibraryWhisper, &whisperDownloader{})
}

type whisperDownloader struct{}

func (d *whisperDownloader) LatestVersion() (string, error) {
	downloader := whisper.NewLibraryDownloader("")
	release, err := downloader.GetLatestRelease()
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

func (d *whisperDownloader) Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error {
	downloader := whisper.NewLibraryDownloader(dest)

	release, err := downloader.GetLatestRelease()
	if err != nil {
		return err
	}

	platform := whisper.DetectPlatform()
	asset, err := downloader.SelectBestLibrary(release, platform)
	if err != nil {
		return err
	}

	_, err = downloader.DownloadWithProgress(asset, func(bytesComplete, totalBytes int64, mbps float64, done bool) {
		if progress != nil && totalBytes > 0 {
			percent := float64(bytesComplete) / float64(totalBytes) * 100
			msg := fmt.Sprintf("Downloading: %.1f%%", percent)
			progress(asset.URL, bytesComplete, totalBytes, mbps, done)
			fmt.Printf("\r%s", msg)
		}
	})
	return err
}

func (d *whisperDownloader) LibraryName(os string) string {
	platform := whisper.DetectPlatform()
	return platform.Prefix + "gowhisper" + platform.Extension
}
