/*
 * Memory RAG Integration Test CLI
 *
 * Tests the REAL Veridium application's Memory services.
 * Uses the same initialization code as main.go via internal/app.
 *
 * Usage:
 *   go run ./cmd/test-memory
 *   go run ./cmd/test-memory -v
 *   go run ./cmd/test-memory -test semantic
 */
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/tools"
	"github.com/kawai-network/veridium/internal/app"
	"github.com/kawai-network/veridium/internal/services"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

var (
	verbose    = flag.Bool("v", false, "verbose output")
	testFilter = flag.String("test", "", "run specific test")
)

type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
	Details  string
	Error    error
}

func main() {
	flag.Parse()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        Memory RAG Integration Test                         ║")
	fmt.Println("║        Testing REAL Veridium Application                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Clean data for fresh test
	fmt.Println("Cleaning data directory...")
	os.RemoveAll("data")
	fmt.Printf("%s✓ Data cleaned%s\n\n", colorGreen, colorReset)

	// Initialize using the SAME code as main.go
	fmt.Println("Initializing services (same as main.go)...")
	appCtx := app.NewContext()
	defer appCtx.Cleanup()

	if err := appCtx.InitAll(); err != nil {
		fmt.Printf("%s✗ Init failed: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Verify what we have
	fmt.Printf("\n%sServices Status:%s\n", colorCyan, colorReset)
	fmt.Printf("  DB: %v\n", appCtx.DB != nil)
	fmt.Printf("  DuckDB: %v\n", appCtx.DuckDBStore != nil)
	fmt.Printf("  Embedder: %v (dims: %d)\n", appCtx.Embedder != nil, safeGetDims(appCtx))
	fmt.Printf("  ChatModel (LLM): %v\n", appCtx.ChatModel != nil)
	fmt.Printf("  MemoryService: %v\n", appCtx.MemoryService != nil)
	fmt.Printf("  MemoryEnrichment: %v\n", appCtx.MemoryEnrichment != nil)
	fmt.Printf("  MemoryIntegration: %v\n", appCtx.MemoryIntegration != nil)
	fmt.Println()

	// Check prerequisites
	if appCtx.MemoryService == nil {
		fmt.Printf("%s✗ MemoryService not initialized - cannot run tests%s\n", colorRed, colorReset)
		os.Exit(1)
	}

	// Run tests
	tests := []struct {
		name string
		fn   func(context.Context, *app.Context) TestResult
	}{
		{"Semantic Search", testSemanticSearch},
		{"Memory Enrichment", testMemoryEnrichment},
		{"Auto-Archive", testAutoArchive},
		{"Tool Execution", testToolExecution},
		{"End-to-End Flow", testEndToEnd},
		{"Vector Similarity Ranking", testVectorSimilarityRanking},
	}

	var results []TestResult
	passCount := 0

	for _, test := range tests {
		if *testFilter != "" && !strings.Contains(strings.ToLower(test.name), strings.ToLower(*testFilter)) {
			continue
		}

		fmt.Printf("%s▶ Running: %s%s\n", colorCyan, test.name, colorReset)
		start := time.Now()

		result := test.fn(context.Background(), appCtx)
		result.Name = test.name
		result.Duration = time.Since(start)

		if result.Passed {
			passCount++
			fmt.Printf("  %s✓ PASSED%s (%v)\n", colorGreen, colorReset, result.Duration)
		} else {
			fmt.Printf("  %s✗ FAILED%s (%v)\n", colorRed, colorReset, result.Duration)
			if result.Error != nil {
				fmt.Printf("    Error: %v\n", result.Error)
			}
		}

		if *verbose && result.Details != "" {
			for _, line := range strings.Split(result.Details, "\n") {
				if line != "" {
					fmt.Printf("    %s\n", line)
				}
			}
		}

		results = append(results, result)
		fmt.Println()
	}

	// Summary
	fmt.Println("════════════════════════════════════════════════════════════")
	total := len(results)
	if total == 0 {
		fmt.Println("No tests matched filter")
		os.Exit(1)
	}

	if passCount == total {
		fmt.Printf("%s✓ All %d tests passed (100%%)%s\n", colorGreen, total, colorReset)
	} else {
		fmt.Printf("%s✗ %d/%d tests passed (%.0f%%)%s\n",
			colorYellow, passCount, total, float64(passCount)/float64(total)*100, colorReset)
		os.Exit(1)
	}
}

func safeGetDims(ctx *app.Context) int {
	if ctx.Embedder == nil {
		return 0
	}
	return ctx.Embedder.Dimensions()
}

// =============================================================================
// TESTS
// =============================================================================

func testSemanticSearch(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	memories := []struct {
		title, summary, topic string
	}{
		{"User Location", "Saya tinggal di Jakarta, Indonesia", "location"},
		{"Programming Skills", "Saya suka programming dengan Go dan Python", "programming"},
		{"Color Preference", "Warna favorit saya adalah biru", "color"},
	}

	createdIDs := make(map[string]string)
	for _, m := range memories {
		mem, err := appCtx.MemoryService.CreateMemory(ctx, &services.Memory{
			Category: services.MemoryCategoryFact,
			Title:    m.title,
			Summary:  m.summary,
		})
		if err != nil {
			return TestResult{Error: fmt.Errorf("create memory: %w", err)}
		}
		createdIDs[mem.ID] = m.topic
		details.WriteString(fmt.Sprintf("Created [%s]: %s\n", m.topic, m.title))
	}

	testCases := []struct {
		query, expectedTopic string
	}{
		{"Dimana user tinggal?", "location"},
		{"Bahasa programming apa yang digunakan?", "programming"},
		{"Apa warna favorit user?", "color"},
	}

	passCount := 0
	for _, tc := range testCases {
		results, err := appCtx.MemoryService.SemanticSearch(ctx, tc.query, 10)
		if err != nil {
			details.WriteString(fmt.Sprintf("✗ '%s': error: %v\n", tc.query, err))
			continue
		}

		found := false
		for i, r := range results {
			if i >= 3 {
				break
			}
			if topic := createdIDs[r.Memory.ID]; topic == tc.expectedTopic {
				found = true
				passCount++
				details.WriteString(fmt.Sprintf("✓ '%s' -> rank %d (%.3f)\n", tc.query, i+1, r.Similarity))
				break
			}
		}
		if !found {
			details.WriteString(fmt.Sprintf("✗ '%s' -> expected %s not in top 3\n", tc.query, tc.expectedTopic))
		}
	}

	return TestResult{
		Passed:  passCount >= 2,
		Details: details.String(),
	}
}

func testMemoryEnrichment(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	messages := []fantasy.Message{
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Nama saya Budi, saya tinggal di Jakarta"}}},
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya suka programming dengan Go"}}},
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Warna favorit saya biru"}}},
	}

	result, err := appCtx.MemoryEnrichment.EnrichMessages(ctx, messages)
	if err != nil {
		return TestResult{Error: fmt.Errorf("enrich: %w", err)}
	}

	details.WriteString(fmt.Sprintf("Extracted %d facts\n", result.FactCount))

	memories, _ := appCtx.MemoryService.ListMemories(ctx, 100, 0)

	patterns := []string{"budi", "jakarta", "programming", "biru"}
	found := 0
	for _, p := range patterns {
		for _, mem := range memories {
			if strings.Contains(strings.ToLower(mem.Summary+" "+mem.Title), p) {
				found++
				details.WriteString(fmt.Sprintf("✓ Found: %s\n", p))
				break
			}
		}
	}

	return TestResult{
		Passed:  found >= 2,
		Details: details.String(),
	}
}

func testAutoArchive(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	messages := []fantasy.Message{
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Nama saya Ahmad"}}},
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya suka programming"}}},
	}
	for i := 0; i < 10; i++ {
		messages = append(messages, fantasy.Message{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: fmt.Sprintf("Filler %d", i)}},
		})
	}

	initial := len(messages)
	remaining, err := appCtx.MemoryEnrichment.AutoArchive(ctx, messages, services.BufferConfig{
		MaxBufferSize:    15,
		ArchiveBatchSize: 3,
		ArchiveThreshold: 8,
	})
	if err != nil {
		return TestResult{Error: err}
	}

	details.WriteString(fmt.Sprintf("Buffer: %d -> %d\n", initial, len(remaining)))
	return TestResult{
		Passed:  len(remaining) < initial,
		Details: details.String(),
	}
}

func testToolExecution(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	mem, _ := appCtx.MemoryService.CreateMemory(ctx, &services.Memory{
		Category: services.MemoryCategoryFact,
		Title:    "Test User",
		Summary:  "User name is Alice, works as developer in Jakarta",
	})
	details.WriteString(fmt.Sprintf("Created: %s\n", mem.ID[:8]))

	registry := tools.NewToolRegistry()
	appCtx.MemoryIntegration.RegisterMemoryTool(registry)

	tool, _ := registry.Get("search_memory")
	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test",
		Name:  "search_memory",
		Input: `{"query": "siapa nama user", "limit": 5}`,
	})
	if err != nil {
		return TestResult{Error: err}
	}

	details.WriteString(fmt.Sprintf("Response: %d bytes\n", len(resp.Content)))

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Content), &data)

	hasInfo := strings.Contains(strings.ToLower(resp.Content), "alice")
	if hasInfo {
		details.WriteString("✓ Found 'alice'\n")
	}

	return TestResult{
		Passed:  !resp.IsError && hasInfo,
		Details: details.String(),
	}
}

