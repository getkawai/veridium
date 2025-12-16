# Development Modes

Veridium menyediakan beberapa mode development untuk optimasi workflow Anda.

## 🔍 Pemahaman Penting

Ketika menjalankan `wails3 dev`, yang terjadi adalah:
1. ✅ `bun install` - Install dependencies
2. ✅ `bun run dev` - Start Vite dev server (background, port 9245) - **NEEDED for hot reload!**
3. ✅ `go mod tidy` - Tidy Go modules
4. ✅ **`bun run build:dev`** - Build frontend ke `dist/` (minify=false) - **⏱️ BOTTLENECK!**
5. ✅ `go build` - Compile Go binary (embeds dist)
6. ✅ Run aplikasi

**Key insights:**
- **Step #2 (Vite dev server)** - Tetap diperlukan untuk hot reload frontend
- **Step #4 (bun run build:dev)** - Ini yang di-skip untuk optimasi! (~10-20s saved)
- Go binary serve dari `dist/` tapi proxy ke Vite untuk hot reload

## 📋 Mode yang Tersedia

### 1. **Full Development Mode** (Default)
```bash
make dev    # Fresh start (hapus DB)
make devd   # Keep existing DB
```

**Kapan digunakan:**
- First run
- Setelah ada perubahan di frontend
- Setelah update dependencies

**Yang terjadi:**
- ✅ Install/update frontend dependencies
- ✅ Start Vite dev server (port 9245)
- ✅ **Build frontend dengan `bun run build:dev`** ⏱️ ~10-20s
- ✅ Generate bindings
- ✅ Compile Go backend
- ✅ Run aplikasi dengan hot reload

**Waktu startup:** ~15-35 detik

---

### 2. **Quick Development Mode** ⚡ (Recommended)
```bash
make dev-quick    # Fresh start (hapus DB)
make devd-quick   # Keep existing DB
```

**Kapan digunakan:**
- Iterasi cepat saat develop backend
- Frontend sudah pernah di-build sebelumnya
- Tidak ada perubahan di frontend

**Yang terjadi:**
- ✅ Check apakah `frontend/dist` exists
- ✅ Jika exists → **skip `bun run build:dev`** ⚡
- ✅ Jika tidak exists → fallback ke full mode
- ✅ Compile Go backend saja
- ✅ Run aplikasi

**Waktu startup:** ~5-8 detik (vs ~15-35 detik)

**Contoh workflow:**
```bash
# First run: build everything
make dev

# Subsequent runs: skip frontend build
make devd-quick  # Fast! ~5-8s
# ... edit Go code ...
make devd-quick  # Fast! ~5-8s
# ... edit Go code ...
make devd-quick  # Fast! ~5-8s
```

---

### 3. **Smart Development Mode** 🧠 (Smartest)
```bash
make dev-smart    # Fresh start (hapus DB)
make devd-smart   # Keep existing DB
```

**Kapan digunakan:**
- Workflow paling optimal
- Kombinasi reuse Vite + skip build

**Yang terjadi:**
1. Check apakah Vite dev server running di port 9245
   - ✅ Jika running → reuse Vite, skip semua frontend setup
2. Jika tidak, check apakah `frontend/dist` exists
   - ✅ Jika exists → skip `bun run build:dev`
3. Jika tidak ada keduanya → fallback ke full mode

**Waktu startup:** 
- ~3-5s (jika Vite running)
- ~5-8s (jika dist exists)
- ~15-35s (full build)

**Contoh workflow:**
```bash
# Terminal 1: Start Vite sekali (optional)
cd frontend && bun run dev -- --port 9245

# Terminal 2: Always use smart mode
make devd-smart  # Auto-detect best strategy
# ... edit code ...
make devd-smart  # Auto-detect best strategy
```

---

### 4. **Skip Build Mode** (Manual Control)
```bash
make dev-skip-build    # Fresh start (hapus DB)
make devd-skip-build   # Keep existing DB
```

**Kapan digunakan:**
- Anda yakin `frontend/dist` sudah up-to-date
- Ingin skip `bun run build:dev` secara eksplisit
- Maximum speed untuk iterasi backend

**Yang terjadi:**
- ⚠️ Check `frontend/dist` exists (error jika tidak ada)
- ✅ Skip `bun run build:dev` completely
- ✅ Compile Go backend
- ✅ Run aplikasi

**Waktu startup:** ~5-8 detik

---

### 5. **Skip Vite Dev Server Mode** (Advanced)
```bash
make dev-skip-frontend    # Fresh start (hapus DB)
make devd-skip-frontend   # Keep existing DB
```

**Kapan digunakan:**
- Anda ingin start Vite manual di terminal terpisah
- Custom Vite config atau debugging
- Kontrol penuh atas frontend

**Yang terjadi:**
- ⚠️ Tidak start Vite (Anda harus start manual)
- ✅ Masih run `bun run build:dev`
- ✅ Compile Go backend
- ✅ Run aplikasi

**Setup manual:**
```bash
# Terminal 1: Start Vite manual
cd frontend
bun run dev -- --port 9245 --strictPort

# Terminal 2: Start backend
make devd-skip-frontend
```

---

## 🎯 Rekomendasi Workflow

### ⭐ Best Practice (Paling Optimal):
```bash
# Always use smart mode - auto-detect best strategy
make devd-smart
# ... edit code ...
make devd-smart
# ... edit code ...
make devd-smart
```

