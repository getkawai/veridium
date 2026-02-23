package stablediffusion

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"unsafe"
)

type SDCacheMode int32

const (
	SDCacheDisabled SDCacheMode = iota
	SDCacheEasycache
	SDCacheUcache
	SDCacheDbcache
	SDCacheTaylorseer
	SDCacheCacheDit
)

type SDTilingParams struct {
	Enabled       bool
	TileSizeX     int32
	TileSizeY     int32
	TargetOverlap float32
	RelSizeX      float32
	RelSizeY      float32
}

type SDImage struct {
	Width   uint32
	Height  uint32
	Channel uint32
	Data    *uint8
}

type SDCacheParams struct {
	Mode                     SDCacheMode
	ReuseThreshold           float32
	StartPercent             float32
	EndPercent               float32
	ErrorDecayRate           float32
	UseRelativeThreshold     bool
	ResetErrorOnCompute      bool
	FnComputeBlocks          int32
	BnComputeBlocks          int32
	ResidualDiffThreshold    float32
	MaxWarmupSteps           int32
	MaxCachedSteps           int32
	MaxContinuousCachedSteps int32
	TaylorseerNDerivatives   int32
	TaylorseerSkipInterval   int32
	ScmMask                  *uint8
	ScmPolicyDynamic         bool
}

// SaveImage saves SDImage as PNG file.
func SaveImage(img *SDImage, path string) error {
	if img == nil || img.Data == nil {
		return fmt.Errorf("invalid image data")
	}

	bounds := image.Rect(0, 0, int(img.Width), int(img.Height))
	rgba := image.NewRGBA(bounds)
	data := unsafe.Slice(img.Data, img.Width*img.Height*img.Channel)

	for i := 0; i < int(img.Width*img.Height); i++ {
		index := i * int(img.Channel)
		x := i % int(img.Width)
		y := i / int(img.Width)

		r := data[index]
		g := data[index+1]
		b := data[index+2]
		rgba.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	return png.Encode(file, rgba)
}

// LoadImage loads image from file and converts to SDImage format.
func LoadImage(path string) (SDImage, error) {
	file, err := os.Open(path)
	if err != nil {
		return SDImage{}, fmt.Errorf("failed to open image file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return SDImage{}, fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	channel := 3
	dataSize := width * height * channel
	data := make([]uint8, dataSize)

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[index] = uint8(r >> 8)
			index++
			data[index] = uint8(g >> 8)
			index++
			data[index] = uint8(b >> 8)
			index++
		}
	}

	return SDImage{
		Width:   uint32(width),
		Height:  uint32(height),
		Channel: uint32(channel),
		Data:    &data[0],
	}, nil
}

// EncodeVideo encodes PNG frame sequence to video using FFmpeg.
func EncodeVideo(inputDir, outputPath string, framerate int) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: %v", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-framerate", strconv.Itoa(framerate),
		"-i", filepath.Join(inputDir, "frame_%04d.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		outputPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}
	return nil
}

// SaveFrames saves all video frames as PNG files.
func SaveFrames(frames []SDImage, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for i, frame := range frames {
		framePath := filepath.Join(outputDir, fmt.Sprintf("frame_%04d.png", i+1))
		if err := SaveImage(&frame, framePath); err != nil {
			return fmt.Errorf("failed to save frame %d: %v", i+1, err)
		}
	}
	return nil
}
