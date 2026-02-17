package libs

import (
	"github.com/kawai-network/veridium/pkg/tools/libs/github"
	"github.com/kawai-network/whisper"
)

func init() {
	registerDownloader(LibraryWhisper, &github.ReleaseDownloader{
		Owner: "kawai-network",
		Repo:  "whisper",
		AssetName: func(os, arch, version string) string {
			return whisper.LibraryName(os)
		},
	})
}