func testEndToEnd(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	details.WriteString("Phase 1: Chat\n")
	chat := []fantasy.Message{
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Nama saya Dewi"}}},
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya tinggal di Bandung"}}},
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya bekerja sebagai data scientist"}}},
	}

	result, _ := appCtx.MemoryEnrichment.EnrichMessages(ctx, chat)
	details.WriteString(fmt.Sprintf("Enriched: %d facts\n", result.FactCount))

	details.WriteString("\nPhase 2: Recall\n")
	queries := []struct{ query, expected string }{
		{"siapa nama user", "dewi"},
		{"dimana user tinggal", "bandung"},
		{"apa pekerjaan user", "data scientist"},
	}

	success := 0
	for _, q := range queries {
		results, _ := appCtx.MemoryService.SemanticSearch(ctx, q.query, 5)
		found := false
		for _, r := range results {
			if strings.Contains(strings.ToLower(r.Memory.Summary), q.expected) {
				found = true
				break
			}
		}
		if found {
			success++
			details.WriteString(fmt.Sprintf("✓ '%s' -> '%s'\n", q.query, q.expected))
		} else {
			details.WriteString(fmt.Sprintf("✗ '%s' -> missing '%s'\n", q.query, q.expected))
		}
	}

	return TestResult{
		Passed:  success >= 2,
		Details: details.String(),
	}
}

