package builtin

import (
	"log"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// RegisterAll registers all builtin tools
func RegisterAll(registry *tools.ToolRegistry) error {
	log.Println("Registering builtin tools (yzma)...")
	
	// Register web search
	if err := RegisterWebSearch(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: web_search")
	
	// Register calculator
	if err := RegisterCalculator(registry); err != nil {
		return err
	}
	log.Println("✅ Registered: calculator")
	
	return nil
}

