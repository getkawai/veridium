# Token Counter Example

Contoh program untuk menghitung jumlah token dalam teks menggunakan library Yzma.

## Fitur

- **Token Counting**: Menghitung jumlah token dari teks input
- **Special Tokens**: Menunjukkan perbedaan tokenisasi dengan dan tanpa special tokens
- **Detailed Analysis**: Menampilkan breakdown lengkap termasuk karakter per token
- **Multiple Examples**: Contoh tokenisasi untuk berbagai jenis teks
- **Special Token Info**: Menampilkan informasi special tokens yang digunakan model

## Instalasi

 Pastikan Yzma sudah di-install dan model GGUF tersedia:

```bash
# Build Yzma
go build

# Download model (contoh)
./yzma installer install microsoft/DialoGPT-medium
```

## Penggunaan

### Mode Examples (Tanpa teks input)

```bash
cd examples/token_counter
go run main.go /path/to/your/model.gguf
```

### Mode Custom Text

```bash
go run main.go /path/to/your/model.gguf "Your text here"
```

### Contoh Lengkap

```bash
# Examples mode
go run main.go models/llama-2-7b-chat.Q4_0.gguf

# Custom text mode  
go run main.go models/llama-2-7b-chat.Q4_0.gguf "Hello world, how are you?"
```

## Output Contoh

```
=== YZMA TOKEN COUNTER ===
Model: models/llama-2-7b-chat.Q4_0.gguf
Vocab Size: 32000 tokens

=== TOKEN COUNTING EXAMPLES ===

Example 1: "Hello world"
  Characters: 11
  Tokens (with special): 3
  Tokens (without special): 2
  Characters per token (with special): 3.67

Example 2: "What is the capital of France?"
  Characters: 30
  Tokens (with special): 7
  Tokens (without special): 6
  Characters per token (with special): 4.29
...

=== SPECIAL TOKENS ===
BOS (Beginning of Sentence): 1
EOS (End of Sentence): 2
EOT (End of Turn): 3
SEP (Separator): 4
NL (New Line): 13
```

## Parameter Tokenisasi

Program ini menggunakan dua mode tokenisasi:

1. **With Special Tokens** (`addSpecial=true, parseSpecial=true`)
   - Menambah BOS, EOS tokens
   - Memproses special tokens dalam input

2. **Without Special Tokens** (`addSpecial=false, parseSpecial=false`)
   - Tidak menambah special tokens
   - Treating all text as regular content

## Technical Details

- Menggunakan `llama.Tokenize()` untuk konversi teks ke tokens
- Menggunakan `llama.TokenToPiece()` untuk konversi tokens kembali ke teks
- Mendukung berbagai model GGUF format
- Output karakter per token untuk analisis efisiensi tokenisasi
