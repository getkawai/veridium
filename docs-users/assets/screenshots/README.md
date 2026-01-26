# Screenshots for Documentation

This directory contains screenshots used in the user documentation.

## Required Screenshots

### macOS Installation

#### 1. `macos-right-click-open.png`
**Description**: Screenshot showing the right-click context menu on Kawai.app in Applications folder

**What to capture**:
- Finder window with Applications folder
- Kawai.app icon visible
- Right-click context menu open
- "Open" option highlighted

**Recommended size**: 800x600px or similar
**Format**: PNG with transparency if possible

**How to capture**:
1. Open Applications folder in Finder
2. Right-click on Kawai.app
3. Take screenshot (Cmd+Shift+4, then Space to capture window)
4. Crop to show relevant area

---

#### 2. `macos-security-settings.png`
**Description**: Screenshot showing the Security & Privacy settings with "Open Anyway" button

**What to capture**:
- System Settings/Preferences window
- Privacy & Security (or Security & Privacy) section
- Message about Kawai being blocked
- "Open Anyway" button visible

**Recommended size**: 1000x700px or similar
**Format**: PNG

**How to capture**:
1. Try to open Kawai.app (it will be blocked)
2. Open System Settings → Privacy & Security
3. Scroll to Security section
4. Take screenshot showing the blocked app message
5. Highlight the "Open Anyway" button

**Variations needed**:
- macOS 13+ (Ventura/Sonoma): `macos-security-settings-ventura.png`
- macOS 12 (Monterey): `macos-security-settings-monterey.png`
- macOS 11 (Big Sur): `macos-security-settings-bigsur.png`

---

### Windows Installation

#### 3. `windows-smartscreen-warning.png`
**Description**: Screenshot showing Windows SmartScreen warning dialog

**What to capture**:
- Full SmartScreen dialog
- "Windows protected your PC" message
- "More info" link visible
- App name and publisher info

**Recommended size**: 600x400px or similar
**Format**: PNG

**How to capture**:
1. Download unsigned Kawai.exe
2. Try to run it
3. SmartScreen warning appears
4. Take screenshot (Win+Shift+S)
5. Capture the entire dialog

---

#### 4. `windows-smartscreen-run-anyway.png`
**Description**: Screenshot showing SmartScreen after clicking "More info"

**What to capture**:
- SmartScreen dialog after clicking "More info"
- "Run anyway" button visible
- App details shown

**Recommended size**: 600x400px or similar
**Format**: PNG

**How to capture**:
1. From previous step, click "More info"
2. Take screenshot showing "Run anyway" button
3. Highlight the button

---

#### 5. `windows-uac-prompt.png`
**Description**: Screenshot showing User Account Control (UAC) prompt

**What to capture**:
- UAC dialog asking for permission
- "Yes" and "No" buttons
- App name and publisher info

**Recommended size**: 500x300px or similar
**Format**: PNG

---

### Linux Installation

#### 6. `linux-appimage-properties.png`
**Description**: Screenshot showing AppImage file properties with executable permission

**What to capture**:
- File manager (Nautilus/Dolphin/Thunar)
- Right-click → Properties
- Permissions tab
- "Allow executing file as program" checkbox

**Recommended size**: 600x500px or similar
**Format**: PNG

---

## Screenshot Guidelines

### Quality Standards
- **Resolution**: At least 72 DPI for web
- **Format**: PNG preferred (supports transparency)
- **Size**: Keep under 500KB per image
- **Clarity**: Text must be readable
- **Language**: English UI preferred

### Editing
- **Annotations**: Use red arrows or circles to highlight important elements
- **Privacy**: Blur any personal information (usernames, paths, etc.)
- **Consistency**: Use same annotation style across all screenshots
- **Borders**: Add subtle border if screenshot blends with background

### Tools Recommended
- **macOS**: Built-in Screenshot tool (Cmd+Shift+4)
- **Windows**: Snipping Tool or Snip & Sketch (Win+Shift+S)
- **Linux**: GNOME Screenshot, Spectacle, or Flameshot
- **Editing**: Preview (macOS), Paint (Windows), GIMP (all platforms)

### Annotation Tools
- **Skitch** (macOS) - Easy arrows and highlights
- **Greenshot** (Windows) - Screenshot + annotation
- **Flameshot** (Linux) - Powerful screenshot tool
- **Photopea** (Web) - Free Photoshop alternative

## File Naming Convention

Use descriptive, lowercase names with hyphens:
- ✅ `macos-right-click-open.png`
- ✅ `windows-smartscreen-warning.png`
- ❌ `Screenshot 2024-01-15.png`
- ❌ `IMG_1234.png`

## Placeholder Images

Until real screenshots are available, you can use placeholder images:

```markdown
![macOS Right-Click Open](../assets/screenshots/macos-right-click-open.png)
*Screenshot coming soon*
```

Or create simple diagrams using:
- **Excalidraw** (https://excalidraw.com) - Hand-drawn style diagrams
- **Draw.io** (https://app.diagrams.net) - Professional diagrams
- **Figma** (https://figma.com) - Design tool with screenshot capabilities

## Contributing Screenshots

If you'd like to contribute screenshots:

1. **Fork** the repository
2. **Capture** screenshots following guidelines above
3. **Edit** and annotate as needed
4. **Save** in this directory with proper naming
5. **Submit** pull request with description

### Checklist for Contributors
- [ ] Screenshot is clear and readable
- [ ] Personal information is blurred/removed
- [ ] File size is under 500KB
- [ ] Filename follows naming convention
- [ ] Annotations are clear and helpful
- [ ] Screenshot matches current OS version

## Current Status

| Screenshot | Status | Priority | Notes |
|------------|--------|----------|-------|
| macos-right-click-open.png | ❌ Missing | High | Critical for macOS users |
| macos-security-settings.png | ❌ Missing | High | Critical for macOS users |
| windows-smartscreen-warning.png | ❌ Missing | High | Critical for Windows users |
| windows-smartscreen-run-anyway.png | ❌ Missing | High | Critical for Windows users |
| windows-uac-prompt.png | ❌ Missing | Medium | Helpful but not critical |
| linux-appimage-properties.png | ❌ Missing | Low | Linux users are tech-savvy |

## Alternative: Video Tutorials

Consider creating video tutorials as an alternative or supplement to screenshots:

- **macOS Installation**: 2-3 minute walkthrough
- **Windows Installation**: 2-3 minute walkthrough
- **Linux Installation**: 1-2 minute walkthrough

Host on:
- YouTube (public, searchable)
- Vimeo (professional, ad-free)
- GitHub (embedded in README)

---

**Need help?** Contact the documentation team or open an issue on GitHub.
