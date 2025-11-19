# Phase 1: Quick Start Guide

> **Get started testing Phase 1 in 5 minutes**

---

## 🚀 Quick Test (5 minutes)

### Step 1: Verify Flag is OFF (Default)

```bash
# Check feature flag
cat frontend/src/config/features.ts | grep USE_BACKEND_CHAT
```

Should show: `USE_BACKEND_CHAT: false,`

### Step 2: Test Original Flow

```bash
# Build and run
cd /Users/yuda/github.com/kawai-network/veridium
wails3 dev
```

**In the app**:
1. Open chat
2. Type: "Hello"
3. Verify response appears
4. ✅ **PASS** if it works

### Step 3: Enable Backend

Edit `frontend/src/config/features.ts`:

```typescript
export const FEATURE_FLAGS = {
  USE_BACKEND_CHAT: true,   // ← Change to true
  DEBUG_MIGRATION: true,    // ← Enable logging
  // ... rest stay false
};
```

### Step 4: Rebuild and Test

```bash
# Rebuild (Wails will hot-reload)
# Or restart: Ctrl+C and wails3 dev again
```

**In the app**:
1. Open browser console (F12)
2. Type: "Hello"
3. Check console for:
   ```
   [Migration] sendMessage: Attempting BACKEND path
   [Migration] sendMessage: Backend SUCCESS
   ```
4. Wait a moment, refresh page if needed
5. ✅ **PASS** if message appears (even if slow/glitchy)

### Step 5: Disable and Verify Rollback

Edit `frontend/src/config/features.ts`:

```typescript
export const FEATURE_FLAGS = {
  USE_BACKEND_CHAT: false,  // ← Change back to false
  DEBUG_MIGRATION: false,
  // ...
};
```

Rebuild and test again - should work like before.

---

## ✅ You're Done with Phase 1 Testing If:

- [ ] Original works (flag OFF)
- [ ] Backend works (flag ON) - even if glitchy
- [ ] Rollback works (flag OFF again)
- [ ] No crashes or errors
- [ ] Issues documented

---

## 🐛 If Something Breaks:

### Backend Not Working?

1. Check console for errors
2. Verify Wails bindings loaded: `window.go` exists
3. Check backend is running
4. Try refresh page
5. **Fallback should work automatically**

### Can't Rollback?

1. Set `USE_BACKEND_CHAT: false`
2. Rebuild
3. Should work like original
4. If still broken: restore backup
   ```bash
   cp frontend/src/store/chat/slices/aiChat/actions/generateAIChat.ts.phase0 \
      frontend/src/store/chat/slices/aiChat/actions/generateAIChat.ts
   ```

---

## 📊 Expected Results

| Test | Flag OFF | Flag ON |
|------|----------|---------|
| Message sent | ✅ Works | ✅ Works (slower?) |
| Response appears | ✅ Immediately | ⚠️ After refresh? |
| Multiple messages | ✅ Works | ⚠️ Order issues? |
| Topic creation | ✅ Works | ⚠️ Doesn't switch? |
| Error handling | ✅ Works | ✅ Falls back |

**⚠️ = Expected issues for Phase 2**

---

## 🎯 Phase 1 Goal

**NOT** to make it perfect.  
**JUST** to prove:
1. Backend path works end-to-end
2. Fallback works if backend fails
3. Original still works (zero risk)
4. We can iterate safely

**Phase 2** will fix the issues.

---

## 📝 Report Your Findings

After testing, update `PHASE1_IMPLEMENTATION.md` with:

```markdown
### Test Results - [Your Name] - [Date]

**Flag OFF**:
- ✅ Works perfectly
- Issues: None

**Flag ON**:
- ✅ Backend called successfully
- ⚠️ Messages appear after refresh
- ⚠️ Topic doesn't switch automatically
- Issues: [list what you found]

**Overall**: READY / NOT READY for Phase 2
```

---

**Questions?** Check `PROPER_MIGRATION_STRATEGY.md` or `PHASE1_IMPLEMENTATION.md`


