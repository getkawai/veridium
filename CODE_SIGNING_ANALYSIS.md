# Code Signing Analysis untuk Veridium

## Status Implementasi Saat Ini

### ❌ Belum Diimplementasikan
Berdasarkan analisis terhadap workflow GitHub Actions dan Taskfile, **code signing belum diimplementasikan sama sekali**:

1. **macOS**: Hanya menggunakan ad-hoc signing (`codesign --sign -`) untuk development
2. **Windows**: Tidak ada signing sama sekali
3. **Linux**: Tidak ada signing untuk DEB/RPM packages

### Workflow Saat Ini
File `.github/workflows/release.yml` hanya melakukan:
- Build binary untuk macOS (ARM64 & AMD64)
- Build binary untuk Linux
- Build binary untuk Windows
- Upload ke GitHub Releases dan Cloudflare R2
- **Tidak ada proses signing atau notarization**

## Analisis Wails v3 Code Signing

### Fitur yang Tersedia

Wails v3 menyediakan built-in tools untuk:

#### 1. **macOS**
- ✅ Code signing dengan Developer ID
- ✅ Notarization otomatis
- ✅ Entitlements management
- ✅ Hardened runtime
- ✅ Keychain integration untuk credentials

#### 2. **Windows**
- ✅ Code signing dengan certificate (.pfx/.p12)
- ✅ Cross-platform signing (bisa dari macOS/Linux)
- ✅ Timestamp server support
- ✅ Support untuk EXE, MSI, MSIX

#### 3. **Linux**
- ✅ PGP signing untuk DEB packages
- ✅ PGP signing untuk RPM packages
- ✅ Cross-platform signing
- ✅ Key generation tools

### Tools yang Disediakan

```bash
wails3 setup signing              # Interactive wizard
wails3 setup entitlements         # macOS entitlements wizard
wails3 signing list               # List certificates
wails3 signing generate-key       # Generate PGP key
wails3 tool sign                  # Low-level signing
```

## Pros & Cons Implementasi

### ✅ PROS - Mengimplementasikan Code Signing

#### Keamanan & Trust
1. **User Trust**: Aplikasi signed tidak akan memunculkan warning "Unknown Developer"
2. **macOS Gatekeeper**: Aplikasi bisa langsung dijalankan tanpa bypass security
3. **Windows SmartScreen**: Mengurangi warning "Windows protected your PC"
4. **Malware Protection**: Membuktikan binary tidak dimodifikasi setelah release
5. **Enterprise Deployment**: Banyak perusahaan hanya mengizinkan signed apps

#### Distribution
1. **App Store Ready**: Requirement untuk Mac App Store dan Microsoft Store
2. **Auto-update**: Signed apps lebih mudah untuk implement auto-update
3. **Package Managers**: Linux repos (apt, yum) memerlukan signed packages
4. **Professional Image**: Menunjukkan profesionalisme dan legitimasi

#### Technical
1. **Built-in Tools**: Wails v3 sudah menyediakan semua tools yang dibutuhkan
2. **Cross-platform**: Windows & Linux bisa di-sign dari platform manapun
3. **CI/CD Ready**: Mudah diintegrasikan dengan GitHub Actions
4. **Secure Credential Storage**: Password disimpan di system keychain, bukan di code

### ❌ CONS - Mengimplementasikan Code Signing

#### Biaya
1. **Apple Developer**: $99/tahun untuk Developer ID certificate (MANDATORY)
2. **Windows Certificate**: $200-400/tahun (DigiCert, Sectigo, dll) - **Ada alternatif gratis, lihat section Windows Alternatives**
3. **Linux**: Gratis (PGP key), tapi perlu setup infrastructure

#### Kompleksitas
1. **Setup Time**: Initial setup memerlukan waktu (1-2 hari)
2. **Certificate Management**: Perlu renew certificates setiap tahun
3. **CI/CD Secrets**: Perlu manage secrets di GitHub Actions
4. **Multiple Platforms**: Perlu setup untuk 3 platform berbeda

#### Maintenance
1. **Certificate Expiry**: Harus monitor dan renew certificates
2. **Notarization Time**: macOS notarization menambah 1-2 menit per build
3. **Troubleshooting**: Signing errors bisa kompleks untuk di-debug
4. **Team Coordination**: Jika ada multiple developers, perlu share credentials

#### Limitations
1. **macOS Signing**: Harus dilakukan di macOS (tidak bisa cross-platform)
2. **Hardware Requirements**: macOS runner di GitHub Actions lebih mahal
3. **Build Time**: Signing + notarization menambah waktu build

