package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// TestLLMProvider implements LLMProvider interface for testing
type TestLLMProvider struct {
	router *llm.TaskRouter
}

func (p *TestLLMProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	if p.router == nil {
		return "", fmt.Errorf("router not initialized")
	}

	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: prompt,
		},
	}

	// Use transcript_cleanup task type
	resp, err := p.router.GenerateWithoutTools(ctx, llm.TaskTranscriptCleanup, messages)
	if err != nil {
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	if resp == nil || resp.Content == "" {
		return "", fmt.Errorf("empty response from LLM")
	}

	return resp.Content, nil
}

func TestTranscriptCleanupWithZhipu(t *testing.T) {
	// Sample transcript with errors (from Whisper)
	sampleTranscript := `**[Segment 1]**
iya oke guys, balik lagi sama gue ya karena data gue bisa punggian, oke ya hari ini Jum 8, bener, 5 Desember 2025 ya
IHSK dalam posisi hari ini oke, iya turun timis gak apa-apa, oke ya gue pengen bahas aturan baru IPO ya John ya
aturan baru IPO dan sebenernya kalau lo liat-liat aturan ini sih lebih ngontongnya nge-retail ya sebenernya ya

**[Segment 2]**
jadi kalau lo baca John OJK resmi menerbitkan SOJK nomor 25 ya blablabla yang mengatur perubahan alokasi efek menawaran umum untuk retail
kalau satu terliunan 2,5% gitu ya lokasinya ya
berarti peraturan baru ini lebih merogikan kepada para pemain yang di atas 100 juta`

	// Create task router with Zhipu config
	config := llm.GetDefaultDevConfig()
	router := llm.BuildTaskRouter(config, nil, nil)

	// Use TaskRouterAdapter which adds system prompt and handles task detection
	adapter := NewTaskRouterAdapter(router)

	// Create file processor service with adapter
	fps := &FileProcessorService{}
	fps.SetLLMProvider(adapter)

	// Test cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	t.Log("Starting transcript cleanup test...")
	t.Logf("Original transcript length: %d", len(sampleTranscript))
	t.Log("---ORIGINAL---")
	t.Log(sampleTranscript)
	t.Log("--------------")

	start := time.Now()
	
	// Call the cleanupTranscription method directly
	cleaned, err := fps.cleanupTranscription(ctx, sampleTranscript)
	
	elapsed := time.Since(start)
	t.Logf("Cleanup took: %v", elapsed)

	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	t.Logf("Cleaned transcript length: %d", len(cleaned))
	t.Log("---CLEANED---")
	t.Log(cleaned)
	t.Log("-------------")

	// Verify some corrections were made
	if cleaned == sampleTranscript {
		t.Error("No corrections were made - transcript unchanged")
	}

	// Check for specific corrections
	corrections := map[string]string{
		"IHSK":       "IHSG",
		"terliunan":  "triliun",
		"merogikan":  "merugikan",
		"SOJK":       "SEOJK",
		"timis":      "tipis",
		"lokasinya":  "alokasinya",
	}

	for wrong, right := range corrections {
		if containsWord(cleaned, wrong) && !containsWord(cleaned, right) {
			t.Logf("Warning: '%s' should be corrected to '%s'", wrong, right)
		}
	}
}

func containsWord(text, word string) bool {
	return len(text) > 0 && len(word) > 0 && 
		(text == word || 
		 len(text) >= len(word) && 
		 (text[:len(word)] == word || 
		  text[len(text)-len(word):] == word ||
		  containsSubstring(text, " "+word+" ") ||
		  containsSubstring(text, " "+word+",") ||
		  containsSubstring(text, " "+word+".")))
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestTranscriptCleanupPromptDetection tests that the adapter detects transcript cleanup
func TestTranscriptCleanupPromptDetection(t *testing.T) {
	config := llm.GetDefaultDevConfig()
	router := llm.BuildTaskRouter(config, nil, nil)
	
	adapter := NewTaskRouterAdapter(router)
	
	// Test transcript detection
	transcriptPrompt := "Fix typos in this Indonesian transcription from Whisper speech-to-text..."
	taskType := adapter.detectTaskType(transcriptPrompt)
	
	if taskType != llm.TaskTranscriptCleanup {
		t.Errorf("Expected TaskTranscriptCleanup, got %s", taskType)
	}
	
	// Test OCR detection (default)
	ocrPrompt := "Clean up this OCR text from document..."
	taskType = adapter.detectTaskType(ocrPrompt)
	
	if taskType != llm.TaskOCRCleanup {
		t.Errorf("Expected TaskOCRCleanup, got %s", taskType)
	}
}

// TestZhipuDirectCall tests Zhipu API directly
func TestZhipuDirectCall(t *testing.T) {
	config := llm.GetDefaultDevConfig()
	
	// Check if API key is configured
	if config.TranscriptCleanup.APIKey == "" {
		t.Skip("No Zhipu API key configured")
	}
	
	router := llm.BuildTaskRouter(config, nil, nil)
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// Simple test prompt
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "Fix typos: IHSK→IHSG, terlion→triliun. Text: IHSK naik 100 terlion",
		},
	}
	
	t.Log("Sending request to Zhipu...")
	start := time.Now()
	
	resp, err := router.GenerateWithoutTools(ctx, llm.TaskTranscriptCleanup, messages)
	
	elapsed := time.Since(start)
	t.Logf("Request took: %v", elapsed)
	
	if err != nil {
		t.Fatalf("Zhipu request failed: %v", err)
	}
	
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Response content: %s", resp.Content)
	t.Logf("Finish reason: %s", resp.FinishReason)
	
	if resp.Content == "" {
		t.Error("Response content is empty")
	}
}
