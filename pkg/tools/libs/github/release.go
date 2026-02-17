package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

// ReleaseDownloader provides a generic downloader for GitHub releases.
// It handles fetching release metadata from GitHub API and downloading assets.
type ReleaseDownloader struct {
	Owner     string
	Repo      string
	AssetName func(os, arch, version string) string
}

// release represents a GitHub release response.
type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// LatestVersion fetches the latest release tag from GitHub.
func (d *ReleaseDownloader) LatestVersion() (string, error) {
	rel, err := d.fetchLatestRelease()
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}

// Download downloads the library asset for the specified OS/architecture.
// If version is empty, it fetches the latest release.
func (d *ReleaseDownloader) Download(
	ctx context.Context,
	arch, os, processor, version, dest string,
	progress download.ProgressCallback,
) error {
	libName := d.AssetName(os, arch, version)
	libFile := fmt.Sprintf("%s/%s", dest, libName)

	var rel *release
	var err error

	if version == "" {
		rel, err = d.fetchLatestRelease()
	} else {
		rel, err = d.fetchRelease(version)
	}

	if err != nil {
		return err
	}

	assetURL := ""
	for _, asset := range rel.Assets {
		if asset.Name == libName {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("library %s not found in release %s", libName, rel.TagName)
	}

	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progress != nil {
			progress(src, currentSize, totalSize, mibPerSec, complete)
		}
	}

	_, err = downloader.Download(ctx, assetURL, libFile, progressFunc, downloader.SizeIntervalMIB)
	return err
}

// LibraryName returns the library file name for the given OS.
// This is a placeholder that should be overridden by specific implementations.
func (d *ReleaseDownloader) LibraryName(os string) string {
	return d.AssetName(os, "", "")
}

func (d *ReleaseDownloader) fetchLatestRelease() (*release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", d.Owner, d.Repo)
	return d.doFetch(url)
}

func (d *ReleaseDownloader) fetchRelease(version string) (*release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", d.Owner, d.Repo, version)
	return d.doFetch(url)
}

func (d *ReleaseDownloader) doFetch(url string) (*release, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &rel, nil
}