### ⚠️ CONS - Tidak Mengimplementasikan Code Signing

#### User Experience
1. **macOS**: User harus klik kanan > Open untuk bypass Gatekeeper
2. **Windows**: User harus klik "More info" > "Run anyway" di SmartScreen
3. **Trust Issues**: User mungkin ragu untuk install aplikasi unsigned
4. **Support Burden**: Lebih banyak support tickets tentang "can't open app"

#### Distribution
1. **App Stores**: Tidak bisa distribute via Mac App Store atau Microsoft Store
2. **Enterprise**: Banyak perusahaan block unsigned apps
3. **Antivirus**: Beberapa antivirus lebih aggressive terhadap unsigned apps

#### Package Manager Limitations
1. **Homebrew Cask**: Sejak Homebrew 5.0.0 (2024), **unsigned casks akan dihapus dari official tap** (deadline: September 2026)
2. **No Workarounds**: Flag `--no-quarantine` sudah deprecated, tidak bisa bypass Gatekeeper lagi
3. **Mandatory Signing**: Homebrew sekarang enforce signing requirements yang sama dengan Apple
4. **User Experience**: Bahkan via Homebrew, unsigned apps tetap akan kena Gatekeeper warning

## ⚠️ Update: Homebrew & Package Manager Distribution

### Apakah Homebrew Bisa Bypass Code Signing Requirements?

**TL;DR: TIDAK** ❌