### Untuk Iterasi Backend Cepat:
```bash
# First run: build everything
make dev

# Subsequent runs: skip frontend build
make devd-quick  # ~5-8s vs ~15-35s
# ... edit Go code ...
make devd-quick
# ... edit Go code ...
make devd-quick
```

### Untuk Development dengan Vite Manual:
```bash
# Terminal 1: Start Vite sekali
cd frontend && bun run dev -- --port 9245

# Terminal 2: Iterasi backend
make devd-skip-frontend
# ... edit code ...
make devd-skip-frontend
```

### Untuk Fresh Start:
```bash
make dev        # Full reset + hapus DB
make devd       # Full reset + keep DB
```

---

## 🔍 Troubleshooting

### Port 9245 sudah digunakan?
```bash
# Kill port manual
killport 9245

# Atau gunakan lsof
lsof -ti:9245 | xargs kill -9
```

### Frontend tidak connect ke backend?
1. Pastikan Vite running di port 9245
2. Check `FRONTEND_DEVSERVER_URL` di environment
3. Restart dengan full mode: `make dev`

### Hot reload tidak bekerja?
- **Frontend changes**: Harus via Vite HMR (instant)
- **Go changes**: Auto-rebuild + restart (~3-5 detik)
- Check `build/config.yml` untuk file watcher settings

---

## 📊 Perbandingan Performa

| Mode | Startup Time | Frontend Build | Vite Server | Use Case |
|------|--------------|----------------|-------------|----------|
| `make dev` | ~15-35s | ✅ Yes | ✅ Start | First run, frontend changes |
| `make devd` | ~15-35s | ✅ Yes | ✅ Start | Full restart + keep DB |
| `make dev-quick` | ~5-8s | ❌ Skip | ✅ Start | Backend iterations |
| `make devd-quick` | ~5-8s | ❌ Skip | ✅ Start | Backend iterations + keep DB |
| `make dev-smart` | ~3-8s | 🤖 Auto | 🤖 Auto | Smartest (recommended!) |
| `make devd-smart` | ~3-8s | 🤖 Auto | 🤖 Auto | Smartest + keep DB |
| `make dev-skip-build` | ~5-8s | ❌ Skip | ✅ Start | Explicit skip build |
| `make dev-skip-frontend` | ~15-35s | ✅ Yes | ❌ Manual | Manual Vite control |

**Legend:**
- ✅ Yes = Always execute
- ❌ Skip = Never execute
- 🤖 Auto = Auto-detect and decide

---

## 🛠️ Technical Details

### What Actually Happens in `wails3 dev`

**Full Mode (`make dev`):**
```
1. bun install                    (~2-5s)
2. bun run dev (background)       (~1s)
3. go mod tidy                    (~1-2s)
4. bun run build:dev              (~10-20s) ⏱️ BOTTLENECK!
5. wails3 generate bindings       (~2-3s)
6. go build                       (~3-5s)
7. Run app                        (~1s)
────────────────────────────────────────
Total: ~15-35s
```

**Quick Mode (`make dev-quick`):**
```
1. Check frontend/dist exists     (~0s)
2. bun install                    (~2-5s)
3. bun run dev (background)       (~1s)
4. go mod tidy                    (~1-2s)
5. [SKIP] bun run build:dev       (saved ~10-20s!) ⚡
6. wails3 generate bindings       (~2-3s)
7. go build                       (~3-5s)
8. Run app                        (~1s)
────────────────────────────────────────
Total: ~5-8s (3-5x faster!)
```

### Configuration Files

1. **`build/config.yml`** - Default full mode
   ```yaml
   executes:
     - bun install
     - bun run dev (Vite server)
     - go mod tidy
     - wails3 task build          # ← includes bun run build:dev
     - run app
   ```

2. **`build/config-skip-frontend-build.yml`** - Skip frontend build (OPTIMIZED!)
   ```yaml
   executes:
     - bun install
     - bun run dev (Vite server)  # ← Still needed for hot reload!
     - go mod tidy
     - wails3 task build:skip     # ← SKIP bun run build:dev
     - run app
   ```
   **Key optimization:** Skip `bun run build:dev` (~10-20s saved!)

3. **`build/config-skip-frontend.yml`** - Skip Vite dev server (Manual control)
   ```yaml
   executes:
     # Skip: bun install
     # Skip: bun run dev
     - go mod tidy
     - wails3 task build          # ← Still runs bun run build:dev
     - run app
   ```
   **Use case:** When you manually start Vite in separate terminal

### File Watcher Configuration
Lihat `build/config.yml` untuk:
- Watched extensions: `*.go`
- Ignored directories: `.git`, `node_modules`, `frontend`, `bin`
- Debounce time: 1000ms

### Architecture
```
┌─────────────┐      Proxy      ┌──────────────────┐
│  Wails App  │ ←──────────────→ │ Vite Dev Server  │
│   (Go)      │                  │ (localhost:9245) │
└─────────────┘                  └──────────────────┘
       ↓                                   ↓
       └────────────→ WebView ←────────────┘
                  (Render UI + Bridge)
```

**Hot Reload:**
- Frontend: Vite HMR (instant, no restart)
- Backend: File watcher → rebuild → restart (~5-8s with quick mode)

