//go:build linux

package stablediffusion

import "log"

// selectBestAsset selects the best asset for Linux
func (sdrm *StableDiffusionReleaseManager) selectBestAsset(assets []Asset) *Asset {
	var patterns []string

	// Linux releases - actual format: sd-master--bin-Linux-Ubuntu-24.04-x86_64.zip
	if isARM64() {
		patterns = append(patterns, ".*linux.*arm64")
		patterns = append(patterns, ".*ubuntu.*arm64")
	} else {
		patterns = append(patterns, ".*linux.*x86_64")
		patterns = append(patterns, ".*ubuntu.*x86_64")
		patterns = append(patterns, ".*linux.*x64")
	}

	log.Printf("Trying %d patterns for Linux asset selection:", len(patterns))
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

	log.Printf("ERROR: No patterns matched any available assets for Linux")
	return nil
}