func testVectorSimilarityRanking(ctx context.Context, appCtx *app.Context) TestResult {
	var details strings.Builder

	mem1, _ := appCtx.MemoryService.CreateMemory(ctx, &services.Memory{
		Category: services.MemoryCategoryFact,
		Title:    "Programming",
		Summary:  "User loves coding with Go and Python",
	})
	mem2, _ := appCtx.MemoryService.CreateMemory(ctx, &services.Memory{
		Category: services.MemoryCategoryFact,
		Title:    "Location",
		Summary:  "User lives in Jakarta Indonesia",
	})

	// Query programming
	results, _ := appCtx.MemoryService.SemanticSearch(ctx, "programming coding software", 10)

	var progSim, locSim float64
	for _, r := range results {
		if r.Memory.ID == mem1.ID {
			progSim = r.Similarity
		} else if r.Memory.ID == mem2.ID {
			locSim = r.Similarity
		}
	}

	details.WriteString(fmt.Sprintf("Query 'programming': prog=%.3f, loc=%.3f\n", progSim, locSim))
	progHigher := progSim > locSim

	// Query location
	results2, _ := appCtx.MemoryService.SemanticSearch(ctx, "dimana tinggal kota", 10)
	for _, r := range results2 {
		if r.Memory.ID == mem1.ID {
			progSim = r.Similarity
		} else if r.Memory.ID == mem2.ID {
			locSim = r.Similarity
		}
	}

	details.WriteString(fmt.Sprintf("Query 'location': prog=%.3f, loc=%.3f\n", progSim, locSim))
	locHigher := locSim > progSim

	if progHigher {
		details.WriteString("✓ Programming ranked higher for programming query\n")
	}
	if locHigher {
		details.WriteString("✓ Location ranked higher for location query\n")
	}

	return TestResult{
		Passed:  progHigher && locHigher,
		Details: details.String(),
	}
}
