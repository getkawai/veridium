package libs

import (
	"context"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
)

func init() {
	registerDownloader(LibraryLlama, &llamaDownloader{})
}

type llamaDownloader struct{}

func (d *llamaDownloader) LatestVersion() (string, error) {
	return download.LlamaLatestVersion()
}

func (d *llamaDownloader) Download(ctx context.Context, arch, os, processor, version, dest string, progress download.ProgressCallback) error {
	return download.GetWithContext(ctx, arch, os, processor, version, dest, progress)
}

func (d *llamaDownloader) LibraryName(os string) string {
	return download.LibraryName(os)
}
