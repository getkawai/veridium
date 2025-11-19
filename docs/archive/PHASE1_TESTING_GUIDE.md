# Phase 1 Testing Guide

> **🎯 Goal**: Verify both frontend (original) and backend (new) paths work correctly

## 📋 Quick Testing Checklist

### ✅ Step 1: Test Original Frontend (Flag OFF)

**Current Status**: `USE_BACKEND_CHAT: false`

1. **Restart application**
   ```bash
   # Stop current app (if running)
   # Rebuild and run:
   cd /Users/yuda/github.com/kawai-network/veridium
   wails dev  # or your usual dev command
   ```

2. **Open browser console** (F12 or Cmd+Option+I)

3. **Create new chat**:
   - Click "New Chat" or equivalent
   - Type: "Hello, how are you?"
   - Press Enter

4. **Verify**:
   - [ ] AI responds normally
   - [ ] Topic is created automatically
   - [ ] Topic has a generated title
   - [ ] NO "[Migration]" logs in console (unless DEBUG_MIGRATION is on)
   - [ ] Messages appear immediately

5. **Send follow-up**:
   - Type: "Tell me more"
   - Verify it goes to same topic

6. **Check persistence**:
   - Refresh page
   - Navigate to topic
   - Verify messages are still there

**✅ If all above works**: Original frontend is intact! Proceed to Step 2.

**❌ If something breaks**: 
- Check console for errors
- Document in `PHASE1_TESTING_RESULTS.md`
- **DO NOT** enable backend yet
- Report issue for fixing

---

### ✅ Step 2: Test Backend Path (Flag ON)

**Change flag**: Edit `frontend/src/config/features.ts`

```typescript
USE_BACKEND_CHAT: true, // Enable backend path
DEBUG_MIGRATION: true,  // Enable verbose logging
```

1. **Restart application**

2. **Open browser console** (should see more logs now)

3. **Create new chat**:
   - Click "New Chat"
   - Type: "Why do people sleep?"
   - Press Enter

4. **Watch console logs** - You should see:
   ```
   [Migration] sendMessage: Attempting BACKEND path { activeId: '...', ... }
   [Migration] sendMessage: Backend SUCCESS { messageId: '...', ... }
   [Migration] sendMessage: Backend created topic { topicId: '...' }
   [Migration] sendMessage: Backend path completed successfully
   ```

5. **Verify in UI**:
   - [ ] User message appears
   - [ ] AI response appears (even if delayed)
   - [ ] Topic is created in sidebar
   - [ ] No crash or blank screen

6. **Check database**:
   - Close app
   - Reopen app
   - Messages should persist

**✅ If backend logs appear and chat works**: Backend integration successful!

**⚠️ Known Issues** (these are EXPECTED in Phase 1):
- UI might not update immediately
- Topic might not switch automatically
- Message might appear twice (temp + real)
- These will be fixed in Phase 2

---

### ✅ Step 3: Test Fallback Mechanism

**Keep flag ON**: `USE_BACKEND_CHAT: true`

**Simulate backend failure** (choose one method):

**Option A: Modify code temporarily**

Edit `frontend/src/services/backend/backendAgentChat.ts`:

```typescript
export async function sendMessage(params: BackendSendMessageParams) {
  // TEMPORARY: Force error for testing
  throw new Error('Simulated backend failure');
  
  // Original code below...
  const response = await Chat(params);
  // ...
}
```

**Option B: Stop backend service**
- Stop Wails backend (Ctrl+C if running in terminal)
- Keep frontend running

**Test**:
1. Send message
2. Check console - should see:
   ```
   [Migration] Backend chat failed, falling back to frontend: Error: ...
   [Migration] sendMessage: Backend FAILED, falling back to frontend
   ```
3. **Verify**: Chat still works (via frontend fallback)

**✅ If chat works despite backend failure**: Fallback mechanism works!

**After testing**: Revert the code changes (if using Option A)

