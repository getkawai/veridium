# SQLite Viewer

SQLite Viewer adalah aplikasi untuk melihat dan mengelola database SQLite `data/veridium.db` dengan antarmuka yang user-friendly.

## 🔒 Privacy Commitment

### Own Your Data, Command Your AI

Fitur SQLite Viewer adalah komitmen Kawai untuk menyediakan **Local AI** yang menghargai privasi pengguna. Semua operasi database dilakukan secara **lokal** di komputer Anda:

- ✅ **Data tetap di lokal** - Database SQLite disimpan lokal, tidak dikirim ke server eksternal
- ✅ **No cloud dependency** - Query dan manipulasi data berjalan sepenuhnya offline
- ✅ **Full control** - Anda memiliki akses penuh dan kontrol atas semua data Anda
- ✅ **Privacy-first** - Tidak ada tracking atau sharing data tanpa persetujuan. Debug logging (console.log) hanya ada di browser console untuk development/debugging saja, tidak ada transmisi eksternal atau telemetry. Query history disimpan lokal (localStorage) saja.
- ✅ **Sovereign computing** - AI Anda bekerja dengan data yang sepenuhnya Anda miliki

## Fitur

### 🗂️ Table Browser
- Daftar semua tabel dan views dalam database
- Search dan filter tabel
- Informasi jumlah rows per tabel
- Indikator tipe (TABLE/VIEW)

### 📊 Data Viewer
- Tampilan data tabel dengan pagination
- Sorting dan filtering data
- Search berdasarkan kolom tertentu
- Support berbagai operator (equals, contains, startsWith, endsWith)

### 🏗️ Structure Viewer
- Detail struktur tabel (kolom, tipe data, constraints)
- Informasi Primary Key dan Foreign Key
- Nullable dan default values
- Relationship antar tabel

### 💻 Query Editor
- Execute raw SQL queries
- Syntax highlighting
- Query history (tersimpan di localStorage)
- Sample queries untuk memulai
- Export hasil query

## Komponen

### Layout
- **SQLiteViewer**: Main container dengan DraggablePanel layout
- **TableSidebar**: Left sidebar dengan daftar tabel (resizable 200px-400px)
- **TableViewer**: Main content area dengan tabs

### Data Components
- **DataTable**: Menampilkan data tabel dengan pagination dan filtering
- **TableDetails**: Menampilkan struktur dan metadata tabel
- **QueryEditor**: Editor untuk menjalankan SQL queries

### Hooks
- **useSQLiteData**: Custom hook untuk mengelola semua operasi database

## API Services

Menggunakan Wails bindings dari:
- `internal/tableviewer/service.ts`
- `internal/tableviewer/models.ts`

### Available Methods
- `GetAllTables()` - Mendapatkan daftar semua tabel
- `GetTableData(tableName, pagination, filters)` - Data tabel dengan pagination
- `GetTableDetails(tableName)` - Struktur detail tabel
- `ExecuteRawQuery(query, args)` - Execute raw SQL

## Usage

1. Buka `/sqlite` di browser
2. Pilih tabel dari sidebar kiri
3. Browse data di tab "Data"
4. Lihat struktur di tab "Structure"  
5. Jalankan custom queries di tab "Query"

## Responsive Design

- Desktop: Sidebar kiri dengan main content
- Mobile: Collapsible sidebar (float mode)
- Draggable panel untuk resize sidebar

## Local Storage

- Query history disimpan di `sqlite-query-history`
- Maksimal 10 query terakhir

## Error Handling

- Graceful error handling untuk semua API calls
- User-friendly error messages
- Retry mechanisms untuk failed requests