# Register Admin Command

Command untuk mendaftarkan admin addresses ke contributor list agar bisa menerima reward (30% admin fee).

## Problem

Ketika `RecordJobReward` dipanggil, sistem akan mendistribusikan reward dengan rasio 70/30:
- 70% untuk contributor
- 30% untuk admin

Jika admin address tidak terdapat di contributor list, maka proses akan gagal dengan error:
```
failed to update admin fee: failed to get account {adminAddress}: failed to get from KV: ...
```

## Solution

Command ini akan:
1. Mengambil semua treasury addresses dari `internal/constant/treasury.go`
2. Mendaftarkan/update setiap address sebagai admin contributor
3. Memastikan admin bisa menerima reward tanpa error

## Usage

### Register semua treasury addresses

```bash
go run cmd/register-admin/main.go
```

### Dry run (preview tanpa melakukan perubahan)

```bash
go run cmd/register-admin/main.go --dry-run
```

### Register specific address

```bash
go run cmd/register-admin/main.go --address 0x1234567890123456789012345678901234567890
```

## Output Example

```
✓ Connected to Cloudflare KV
Admin addresses to register count=4

[1/4] Processing: 0x1234567890123456789012345678901234567890
  ✅ Successfully registered/updated as admin

[2/4] Processing: 0x2345678901234567890123456789012345678901
  ⏭️  Already registered as admin

[3/4] Processing: 0x3456789012345678901234567890123456789012
  ⚠️  Exists as regular contributor, will update to admin
  ✅ Successfully registered/updated as admin

[4/4] Processing: 0x4567890123456789012345678901234567890123
  ✅ Successfully registered/updated as admin

==================================================
Registration Summary:
  Total addresses: 4
  ✅ Success: 3
  ⏭️  Skipped (already admin): 1
  ❌ Failed: 0
==================================================
```

## Features

- ✅ Auto-creates admin accounts if they don't exist
- ✅ Updates existing contributors to admin status
- ✅ Preserves existing balances when updating
- ✅ Validates Ethereum addresses
- ✅ Dry-run mode for safe testing
- ✅ Detailed progress reporting
- ✅ Summary statistics

## Admin Account Properties

Admin accounts yang dibuat memiliki properties:
- `IsAdmin`: true
- `Status`: "admin"
- `HardwareSpecs`: "Admin Account"
- `EndpointURL`: "" (tidak perlu endpoint)
- `AccumulatedRewards`: "0" (atau preserved jika sudah ada)
- `AccumulatedUSDT`: "0" (atau preserved jika sudah ada)
- `IsActive`: true

## Integration

Fungsi `RecordJobReward` di `pkg/store/contributor.go` sudah diupdate untuk:
1. Otomatis memanggil `EnsureAdminExists()` sebelum update admin balance
2. Membuat admin account jika belum ada
3. Mencegah error saat admin address tidak terdaftar

## When to Run

Run command ini:
- ✅ Setelah deploy pertama kali
- ✅ Ketika menambah treasury address baru
- ✅ Ketika ada error "failed to update admin fee"
- ✅ Sebelum menjalankan settlement/dividend process

