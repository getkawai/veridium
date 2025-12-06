package builtin

import (
	"database/sql"
	"log"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// RegisterAll registers all builtin tools
func RegisterAll(registry *tools.ToolRegistry) error {
	return RegisterAllWithDB(registry, nil)
}

// RegisterAllWithDB registers all builtin tools with optional database connection
// Some tools (like image describe) require database access
func RegisterAllWithDB(registry *tools.ToolRegistry, sqlDB *sql.DB) error {
	log.Println("Registering builtin tools (yzma)...")

	// Register lobe-web-browsing (search, crawlSinglePage, crawlMultiPages)
	if err := RegisterWebBrowsing(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: lobe-web-browsing (search, crawlSinglePage, crawlMultiPages)")

	// Register lobe-local-system (file operations)
	if err := RegisterLocalSystem(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: lobe-local-system (list, read, search, write, rename, move)")

	// Register lobe-image-designer (DALL-E compatible)
	if err := RegisterImageDesigner(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: lobe-image-designer (text2image)")

	// Register lobe-code-interpreter (Python execution)
	if err := RegisterCodeInterpreter(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: lobe-code-interpreter (python)")

	// Register calculator
	if err := RegisterCalculator(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: calculator")

	// Register lobe-image-describe (requires DB for querying VL descriptions)
	if sqlDB != nil {
		if err := RegisterImageDescribe(registry, sqlDB); err != nil {
			return err
		}
		log.Println("✅ Registered: lobe-image-describe (getImageDescription)")

		// Register lobe-video-describe (requires DB for querying Whisper transcriptions)
		if err := RegisterVideoDescribe(registry, sqlDB); err != nil {
			return err
		}
		log.Println("✅ Registered: lobe-video-describe (getVideoTranscription)")
	} else {
		log.Println("⚠️  Skipped: lobe-image-describe (no database connection)")
		log.Println("⚠️  Skipped: lobe-video-describe (no database connection)")
	}

	return nil
}

