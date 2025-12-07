package main

import (
	"log"

	"github.com/kawai-network/veridium/internal/stablediffusion"
)

func main() {
	manager := stablediffusion.NewStableDiffusionReleaseManager()

	log.Println("Downloading Stable Diffusion binary...")
	err := manager.DownloadRelease("latest", func(progress float64) {
		if int(progress)%10 == 0 {
			log.Printf("Progress: %.1f%%", progress)
		}
	})

	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	log.Println("Download completed successfully!")
}