---

## 🎨 Visual Testing

### What to Look For:

1. **Original Frontend (Flag OFF)**:
   ```
   User types message
      ↓
   Message appears instantly (optimistic UI)
      ↓
   Loading indicator
      ↓
   AI response streams in
      ↓
   Topic auto-created with title
   ```

2. **Backend Path (Flag ON)**:
   ```
   User types message
      ↓
   [May have slight delay] ← EXPECTED IN PHASE 1
      ↓
   Message appears
      ↓
   AI response appears (full, no streaming yet)
      ↓
   Topic created (might need manual switch) ← EXPECTED
   ```

---

## 📊 Expected Results Summary

| Scenario | Flag | Expected Result |
|----------|------|-----------------|
| Normal chat | OFF | Works perfectly (original) |
| Normal chat | ON | Works, minor UI delays OK |
| Backend fails | ON | Falls back to frontend |
| Topic creation | OFF | Auto-creates with title |
| Topic creation | ON | Creates, might need manual switch |
| Persistence | Both | Messages saved to DB |

---

## 🐛 When to Report Issues

**Report if you see**:
- ❌ Crash or blank screen
- ❌ No response from AI (with either flag)
- ❌ Console errors (red, not warnings)
- ❌ Messages not persisting to database
- ❌ Backend doesn't get called at all (with flag ON)

**Don't report** (these are expected):
- ⚠️ Slight delay in UI updates (Flag ON)
- ⚠️ Topic not switching automatically (Flag ON)
- ⚠️ No streaming (Flag ON) - Phase 5 feature
- ⚠️ Duplicate messages (Flag ON) - Fixed in Phase 2

---

## 💡 Quick Debugging

### Can't see migration logs?

Check `features.ts`:
```typescript
DEBUG_MIGRATION: true,  // Must be true
```

### Backend not being called?

Check console:
```typescript
// Run in browser console:
import { isBackendAvailable } from '@/config/features';
console.log('Backend available:', isBackendAvailable());

// Should return: true (in Wails app) or false (in web browser)
```

### Want to see all feature flags?

```typescript
import { getFeatureFlagsStatus } from '@/config/features';
console.log(getFeatureFlagsStatus());
```

---

## 📝 Recording Results

As you test, update `PHASE1_TESTING_RESULTS.md`:

```markdown
### Test 1: Flag OFF ✅
- [x] Chat works
- [x] Topics created
- [x] No issues

Console logs:
[paste relevant logs]

### Test 2: Flag ON ✅
- [x] Backend called
- [x] Response received
- [x] Minor UI delay (expected)

Issues:
- Topic not switching (expected, Phase 2)

### Test 3: Fallback ✅
- [x] Falls back to frontend
- [x] No data loss
```

---

## 🚀 After All Tests Pass

### Success Criteria Met:
- [x] Original frontend works (Flag OFF)
- [x] Backend integration works (Flag ON)
- [x] Fallback mechanism works
- [x] No critical bugs
- [x] Database persistence works

### Decision:
**Keep flag OFF** for now (`USE_BACKEND_CHAT: false`) until Phase 2 is ready.

### Next Steps:
1. Document all findings in `PHASE1_TESTING_RESULTS.md`
2. Plan Phase 2: State Synchronization
3. Fix UI update issues
4. Add optimistic UI to backend path

---

## ⏱️ Estimated Time

- **Step 1** (Flag OFF): 10 minutes
- **Step 2** (Flag ON): 10 minutes
- **Step 3** (Fallback): 10 minutes
- **Documentation**: 10 minutes

**Total**: ~40 minutes

---

## 🎯 Your Turn!

**Start with**:
1. Restart app (flag is already OFF)
2. Test basic chat
3. Verify original frontend works
4. Document results

**Then report back with**:
- ✅ What worked
- ❌ What didn't work
- 📊 Console logs (screenshot or text)

Let me know when you're ready to proceed! 🚀