Berdasarkan [Homebrew 5.0.0 announcement](https://workbrew.com/blog/what-homebrew-5-0-0-means-for-your-mac-fleet) (released 2024):

#### Perubahan Kebijakan Homebrew

1. **Mandatory Signing (Deadline: September 2026)**
   - Casks tanpa code signing & notarization akan **dihapus dari official tap**
   - Tidak bisa lagi `brew install` untuk unsigned apps
   - Homebrew sekarang enforce security requirements yang sama dengan Apple

2. **Deprecation of `--no-quarantine`**
   - Flag untuk bypass Gatekeeper sudah deprecated
   - Homebrew "does not wish to easily provide circumvention to macOS security features"
   - User tidak bisa lagi bypass security warnings via Homebrew

3. **macOS 15 (Sequoia) Changes**
   - Control-click workaround untuk bypass unsigned apps **sudah dihapus**
   - Semakin sulit untuk user menjalankan unsigned apps

#### Implikasi untuk Veridium

**Scenario 1: Distribusi via Homebrew TANPA Signing**
```bash
# User experience:
brew install kawai  # ❌ GAGAL (setelah Sept 2026)
# Error: Cask 'kawai' requires code signing and notarization
```

**Scenario 2: Distribusi via Homebrew DENGAN Signing**
```bash
# User experience:
brew install kawai  # ✅ SUKSES
# App langsung bisa dibuka tanpa warning
```

#### Kesimpulan Package Manager

| Package Manager | Unsigned App Support | Notes |
|----------------|---------------------|-------|
| **Homebrew Cask** | ❌ Tidak (sejak 2026) | Mandatory signing & notarization |
| **MacPorts** | ⚠️ Terbatas | Masih allow, tapi user tetap kena Gatekeeper |
| **Direct Download** | ⚠️ Possible | User harus manual bypass (semakin sulit) |
| **App Store** | ❌ Tidak | Mandatory signing |

**Bottom Line**: Package managers **TIDAK menyelesaikan** masalah code signing. Bahkan Homebrew sekarang **memaksa** developers untuk sign & notarize apps mereka.

## Rekomendasi

### 🎯 Prioritas Implementasi

#### Phase 1: macOS (HIGHEST PRIORITY)
**Alasan**: macOS Gatekeeper paling strict, user experience paling terpengaruh

**Langkah**:
1. Beli Apple Developer account ($99/tahun)
2. Generate Developer ID Application certificate
3. Setup notarization credentials
4. Update `build/darwin/Taskfile.yml` dengan signing tasks
5. Update GitHub Actions workflow

**Estimasi**: 1-2 hari setup, 5-10 menit tambahan per build

#### Phase 2: Windows (MEDIUM PRIORITY)
**Alasan**: SmartScreen warning bisa di-bypass, tapi tetap mengganggu UX

**Langkah**:
1. Beli code signing certificate ($200-400/tahun)
2. Setup certificate di GitHub Actions secrets
3. Update `build/windows/Taskfile.yml` dengan signing tasks
4. Implement cross-platform signing dari Linux runner

**Estimasi**: 1 hari setup, 2-3 menit tambahan per build

#### Phase 3: Linux (LOW PRIORITY)
**Alasan**: Linux users lebih tech-savvy, signing tidak mandatory

**Langkah**:
1. Generate PGP key dengan `wails3 signing generate-key`
2. Setup PGP key di GitHub Actions secrets
3. Update `build/linux/Taskfile.yml` dengan signing tasks
4. Publish public key untuk user verification

**Estimasi**: 4 jam setup, 1-2 menit tambahan per build

### 💡 Implementasi Bertahap

#### Opsi A: Implement Semua Sekaligus
- **Pros**: Konsisten, professional dari awal
- **Cons**: Biaya upfront tinggi ($300-500), setup time 3-4 hari
- **Cocok untuk**: Production-ready release, enterprise customers

#### Opsi B: Implement Bertahap (RECOMMENDED)
- **Pros**: Spread cost, learn as you go, validate need
- **Cons**: Inconsistent UX di awal
- **Cocok untuk**: Beta/early access, iterative development

**Timeline Bertahap**:
```
Week 1: macOS signing + notarization
Week 2-3: Test dengan beta users
Week 4: Windows signing (jika feedback positif)
Week 5-6: Linux signing (optional)
```

#### Opsi C: Delay Sampai v1.0
- **Pros**: Save cost di early stage, focus on features
- **Cons**: Bad UX, harder to migrate later, **tidak bisa distribute via Homebrew setelah Sept 2026**
- **Cocok untuk**: MVP, internal testing only
- **⚠️ WARNING**: Jika target distribusi via Homebrew, signing **WAJIB** sebelum September 2026

## Windows Code Signing: Alternatif Gratis

### ⚠️ Reality Check: Tidak Ada Solusi Gratis yang Perfect

**TL;DR**: Self-signed certificates **TIDAK menghilangkan** SmartScreen warning. Bahkan certificate berbayar pun butuh "reputation building".

### Opsi yang Tersedia

#### 1. ❌ Self-Signed Certificate (TIDAK RECOMMENDED)
**Cara**:
```bash
# Generate self-signed certificate
New-SelfSignedCertificate -Type CodeSigningCert -Subject "CN=Kawai" -CertStoreLocation Cert:\CurrentUser\My
```

**Realita**:
- ❌ SmartScreen tetap muncul (bahkan lebih buruk)
- ❌ User harus install certificate secara manual
- ❌ Tidak ada trust dari Windows
- ❌ Lebih mencurigakan daripada unsigned

**Verdict**: Buang-buang waktu, lebih baik unsigned.

#### 2. ⚠️ OV Certificate ($200-300/tahun) - MASIH ADA WARNING
**Providers**: DigiCert, Sectigo, SSL.com

**Realita**:
- ⚠️ SmartScreen **TETAP MUNCUL** untuk app baru
- ⏱️ Butuh "reputation building" (ribuan downloads)
- 📊 Warning hilang setelah 3-6 bulan dengan volume tinggi
- 🔄 Warning muncul lagi setiap release versi baru

**Verdict**: Lebih baik dari unsigned, tapi bukan solusi instant.

#### 3. ✅ EV Certificate ($400-600/tahun) - INSTANT TRUST
**Providers**: DigiCert, Sectigo (butuh hardware token)

**Keuntungan**:
- ✅ **Instant trust**, no SmartScreen warning
- ✅ Tidak perlu reputation building
- ✅ Langsung trusted untuk semua versi baru

**Kekurangan**:
- 💰 Lebih mahal ($400-600/tahun)
- 🔑 Butuh hardware USB token (dikirim via pos)
- 📋 Verifikasi bisnis lebih ketat (butuh dokumen legal)
- 🚫 Tidak bisa untuk individual developer (butuh registered company)

**Verdict**: Best solution jika punya budget dan registered company.

#### 4. 🆓 Microsoft Store (GRATIS tapi terbatas)
**Cara**: Publish via Microsoft Store

**Keuntungan**:
- ✅ Gratis, Microsoft yang handle signing
- ✅ No SmartScreen warning
- ✅ Auto-update built-in
- ✅ Trusted distribution

**Kekurangan**:
- 📱 Hanya untuk MSIX packages (bukan EXE tradisional)
- 🔒 Harus follow Microsoft Store policies
- 💰 Microsoft ambil 15% revenue (jika paid app)
- 🐌 Review process bisa lambat
- 🚫 Tidak bisa distribute di luar Store

**Verdict**: Good untuk consumer apps, bad untuk enterprise/developer tools.

#### 5. 🎯 Hybrid Approach: Unsigned + User Education (RECOMMENDED untuk early stage)
**Strategi**:
1. Release unsigned binary
2. Provide clear documentation untuk bypass SmartScreen
3. Build reputation organically
4. Upgrade ke OV/EV certificate setelah traction

**Documentation untuk users**:
```markdown
## Installation on Windows

When you download Kawai, Windows SmartScreen may show a warning. This is normal for new applications.

**To install**:
1. Click "More info"
2. Click "Run anyway"

We're working on getting our application signed to remove this warning.
```

**Keuntungan**:
- ✅ Zero cost
- ✅ Fokus ke product development dulu
- ✅ Bisa upgrade nanti setelah validated
- ✅ Honest dengan users

**Kekurangan**:
- ⚠️ Conversion rate lebih rendah
- 📞 More support tickets
- 🎯 Tidak cocok untuk enterprise customers

### 📊 Comparison Table

| Opsi | Biaya/tahun | SmartScreen Warning | Setup Time | Best For |
|------|-------------|---------------------|------------|----------|
| Unsigned | $0 | ⚠️ Ya | 0 | MVP, beta testing |
| Self-signed | $0 | ⚠️⚠️ Ya (worse) | 1 jam | ❌ Jangan |
| OV Cert | $200-300 | ⚠️ Ya (3-6 bulan) | 1-2 hari | Growing apps |
| EV Cert | $400-600 | ✅ Tidak | 1 minggu | Production, enterprise |
| MS Store | $0 (15% cut) | ✅ Tidak | 1-2 minggu | Consumer apps |

### 💡 Rekomendasi Berdasarkan Stage

#### Stage 1: MVP/Beta (0-1000 users)
**Pilihan**: Unsigned + documentation
- **Biaya**: $0
- **Reasoning**: Focus on product-market fit, bukan certificate
- **Action**: Buat clear installation guide

#### Stage 2: Early Traction (1000-10000 users)
**Pilihan**: OV Certificate
- **Biaya**: $200-300/tahun
- **Reasoning**: Mulai build reputation, show professionalism
- **Action**: Beli OV cert, submit app ke Microsoft untuk review

#### Stage 3: Scale (10000+ users atau enterprise customers)
**Pilihan**: EV Certificate
- **Biaya**: $400-600/tahun
- **Reasoning**: Instant trust, no friction, professional image
- **Action**: Upgrade ke EV, register company jika belum

### 🎯 Specific Recommendation untuk Kawai

Berdasarkan analisis codebase (blockchain platform dengan mining rewards):

**Phase 1 (Now)**: 
- ✅ macOS: Beli Apple Developer ($99) - **MANDATORY**
- ⏸️ Windows: Unsigned dengan clear docs - **SAVE $200-400**
- ⏸️ Linux: Unsigned - **SAVE TIME**

**Phase 2 (After 1000+ active users)**:
- ✅ Windows: OV Certificate ($200-300)
- ✅ Linux: PGP signing (gratis)

**Phase 3 (Enterprise ready)**:
- ✅ Windows: Upgrade ke EV Certificate ($400-600)

**Total Savings Phase 1**: $200-400/tahun (Windows cert)

### 🔗 Dimana Beli Code Signing Certificate?

Code signing certificate **BERBEDA** dengan SSL/TLS certificate (domain certificate). Berikut perbedaannya:

#### Perbedaan Code Signing vs SSL Certificate

| Aspek | SSL/TLS Certificate | Code Signing Certificate |
|-------|---------------------|--------------------------|
| **Fungsi** | Encrypt komunikasi website (HTTPS) | Sign software/executable files |
| **Digunakan untuk** | Website, API, web server | .exe, .msi, .app, .dmg, drivers |
| **Validasi** | Domain ownership | Identity verification (individual/company) |
| **Lokasi** | Installed di web server | Embedded dalam software |
| **Harga** | $10-100/tahun | $200-600/tahun |
| **Verification** | Domain validation (instant) | Identity validation (1-7 hari) |
| **Contoh** | Let's Encrypt, Cloudflare | DigiCert, Sectigo |

**TL;DR**: SSL untuk website, Code Signing untuk aplikasi desktop/mobile.

#### Certificate Providers & Harga (2025)

##### 1. 💰 Budget Option: SSL.com
- **OV**: $199/tahun (⚠️ masih ada SmartScreen warning)
- **EV**: $399/tahun (✅ instant trust)
- **Support**: Email only
- **Verification**: 1-3 hari (OV), 3-7 hari (EV)
- **Hardware Token**: Included (dikirim via pos)
- **Link**: [ssl.com/certificates/code-signing](https://www.ssl.com/certificates/code-signing/)

**Pros**: Paling murah, good untuk individual developer
**Cons**: Support terbatas, verification bisa lambat

##### 2. 🎯 Mid-tier: Sectigo (RECOMMENDED)
- **OV**: $211-299/tahun
- **EV**: $499/tahun
- **Support**: Email + phone
- **Verification**: 1-2 hari (OV), 3-5 hari (EV)
- **Hardware Token**: Included
- **Link**: [sectigo.com/ssl-certificates-tls/code-signing](https://sectigo.com/ssl-certificates-tls/code-signing)

**Pros**: Balance antara harga dan service, trusted CA
**Cons**: Tidak se-premium DigiCert

**Reseller Murah**:
- [codesigningstore.com](https://codesigningstore.com/code-signing/sectigo-code-signing-certificate) - $211/tahun
- [cheapsslsecurity.com](https://cheapsslsecurity.com) - $220/tahun

##### 3. 👑 Premium: DigiCert
- **OV**: $399/tahun
- **EV**: $599/tahun
- **Support**: 24/7 phone + priority
- **Verification**: 1 hari (OV), 2-3 hari (EV)
- **Hardware Token**: Included (YubiKey)
- **Link**: [digicert.com/signing/code-signing-certificates](https://www.digicert.com/signing/code-signing-certificates)

**Pros**: Best reputation, fastest support, enterprise-grade
**Cons**: Paling mahal

##### 4. 🆓 Alternative: Microsoft Store
- **Harga**: Gratis (Microsoft yang sign)
- **Support**: Microsoft Store support
- **Verification**: App review process
- **Link**: [partner.microsoft.com](https://partner.microsoft.com/en-us/dashboard/home)

**Pros**: Gratis, trusted distribution
**Cons**: Hanya untuk MSIX, tidak bisa distribute di luar Store

#### 📋 Proses Pembelian & Verification

##### OV (Organization Validation) Certificate

**Requirements**:
- ✅ Valid email address
- ✅ Phone number yang bisa dihubungi
- ✅ Business registration (jika company)
- ✅ ATAU Government-issued ID (jika individual)

**Process**:
1. **Order** (5 menit): Pilih certificate, isi form
2. **Payment** (instant): Credit card, PayPal, wire transfer
3. **Verification** (1-3 hari):
   - CA akan email/call untuk verify identity
   - Jika company: verify business registration
   - Jika individual: verify ID (KTP, passport)
4. **Issuance** (instant setelah approved):
   - Download certificate (.pfx file)
   - ATAU hardware token dikirim via pos (3-7 hari)

**Timeline**: 1-3 hari verification + 3-7 hari shipping (jika hardware token)

##### EV (Extended Validation) Certificate

**Requirements** (lebih ketat):
- ✅ Registered company (TIDAK bisa individual)
- ✅ Business registration documents
- ✅ D-U-N-S number (optional tapi recommended)
- ✅ Verified business address
- ✅ Authorized representative

**Process**:
1. **Order** (10 menit): Isi detailed form
2. **Payment** (instant)
3. **Verification** (3-7 hari):
   - CA akan verify company existence
   - Check business registration dengan government database
   - Call company phone number untuk verify
   - Verify authorized representative
4. **Issuance**:
   - Hardware token (USB) dikirim via courier
   - HARUS menggunakan hardware token (security requirement)

**Timeline**: 3-7 hari verification + 3-7 hari shipping

#### ⚠️ Important Notes (2025 Updates)

1. **Hardware Token Mandatory** (sejak 2023):
   - Semua code signing certificate HARUS disimpan di hardware token
   - Tidak bisa lagi download .pfx file saja
   - Token included dalam harga (USB key atau YubiKey)

2. **Certificate Lifespan** (sejak Dec 2025):
   - Maximum 1 tahun (366 hari)
   - Tidak bisa lagi beli 2-3 tahun
   - Harus renew setiap tahun

3. **Verification Process**:
   - OV: 1-3 hari (bisa lebih cepat jika dokumen lengkap)
   - EV: 3-7 hari (lebih ketat)
   - Jika dokumen tidak lengkap, bisa 1-2 minggu

#### 💡 Rekomendasi Praktis

##### Untuk Individual Developer (Kawai - Early Stage)
**Pilihan**: Sectigo OV via reseller
- **Provider**: [codesigningstore.com](https://codesigningstore.com/code-signing/sectigo-code-signing-certificate)
- **Harga**: $211/tahun
- **Timeline**: Order hari ini, approved 2-3 hari, token sampai 1 minggu
- **Requirements**: KTP/Passport + email + phone

**Steps**:
1. Order di codesigningstore.com
2. Upload KTP/Passport scan
3. Tunggu verification call/email (1-2 hari)
4. Approved → token dikirim (3-5 hari)
5. Total: ~1 minggu dari order sampai bisa sign

##### Untuk Registered Company (Kawai - Scale)
**Pilihan**: Sectigo EV atau DigiCert EV
- **Provider**: Sectigo EV ($499) atau DigiCert EV ($599)
- **Timeline**: 1-2 minggu (verification + shipping)
- **Requirements**: Business registration, company docs

#### 🛒 Step-by-Step: Beli Sectigo OV (Recommended)

1. **Visit**: [codesigningstore.com/code-signing/sectigo-code-signing-certificate](https://codesigningstore.com/code-signing/sectigo-code-signing-certificate)

2. **Select**:
   - Certificate Type: "Sectigo Code Signing Certificate"
   - Validity: 1 year ($211)
   - Add to cart

3. **Checkout**:
   - Fill personal/company info
   - Email: your@email.com
   - Phone: +62-xxx-xxx-xxxx

4. **Upload Documents**:
   - Individual: KTP atau Passport (scan/photo)
   - Company: Business registration + authorized rep ID

5. **Verification** (1-3 hari):
   - CA akan email untuk confirm
   - Mungkin ada phone call untuk verify
   - Respond cepat untuk speed up process

6. **Receive Token** (3-7 hari):
   - Hardware USB token dikirim via courier
   - Track shipping via email

7. **Setup**:
   - Plug in USB token
   - Install certificate
   - Configure Wails Taskfile
   - Sign your app!

#### ❓ FAQ

**Q: Apakah bisa pakai SSL certificate untuk code signing?**
A: ❌ TIDAK. SSL certificate hanya untuk website (HTTPS). Code signing butuh certificate khusus.

**Q: Apakah bisa beli sekali untuk semua platform?**
A: ❌ TIDAK. Windows code signing certificate hanya untuk Windows. macOS butuh Apple Developer account terpisah.

**Q: Berapa lama certificate berlaku?**
A: Maximum 1 tahun (366 hari) sejak Dec 2025. Harus renew setiap tahun.

**Q: Apakah bisa refund jika tidak cocok?**
A: Tergantung provider. Biasanya ada 30-day money-back guarantee sebelum certificate issued.

**Q: Apakah bisa transfer certificate ke developer lain?**
A: ❌ TIDAK. Certificate tied ke identity yang verified. Tidak bisa transfer.

**Q: Apakah hardware token bisa dipakai untuk multiple projects?**
A: ✅ YA. Satu certificate bisa sign unlimited apps/versions.

## Implementasi Detail

### 1. macOS Signing Setup

#### A. Persiapan
```bash
# 1. Beli Apple Developer account
# 2. Generate Developer ID certificate di developer.apple.com
# 3. Download dan install certificate ke Keychain

# 4. Verify certificate
wails3 signing list
```

#### B. Update Taskfile
Edit `build/darwin/Taskfile.yml`, tambahkan:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Kawai (TEAMID)"
  KEYCHAIN_PROFILE: "kawai-notarize"
  ENTITLEMENTS: "build/darwin/entitlements.plist"

tasks:
  sign:
    summary: Build, package, and sign
    deps:
      - task: build
        vars:
          PRODUCTION: "true"
    cmds:
      - task: create:app:bundle
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}.app --identity "{{.SIGN_IDENTITY}}" --entitlements "{{.ENTITLEMENTS}}" --hardened-runtime

  sign:notarize:
    summary: Build, package, sign, and notarize
    deps:
      - task: sign
    cmds:
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}.app --identity "{{.SIGN_IDENTITY}}" --entitlements "{{.ENTITLEMENTS}}" --hardened-runtime --notarize --keychain-profile "{{.KEYCHAIN_PROFILE}}"
```

#### C. Generate Entitlements
```bash
wails3 setup entitlements
# Pilih "Both" untuk generate dev dan production entitlements
```

#### D. Setup Notarization
```bash
# Generate app-specific password di appleid.apple.com
wails3 signing credentials \
  --apple-id "your@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password" \
  --profile "kawai-notarize"
```

#### E. Update GitHub Actions
```yaml
build-macos:
  steps:
    # ... existing steps ...
    
    - name: Import Certificate
      env:
        CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
        CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
      run: |
        echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
        security create-keychain -p "" build.keychain
        security default-keychain -s build.keychain
        security unlock-keychain -p "" build.keychain
        security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
        security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

    - name: Store Notarization Credentials
      env:
        APPLE_ID: ${{ secrets.APPLE_ID }}
        APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
        APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
      run: |
        wails3 signing credentials \
          --apple-id "$APPLE_ID" \
          --team-id "$APPLE_TEAM_ID" \
          --password "$APPLE_APP_PASSWORD" \
          --profile "kawai-notarize"

    - name: Build, Sign, and Notarize ARM64
      run: |
        PRODUCTION=true ARCH=arm64 wails3 task darwin:build
        wails3 task darwin:sign:notarize
        mv bin/Kawai.app bin/Kawai-arm64.app
```

### 2. Windows Signing Setup

#### A. Update Taskfile
Edit `build/windows/Taskfile.yml`, tambahkan:

```yaml
vars:
  SIGN_CERTIFICATE: "certificate.pfx"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"

tasks:
  sign:
    summary: Build and sign executable
    deps:
      - task: build
        vars:
          PRODUCTION: "true"
    cmds:
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}.exe --certificate "{{.SIGN_CERTIFICATE}}" --password "$WAILS_WINDOWS_CERT_PASSWORD" --timestamp "{{.TIMESTAMP_SERVER}}"

  sign:installer:
    summary: Build and sign NSIS installer
    deps:
      - task: create:nsis:installer
    cmds:
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}-installer.exe --certificate "{{.SIGN_CERTIFICATE}}" --password "$WAILS_WINDOWS_CERT_PASSWORD" --timestamp "{{.TIMESTAMP_SERVER}}"
```

#### B. Update GitHub Actions
```yaml
build-windows:
  steps:
    # ... existing steps ...
    
    - name: Import Certificate
      env:
        CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
      run: |
        $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
        [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

    - name: Build and Sign
      env:
        WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
      run: |
        wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx
```

### 3. Linux Signing Setup

#### A. Generate PGP Key
```bash
wails3 signing generate-key \
  --name "Kawai" \
  --email "security@getkawai.com" \
  --comment "Package Signing Key" \
  --bits 4096 \
  --expiry "2y" \
  --output-private signing-key.asc \
  --output-public signing-key.pub.asc
```

#### B. Update Taskfile
Edit `build/linux/Taskfile.yml`, tambahkan:

```yaml
vars:
  PGP_KEY: "signing-key.asc"
  SIGN_ROLE: "builder"

tasks:
  sign:deb:
    summary: Build and sign DEB package
    deps:
      - task: create:deb
    cmds:
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}.deb --pgp-key "{{.PGP_KEY}}" --pgp-password "$WAILS_PGP_PASSWORD" --role "{{.SIGN_ROLE}}"

  sign:rpm:
    summary: Build and sign RPM package
    deps:
      - task: create:rpm
    cmds:
      - wails3 tool sign --input {{.BIN_DIR}}/{{.APP_NAME}}.rpm --pgp-key "{{.PGP_KEY}}" --pgp-password "$WAILS_PGP_PASSWORD"

  sign:packages:
    summary: Build and sign all packages
    deps:
      - task: sign:deb
      - task: sign:rpm
```

#### C. Update GitHub Actions
```yaml
build-linux:
  steps:
    # ... existing steps ...
    
    - name: Import PGP Key
      env:
        PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
      run: |
        echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

    - name: Build and Sign Packages
      env:
        WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
      run: |
        wails3 task linux:sign:packages PGP_KEY=signing-key.asc
```

## Biaya Total

### One-time Setup
- Developer time: 3-4 hari × $500/hari = **$1,500-2,000**

### Annual Recurring

#### Minimum (Phase 1 - MVP/Beta)
- Apple Developer: **$99/tahun** (mandatory)
- Windows Certificate: **$0** (unsigned dengan docs)
- Linux: **$0** (unsigned atau PGP gratis)
- **Total: $99/tahun**

#### Recommended (Phase 2 - Growing)
- Apple Developer: **$99/tahun**
- Windows OV Certificate: **$200-300/tahun**
- Linux PGP: **$0**
- **Total: $300-400/tahun**

#### Enterprise (Phase 3 - Scale)
- Apple Developer: **$99/tahun**
- Windows EV Certificate: **$400-600/tahun**
- Linux PGP: **$0**
- **Total: $500-700/tahun**

### CI/CD Costs
- macOS runner: ~$0.08/menit
- Signing + notarization: ~10 menit/build
- Estimasi: **$0.80/build** (hanya macOS yang mahal)

## Kesimpulan

### ✅ RECOMMENDED: Phased Implementation

**Alasan**:
1. **User Experience**: Drastically better, especially on macOS
2. **Professional Image**: Essential untuk production app
3. **Security**: Protects users dari tampered binaries
4. **Future-proof**: Required untuk app stores dan enterprise
5. **Wails v3 Support**: Tools sudah tersedia, implementation straightforward
6. **Cost-Effective**: Start dengan $99/tahun, scale sesuai kebutuhan
7. **No Workarounds**: Package managers tidak bisa bypass signing requirements lagi

**Timeline**:
- **Week 1**: macOS signing (highest impact)
- **Week 2-3**: Beta testing
- **Week 4**: Windows signing
- **Week 5**: Linux signing (optional)

**Budget**:
- Initial: $1,500-2,000 (developer time)
- Annual: $300-500 (certificates)
- Per-build: $0.80 (CI/CD)

### 🎯 Action Items

#### Phase 1: MVP/Beta (Week 1-2) - PRIORITY
1. **macOS Signing** (Week 1):
   - [ ] Beli Apple Developer account ($99)
   - [ ] Generate Developer ID certificate
   - [ ] Setup notarization credentials
   - [ ] Implement macOS signing di Taskfile
   - [ ] Update GitHub Actions workflow
   - [ ] Test signed build

2. **Windows Documentation** (Week 1):
   - [ ] Buat clear installation guide untuk bypass SmartScreen
   - [ ] Add FAQ section tentang security warning
   - [ ] Update README dengan installation instructions
   - [ ] **SKIP buying certificate untuk sekarang** (save $200-400)

3. **Linux** (Week 1):
   - [ ] **SKIP signing untuk sekarang** (unsigned OK untuk beta)
   - [ ] Focus on functionality testing

#### Phase 2: Growing (After 1000+ users) - OPTIONAL
1. **Windows Signing** (Month 2-3):
   - [ ] Evaluate user feedback tentang SmartScreen friction
   - [ ] Jika conversion rate terpengaruh, beli OV certificate ($200-300)
   - [ ] Implement Windows signing
   - [ ] Submit app ke Microsoft untuk reputation building

2. **Linux Signing** (Month 2-3):
   - [ ] Generate PGP key (gratis)
   - [ ] Implement DEB/RPM signing
   - [ ] Publish public key untuk verification

#### Phase 3: Enterprise (After validated) - OPTIONAL
1. **Windows EV Upgrade** (Month 6+):
   - [ ] Jika targeting enterprise customers, upgrade ke EV ($400-600)
   - [ ] Register company jika belum (required untuk EV)
   - [ ] Order hardware token
   - [ ] Update signing workflow

2. **Maintenance**:
   - [ ] Setup auto-renewal reminders untuk certificates
   - [ ] Document signing process untuk team
   - [ ] Monitor certificate expiry dates

### 📚 Resources

- [Wails v3 Signing Docs](https://wails.io/docs/guides/build/signing)
- [Apple Code Signing Guide](https://developer.apple.com/support/code-signing/)
- [Apple Notarization](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft Code Signing](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Homebrew 5.0.0 Changes](https://workbrew.com/blog/what-homebrew-5-0-0-means-for-your-mac-fleet) - Mandatory signing requirements

---

## FAQ: Package Manager Distribution

### Q: Apakah distribusi via Homebrew bisa bypass code signing?
**A: TIDAK.** Sejak Homebrew 5.0.0, unsigned casks akan dihapus dari official tap (deadline: September 2026). Homebrew sekarang enforce signing requirements yang sama dengan Apple.

### Q: Bagaimana dengan MacPorts atau package manager lain?
**A: Tetap kena Gatekeeper.** Package manager lain mungkin masih allow unsigned apps, tapi user tetap akan mendapat Gatekeeper warning saat membuka app. macOS 15+ semakin mempersulit bypass.

### Q: Apakah bisa distribute via custom Homebrew tap tanpa signing?
**A: Secara teknis bisa,** tapi:
- User tetap kena Gatekeeper warning
- Bad user experience
- Tidak recommended oleh Homebrew
- Custom tap kurang discoverable

### Q: Apakah ada cara lain untuk distribute tanpa signing?
**A: Direct download masih possible,** tapi:
- User harus manual bypass Gatekeeper (semakin sulit di macOS 15+)
- Bad UX, banyak support tickets
- Tidak professional
- Tidak recommended untuk production app

### Q: Kapan deadline untuk implement signing jika mau distribute via Homebrew?
**A: September 2026** - Setelah itu, unsigned casks akan dihapus dari official Homebrew tap.
