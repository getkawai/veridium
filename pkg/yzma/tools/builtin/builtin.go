package builtin

import (
	"log"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// RegisterAll registers all builtin tools
func RegisterAll(registry *tools.ToolRegistry) error {
	log.Println("Registering builtin tools (yzma)...")

	// Register web search (simple version)
	if err := RegisterWebSearch(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: web_search")

	// Register lobe-web-browsing (full version with search, crawlSinglePage, crawlMultiPages)
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

	return nil
}

