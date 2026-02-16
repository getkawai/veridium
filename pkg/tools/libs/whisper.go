package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
	whisper "github.com/kawai-network/whisper"
)

const whisperLibReleaseURL = "https://api.github.com/repos/kawai-network/whisper/releases/latest"

func init() {
	registerDownloader(LibraryWhisper, &whisperDownloader{})
}

type whisperDownloader struct{}

type whisperRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

func (d *whisperDownloader) LatestVersion() (string, error) {
	release, err := d.getLatestRelease()
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

func (d *whisperDownloader) Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error {
	libName := whisper.LibraryName(runtime.GOOS)
	libFile := fmt.Sprintf("%s/%s", dest, libName)

	release, err := d.getLatestRelease()
	if err != nil {
		return err
	}

	assetURL := ""
	for _, asset := range release.Assets {
		if asset.Name == libName {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("library %s not found in release %s", libName, release.TagName)
	}

	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progress != nil {
			progress(src, currentSize, totalSize, mibPerSec, complete)
		}
	}

	_, err = downloader.Download(ctx, assetURL, libFile, progressFunc, downloader.SizeIntervalMIB)
	return err
}

func (d *whisperDownloader) LibraryName(os string) string {
	return whisper.LibraryName(os)
}

func (d *whisperDownloader) getLatestRelease() (*whisperRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(whisperLibReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release whisperRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}
