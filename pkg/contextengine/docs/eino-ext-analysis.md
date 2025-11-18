# Eino-Ext Integration Analysis

## Overview

`github.com/cloudwego/eino-ext` menyediakan berbagai ekstensi untuk framework Eino, termasuk komponen-komponen yang bisa digunakan untuk memperbaiki implementasi `context-engine-go`.

## Komponen yang Relevan

### 1. ChatTemplate (DefaultChatTemplate)

**Kegunaan:**
- Bisa digunakan untuk menggantikan implementasi manual template processing di `InputTemplateProcessor`
- Menyediakan template engine yang lebih powerful daripada regex sederhana

**Current Implementation:**
- Menggunakan regex sederhana: `{{text}}` replacement
- TypeScript version menggunakan lodash template dengan interpolate pattern

**Potential Improvement:**
- Menggunakan Eino's ChatTemplate untuk template processing yang lebih robust
- Mendukung variabel template yang lebih kompleks

**Trade-off:**
- ✅ Lebih powerful dan standardized
- ✅ Terintegrasi dengan Eino ecosystem
- ❌ Mungkin overkill untuk use case sederhana (hanya `{{text}}`)
- ❌ Menambah dependency

### 2. JSONMessageParser (Lambda)

**Kegunaan:**
- Parsing JSON message content
- Bisa berguna untuk MessageContentProcessor

**Current Implementation:**
- Manual parsing di berbagai processors
- Tidak ada JSON parsing khusus

**Potential Improvement:**
- Menggunakan JSONMessageParser untuk parsing structured message content
- Lebih robust error handling

**Trade-off:**
- ✅ Standardized parsing
- ✅ Better error handling
- ❌ Mungkin tidak diperlukan jika message sudah dalam format schema.Message

## Rekomendasi

### ✅ **RECOMMENDED: Gunakan ChatTemplate untuk InputTemplateProcessor**

**Alasan:**
1. Template processing adalah core functionality yang penting
2. Eino ChatTemplate lebih robust daripada regex sederhana
3. Terintegrasi dengan Eino ecosystem
4. Bisa mendukung template yang lebih kompleks di masa depan

**Implementation Plan:**
```go
// Instead of manual regex replacement:
templatePattern := regexp.MustCompile(`{{\s*text\s*}}`)

// Use Eino ChatTemplate:
import "github.com/cloudwego/eino-ext/prompt"

chatTemplate := prompt.NewDefaultChatTemplate(config.InputTemplate)
// Use in workflow as ChatTemplateNode
```

### ⚠️ **OPTIONAL: JSONMessageParser**

**Alasan:**
1. Saat ini tidak ada kebutuhan khusus untuk JSON parsing
2. Message sudah dalam format schema.Message yang terstruktur
3. Bisa ditambahkan nanti jika diperlukan

## Implementation Steps

1. **Add eino-ext dependency:**
   ```bash
   go get github.com/cloudwego/eino-ext
   ```

2. **Refactor InputTemplateProcessor:**
   - Replace regex-based template dengan ChatTemplateNode
   - Update graph builder untuk menggunakan ChatTemplateNode

3. **Testing:**
   - Ensure backward compatibility
   - Test dengan berbagai template formats

## Conclusion

**YES, bisa utilize eino-ext**, terutama untuk:
- ✅ ChatTemplate untuk InputTemplateProcessor (RECOMMENDED)
- ⚠️ JSONMessageParser (OPTIONAL, bisa ditambahkan nanti)

**Priority:**
1. High: ChatTemplate integration
2. Low: JSONMessageParser (jika diperlukan)

