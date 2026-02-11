// Package downloader provides support for downloading files.
package downloader

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/kawai-network/grab"
	"github.com/kawai-network/x/constant"
)

// SizeInterval are pre-calculated size interval values.
const (
	SizeIntervalMIB    = 1024 * 1024
	SizeIntervalMIB10  = SizeIntervalMIB * 10
	SizeIntervalMIB100 = SizeIntervalMIB * 100
)

// ProgressFunc provides feedback on the progress of a file download.
type ProgressFunc func(src string, currentSize int64, totalSize int64, mibPerSec float64, complete bool)

// Download pulls down a single file from a url to a specified destination.
func Download(ctx context.Context, src string, dest string, progress ProgressFunc, sizeInterval int64) (bool, error) {
	if !hasNetwork() {
		return false, errors.New("download: no network available")
	}

	// Create grab request
	req, err := grab.NewRequest(dest, src)
	if err != nil {
		return false, fmt.Errorf("download: failed to create request: %w", err)
	}

	// Set context for cancellation
	req = req.WithContext(ctx)

	// Add authorization header if HF token is set
	if token := constant.GetRandomHfApiKey(); token != "" {
		req.HTTPRequest.Header.Set("Authorization", "Bearer "+token)
	}

	// Create grab client
	client := grab.NewClient()

	// Start download
	resp := client.Do(req)

	// Monitor progress if callback provided
	if progress != nil {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var lastReported int64

		for {
			select {
			case <-ticker.C:
				currentSize := resp.BytesComplete()
				totalSize := resp.Size()

				// Report progress at specified intervals
				if currentSize-lastReported >= sizeInterval {
					lastReported = currentSize
					mibPerSec := resp.BytesPerSecond() / SizeIntervalMIB
					progress(src, currentSize, totalSize, mibPerSec, false)
				}

			case <-resp.Done:
				// Download completed
				if err := resp.Err(); err != nil {
					return false, fmt.Errorf("download: failed to download: %w", err)
				}

				// Final progress report
				currentSize := resp.BytesComplete()
				totalSize := resp.Size()
				mibPerSec := resp.BytesPerSecond() / SizeIntervalMIB
				progress(src, currentSize, totalSize, mibPerSec, true)

				if currentSize == 0 {
					return false, nil
				}

				return true, nil
			}
		}
	}

	// Wait for download to complete if no progress callback
	<-resp.Done
	if err := resp.Err(); err != nil {
		return false, fmt.Errorf("download: failed to download: %w", err)
	}

	if resp.BytesComplete() == 0 {
		return false, nil
	}

	return true, nil
}

// =============================================================================

func hasNetwork() bool {
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 3*time.Second)
	if err != nil {
		return false
	}

	conn.Close()

	return true
}
