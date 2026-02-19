package batch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOutputPath_DefaultPattern(t *testing.T) {
	g := &Generator{}
	
	path := g.generateOutputPath("", 0, 42)
	assert.Equal(t, "batch_000.png", path)
	
	path = g.generateOutputPath("", 5, 123)
	assert.Equal(t, "batch_005.png", path)
}

func TestGenerateOutputPath_IndexPlaceholder(t *testing.T) {
	g := &Generator{}
	
	path := g.generateOutputPath("img_%d.png", 10, 42)
	assert.Equal(t, "img_10.png", path)
	
	path = g.generateOutputPath("image_%03d.png", 5, 42)
	assert.Equal(t, "image_005.png", path)
}

func TestGenerateOutputPath_SeedPlaceholder(t *testing.T) {
	g := &Generator{}
	
	// When only seed placeholder is used, index is appended before extension
	path := g.generateOutputPath("img_{seed}.png", 1, 12345)
	assert.Equal(t, "img_12345_1.png", path)
	
	path = g.generateOutputPath("seed_{seed}_img.png", 2, 999)
	assert.Equal(t, "seed_999_img_2.png", path)
}

func TestGenerateOutputPath_BothPlaceholders(t *testing.T) {
	g := &Generator{}
	
	path := g.generateOutputPath("img_%03d_seed{seed}.png", 7, 42)
	assert.Equal(t, "img_007_seed42.png", path)
	
	path = g.generateOutputPath("{seed}_%04d.png", 100, 1234)
	assert.Equal(t, "1234_0100.png", path)
}

func TestGenerateOutputPath_NoPlaceholders(t *testing.T) {
	g := &Generator{}
	
	// Pattern without placeholders should just return as-is with fmt.Sprintf behavior
	path := g.generateOutputPath("fixed_name.png", 1, 42)
	assert.True(t, strings.HasSuffix(path, ".png"))
}

func TestGenerateOutputPath_SeedPlaceholderOrder(t *testing.T) {
	g := &Generator{}
	
	// Test that seed placeholder is replaced before index formatting
	// This is the bug fix - previously fmt.Sprintf was called twice which would fail
	path := g.generateOutputPath("batch_%03d_seed{seed}.png", 5, 123)
	assert.Equal(t, "batch_005_seed123.png", path)
	
	// Test with seed first in pattern
	path = g.generateOutputPath("{seed}_%03d.png", 10, 999)
	assert.Equal(t, "999_010.png", path)
}

func TestGenerateOutputPath_MultipleSeedPlaceholders(t *testing.T) {
	g := &Generator{}
	
	// When only seed placeholders are used, index is appended before extension
	path := g.generateOutputPath("img_{seed}_seed{seed}.png", 1, 42)
	assert.Equal(t, "img_42_seed42_1.png", path)
}

func TestGenerateOutputPath_DotsInDirectoryName(t *testing.T) {
	g := &Generator{}
	
	// Pattern with dots in directory name but no extension - index appended at end
	path := g.generateOutputPath("outputs/v1.0/batch", 5, 42)
	assert.Equal(t, "outputs/v1.0/batch_5", path)
	
	// Pattern with dots in directory name AND extension - index inserted before extension
	path = g.generateOutputPath("outputs/v1.0/batch.png", 5, 42)
	assert.Equal(t, "outputs/v1.0/batch_5.png", path)
	
	// Pattern with multiple dots in directory name - only filename extension is preserved
	path = g.generateOutputPath("outputs/v1.0.2024/batch.png", 10, 99)
	assert.Equal(t, "outputs/v1.0.2024/batch_10.png", path)
	
	// Dots in directory should NOT be treated as extension
	path = g.generateOutputPath("outputs/v1.0/batch", 1, 1)
	assert.Equal(t, "outputs/v1.0/batch_1", path)
}

func TestGenerateOutputPath_NoExtension(t *testing.T) {
	g := &Generator{}
	
	// Pattern without extension - index appended at end
	path := g.generateOutputPath("outputs/batch", 5, 42)
	assert.Equal(t, "outputs/batch_5", path)
}
