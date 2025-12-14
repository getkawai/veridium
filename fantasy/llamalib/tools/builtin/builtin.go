package builtin

import (
	"database/sql"

	"github.com/kawai-network/veridium/pkg/xlog"

	"github.com/kawai-network/veridium/fantasy/llamalib/tools"
)

// RegisterAll registers all builtin tools
func RegisterAll(registry *tools.ToolRegistry) error {
	return RegisterAllWithDB(registry, nil)
}

// RegisterAllWithDB registers all builtin tools with optional database connection
// Some tools (like image describe) require database access
func RegisterAllWithDB(registry *tools.ToolRegistry, sqlDB *sql.DB) error {
	xlog.Info("Registering builtin tools (yzma)...")

	// Register lobe-web-browsing (search, crawlSinglePage, crawlMultiPages)
	if err := RegisterWebBrowsing(registry); err != nil {
		return err
	}
	xlog.Info("✅ Registered tool", "tool", "lobe-web-browsing", "capabilities", "search, crawlSinglePage, crawlMultiPages")

	// Register lobe-local-system (file operations)
	if err := RegisterLocalSystem(registry); err != nil {
		return err
	}
	xlog.Info("✅ Registered tool", "tool", "lobe-local-system", "capabilities", "list, read, search, write, rename, move")

	// Register lobe-image-designer (DALL-E compatible)
	if err := RegisterImageDesigner(registry); err != nil {
		return err
	}
	xlog.Info("✅ Registered tool", "tool", "lobe-image-designer", "capabilities", "text2image")

	// Register lobe-code-interpreter (Python execution)
	if err := RegisterCodeInterpreter(registry); err != nil {
		return err
	}
	xlog.Info("✅ Registered tool", "tool", "lobe-code-interpreter", "capabilities", "python")

	// Register calculator
	if err := RegisterCalculator(registry); err != nil {
		return err
	}
	xlog.Info("✅ Registered tool", "tool", "calculator")

	// Register lobe-image-describe (requires DB for querying VL descriptions)
	if sqlDB != nil {
		if err := RegisterImageDescribe(registry, sqlDB); err != nil {
			return err
		}
		xlog.Info("✅ Registered tool", "tool", "lobe-image-describe", "capabilities", "getImageDescription")

		// Register lobe-video-describe (requires DB for querying Whisper transcriptions)
		if err := RegisterVideoDescribe(registry, sqlDB); err != nil {
			return err
		}
		xlog.Info("✅ Registered tool", "tool", "lobe-video-describe", "capabilities", "getVideoTranscription")
	} else {
		xlog.Warn("⚠️  Skipped tool", "tool", "lobe-image-describe", "reason", "no database connection")
		xlog.Warn("⚠️  Skipped tool", "tool", "lobe-video-describe", "reason", "no database connection")
	}

	return nil
}
