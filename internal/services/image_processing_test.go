package services_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

func TestHybridImageProcessing_OCRPath(t *testing.T) {
	// Use wails.png which has "WAILS" text - should trigger OCR fast path
	imagePath := "/Users/yuda/github.com/kawai-network/veridium/frontend/dist/wails.png"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Skip("Test image not found, skipping")
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "image_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy image to temp dir
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read image: %v", err)
	}
	testImagePath := filepath.Join(tempDir, "test_image.png")
	if err := os.WriteFile(testImagePath, imageData, 0644); err != nil {
		t.Fatalf("Failed to write test image: %v", err)
	}

	// Initialize test database
	dbPath := filepath.Join(tempDir, "test.db")
	dbService, err := database.NewServiceWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer dbService.Close()

	ctx := context.Background()
	queries := db.New(dbService.DB())
	const testUserID = "DEFAULT_LOBE_CHAT_USER"

	// Initialize file processor (without VL model - test OCR only path)
	fileLoader := services.NewFileLoader()
	processor := services.NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		nil, // No RAG processor
		nil, // No VL model - force OCR path
		nil, // No whisper
	)

	// Process image
	t.Log("Processing image (expecting OCR fast path)...")
	start := time.Now()

	result, err := processor.ProcessFile(ctx, services.ProcessFileRequest{
		FilePath:  testImagePath,
		Filename:  "test_image.png",
		UserID:    testUserID,
		Source:    testImagePath,
		EnableRAG: false,
		IsShared:  false,
	})
	if err != nil {
		t.Fatalf("Failed to process image: %v", err)
	}

	t.Logf("Image processed: FileID=%s, DocumentID=%s", result.FileID, result.DocumentID)

	// Wait for async processing (should be very fast with OCR)
	deadline := time.Now().Add(10 * time.Second)
	var content string

	for time.Now().Before(deadline) {
		doc, err := queries.GetDocumentByFileID(ctx, sql.NullString{String: result.FileID, Valid: true})
		if err == nil && doc.Content.Valid {
			if strings.Contains(doc.Content.String, "OCR Text") ||
				strings.Contains(doc.Content.String, "Image Description") {
				content = doc.Content.String
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	elapsed := time.Since(start)

	if content == "" {
		t.Fatal("Image content not generated within timeout")
	}

	t.Logf("Processing completed in %v", elapsed)
	t.Logf("Content type: %s", getContentType(content))
	t.Logf("Content length: %d chars", len(content))

	// Check if it used OCR (fast path)
	if strings.Contains(content, "OCR Text (Tesseract") {
		t.Log("✅ Used OCR fast path")
	} else if strings.Contains(content, "Image Description") {
		t.Log("Used VL slow path (OCR didn't find enough text)")
	}

	// Verify content contains expected text
	if strings.Contains(strings.ToUpper(content), "WAILS") {
		t.Log("✅ OCR correctly extracted 'WAILS' text")
	}

	// For OCR path, should be very fast (< 2 seconds)
	if elapsed < 2*time.Second && strings.Contains(content, "OCR Text") {
		t.Logf("✅ Fast path confirmed: %v", elapsed)
	}
}

func getContentType(content string) string {
	if strings.Contains(content, "OCR Text (Tesseract + LLM cleanup)") {
		return "OCR (Tesseract + LLM cleanup)"
	}
	if strings.Contains(content, "OCR Text (Tesseract") {
		return "OCR (Tesseract)"
	}
	if strings.Contains(content, "Image Description (VL Model)") {
		return "VL Model"
	}
	if strings.Contains(content, "Image Description (AI Generated)") {
		return "AI Generated"
	}
	return "Unknown"
}

// TestOCRWithLLMCleanup tests OCR extraction with LLM cleanup on a document image
func TestOCRWithLLMCleanup(t *testing.T) {
	// Use docparsing_example1.jpg which has table and text - good for OCR cleanup test
	imagePath := "/Users/yuda/github.com/kawai-network/veridium/internal/llama/docparsing_example1.jpg"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Skip("Test image not found, skipping")
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "ocr_cleanup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy image to temp dir
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("Failed to read image: %v", err)
	}
	testImagePath := filepath.Join(tempDir, "test_document.jpg")
	if err := os.WriteFile(testImagePath, imageData, 0644); err != nil {
		t.Fatalf("Failed to write test image: %v", err)
	}

	// Initialize test database
	dbPath := filepath.Join(tempDir, "test.db")
	dbService, err := database.NewServiceWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer dbService.Close()

	ctx := context.Background()
	queries := db.New(dbService.DB())
	const testUserID = "DEFAULT_LOBE_CHAT_USER"

	// Initialize file processor
	fileLoader := services.NewFileLoader()
	processor := services.NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		nil, // No RAG processor
		nil, // No VL model
		nil, // No whisper
	)

	// Create TaskRouter with Zhipu for OCR cleanup
	toolRegistry := tools.NewToolRegistry()
	config := llm.GetDefaultDevConfig()
	taskRouter := llm.BuildTaskRouter(config, toolRegistry, nil)

	// Connect LLM provider for OCR cleanup
	llmAdapter := services.NewTaskRouterAdapter(taskRouter)
	processor.SetLLMProvider(llmAdapter)

	// Process image
	t.Log("Processing document image with OCR + LLM cleanup...")
	start := time.Now()

	result, err := processor.ProcessFile(ctx, services.ProcessFileRequest{
		FilePath:  testImagePath,
		Filename:  "test_document.jpg",
		UserID:    testUserID,
		Source:    testImagePath,
		EnableRAG: false,
		IsShared:  false,
	})
	if err != nil {
		t.Fatalf("Failed to process image: %v", err)
	}

	t.Logf("Image processed: FileID=%s, DocumentID=%s", result.FileID, result.DocumentID)

	// Wait for async processing (OCR + LLM cleanup may take up to 60 seconds with remote API)
	deadline := time.Now().Add(60 * time.Second)
	var content string

	for time.Now().Before(deadline) {
		doc, err := queries.GetDocumentByFileID(ctx, sql.NullString{String: result.FileID, Valid: true})
		if err == nil && doc.Content.Valid {
			if strings.Contains(doc.Content.String, "OCR Text") ||
				strings.Contains(doc.Content.String, "Image Description") {
				content = doc.Content.String
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	elapsed := time.Since(start)

	if content == "" {
		t.Fatal("Image content not generated within timeout")
	}

	contentType := getContentType(content)
	t.Logf("Processing completed in %v", elapsed)
	t.Logf("Content type: %s", contentType)
	t.Logf("Content length: %d chars", len(content))

	// Check if LLM cleanup was used
	if strings.Contains(content, "LLM cleanup") {
		t.Log("✅ LLM cleanup was applied")
	} else {
		t.Log("⚠️  LLM cleanup was NOT applied (raw OCR or VL)")
	}

	// Check for markdown formatting (tables, headers, etc.)
	if strings.Contains(content, "|") || strings.Contains(content, "##") || strings.Contains(content, "**") {
		t.Log("✅ Markdown formatting detected")
	}

	// Check content contains expected text
	if strings.Contains(content, "Qwen") || strings.Contains(content, "MMLU") {
		t.Log("✅ OCR correctly extracted benchmark table content")
	}

	// Print first 500 chars of content for inspection
	preview := content
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	t.Logf("\n--- Content Preview ---\n%s\n--- End Preview ---", preview)
}
