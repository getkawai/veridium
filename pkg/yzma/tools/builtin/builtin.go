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

	// Register calculator
	if err := RegisterCalculator(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: calculator")

	return nil
}

