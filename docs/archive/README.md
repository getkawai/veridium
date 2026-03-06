# Archived Documentation

**Location:** `veridium/docs/archive/`  
**Last Updated:** March 1, 2026  
**Total Documents:** 7

---

## 📋 Overview

This folder contains historical documentation that is no longer actively maintained but preserved for reference.

### Why Documents Are Archived

Documents are moved here when:
- ✅ Migration/upgrade completed
- ✅ Project phase finished
- ✅ Design decisions finalized
- ✅ Temporary working docs
- ✅ Superseded by newer documentation
- ✅ Historical test results/deployment logs

---

## 📚 Archived Documents

### Active References (Still Useful)

#### 1. Mainnet Deployment Summary
**File:** [`MAINNET_DEPLOYMENT_SUMMARY.md`](../MAINNET_DEPLOYMENT_SUMMARY.md)  
**Date:** January 23, 2026  
**Status:** ✅ **MOVED TO ACTIVE DOCS** - Production reference  
**Description:** Complete mainnet deployment summary with all 7 contract addresses

**Why Active:** Production deployment addresses still in use by backend

---

#### 2. Testing Guide
**File:** [`TESTING_GUIDE.md`](../TESTING_GUIDE.md)  
**Date:** March 1, 2026  
**Status:** ✅ **NEW ACTIVE DOC** - Testing procedures  
**Description:** Comprehensive testing procedures for all reward systems

**Why Active:** Still used for weekly settlement testing

---

### Historical Archive

#### 3. Cloudflare SDK Migration
**File:** `CLOUDFLARE_SDK_MIGRATION.md`  
**Date:** January 2026  
**Status:** ✅ Migration Completed  
**Description:** Migration guide from Cloudflare Go SDK v0.114.0 to v6.6.0

**Why Archived:** Migration completed, no longer needed for active development

---

#### 4. Deployment Fixes Summary
**File:** `DEPLOYMENT_FIXES_SUMMARY.md`  
**Date:** January 23, 2026  
**Status:** ✅ All Issues Resolved  
**Description:** Fixes for MINTER_ROLE grant failures and hardcoded testnet config

**Why Archived:** All fixes applied, deployment completed successfully

---

#### 5. macOS Launch Fix
**File:** `MACOS_LAUNCH_FIX.md`  
**Date:** January 2026  
**Status:** ✅ Fix Implemented  
**Description:** Root cause analysis and fix for macOS app launch issues

**Why Archived:** Fix implemented in production, no pending actions

---

#### 6. Solidity Compiler Upgrade
**Files:** 
- `SOLC_0.8.33_UPGRADE_RESULTS.md`
- `SOLC_UPGRADE_PLAN.md`  
**Date:** January 2026  
**Status:** ✅ Upgrade Completed  
**Description:** Plan and results for upgrading Solidity compiler to v0.8.33

**Why Archived:** Upgrade completed, historical reference only

---

#### 7. Testing Results (Historical)
**File:** `TESTING_RESULTS.md`  
**Date:** January 12, 2026  
**Status:** ⚠️ Superseded by `TESTING_GUIDE.md`  
**Description:** Detailed test results from January 2026 testing session

**Why Archived:** 
- Test results are date-specific and outdated
- Testing procedures extracted to active `TESTING_GUIDE.md`
- Historical reference only

---

## 🔍 How to Reference Archived Docs

### In Documentation

```markdown
<!-- Bad: Links to archived doc without context -->
See [Migration Guide](docs/archive/CLOUDFLARE_SDK_MIGRATION.md)

<!-- Good: Explains status and provides context -->
See [Migration Guide](docs/archive/CLOUDFLARE_SDK_MIGRATION.md) (archived, migration completed Jan 2026)

<!-- Better: Link to active doc instead -->
See [Testing Guide](TESTING_GUIDE.md) for current testing procedures
```

### In Code Comments

```go
// BAD: References archived doc without context
// See docs/archive/DEPLOYMENT_FIXES_SUMMARY.md

// GOOD: Explains what was fixed
// Fixed: MINTER_ROLE grant failures (Jan 2026)
// See: docs/archive/DEPLOYMENT_FIXES_SUMMARY.md for historical context
```

---

## 🔄 Restoration Process

If you need to restore an archived document (e.g., for similar migration):

1. **Copy** (don't move) back to active docs folder
2. **Update** status and date at the top
3. **Add note** about why it's being restored
4. **Archive again** when complete

**Example:**
```bash
# Restore for reference
cp docs/archive/CLOUDFLARE_SDK_MIGRATION.md docs/SDK_MIGRATION_REFERENCE.md

# Edit top of file
# **Status:** ⚠️ RESTORED - Referencing for v7.0.0 migration
# **Restored:** March 2026
```

---

## 🗑️ Cleanup Policy

Archived documents are kept indefinitely unless:

### Remove Immediately
- ❌ Contains sensitive information (API keys, private keys, etc.)
- ❌ Security vulnerabilities disclosed

### Review Annually
- ⚠️ Completely obsolete with no historical value
- ⚠️ Superseded by comprehensive new documentation

### Keep Indefinitely
- ✅ Historical deployment records
- ✅ Major migration/upgrade documentation
- ✅ Significant bug fixes and resolutions

---

## 📊 Archive Statistics

| Year | Documents Archived | Reason |
|------|-------------------|--------|
| 2026 Q1 | 7 | Migration completed, deployments finished |

**Total Archive Size:** ~50 KB (text only)

---

## 📝 Related Documentation

### Active Documentation
- [`README.md`](../README.md) - Main project documentation
- [`DEPLOYMENT.md`](../DEPLOYMENT.md) - Contract deployment guide
- [`TESTING_GUIDE.md`](../TESTING_GUIDE.md) - Testing procedures
- [`MAINNET_DEPLOYMENT_SUMMARY.md`](../MAINNET_DEPLOYMENT_SUMMARY.md) - Production deployment

### External Archives
- `x/` packages - Shared infrastructure moved to respective packages
- `kawai-contributor/` - Contributor-specific docs moved

---

**Maintained By:** Development Team  
**Archive Policy:** Last reviewed March 1, 2026
