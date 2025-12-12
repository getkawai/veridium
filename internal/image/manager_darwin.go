//go:build darwin

package image

import "log"

// selectBestAsset selects the best asset for macOS
func (sdrm *StableDiffusion) selectBestAsset(assets []Asset) *Asset {
	var patterns []string

	// macOS releases - actual format: sd-master--bin-Darwin-macOS-15.5-arm64.zip
	// Check architecture
	if isARM64() {
		patterns = append(patterns, ".*darwin.*arm64")
		patterns = append(patterns, ".*macos.*arm64")
	} else {
		patterns = append(patterns, ".*darwin.*x64")
		patterns = append(patterns, ".*darwin.*x86_64")
		patterns = append(patterns, ".*macos.*x64")
	}

	log.Printf("Trying %d patterns for macOS asset selection:", len(patterns))
	for i, pattern := range patterns {
		log.Printf("  Pattern %d: %s", i+1, pattern)
	}

	// Try each pattern in priority order
	for _, pattern := range patterns {
		log.Printf("Trying pattern: %s", pattern)
		for _, asset := range assets {
			if sdrm.matchesPattern(asset.Name, pattern) {
				log.Printf("SUCCESS: Matched pattern '%s' with asset '%s'", pattern, asset.Name)
				return &asset
			}
		}
		log.Printf("No match found for pattern: %s", pattern)
	}

	log.Printf("ERROR: No patterns matched any available assets for macOS")
	return nil
}
