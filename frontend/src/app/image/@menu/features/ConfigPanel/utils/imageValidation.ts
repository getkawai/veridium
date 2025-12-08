/**
 * Utility functions untuk validasi file gambar
 * Menyediakan validasi ukuran file, jumlah file, dan formatting
 */

/**
 * Format ukuran file ke format yang mudah dibaca manusia
 * @param bytes - Ukuran file dalam bytes
 * @returns String terformat seperti "1.5 MB"
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
};

/**
 * Hasil validasi dengan detail tambahan untuk pesan error
 */
export interface ValidationResult {
  actualSize?: number; // Ukuran file aktual
  error?: string; // Tipe error
  fileName?: string; // Nama file
  maxSize?: number; // Ukuran maksimal yang diizinkan
  valid: boolean; // Apakah validasi berhasil
}

/**
 * Validasi ukuran file gambar tunggal
 * @param file - File yang akan divalidasi
 * @param maxSize - Ukuran maksimal file dalam bytes, default 10MB jika tidak disediakan
 * @returns Hasil validasi
 */
export const validateImageFileSize = (file: File, maxSize?: number): ValidationResult => {
  const defaultMaxSize = 10 * 1024 * 1024; // 10MB default limit
  const actualMaxSize = maxSize ?? defaultMaxSize;

  if (file.size > actualMaxSize) {
    return {
      actualSize: file.size,
      error: 'fileSizeExceeded',
      fileName: file.name,
      maxSize: actualMaxSize,
      valid: false,
    };
  }

  return { valid: true };
};

/**
 * Validasi jumlah gambar
 * @param count - Jumlah gambar saat ini
 * @param maxCount - Jumlah maksimal yang diizinkan, skip validasi jika tidak disediakan
 * @returns Hasil validasi
 */
export const validateImageCount = (count: number, maxCount?: number): ValidationResult => {
  if (!maxCount) return { valid: true };

  if (count > maxCount) {
    return {
      error: 'imageCountExceeded',
      valid: false,
    };
  }

  return { valid: true };
};

/**
 * Validasi list file gambar
 * @param files - Array file yang akan divalidasi
 * @param constraints - Konfigurasi constraint (maxAddedFiles, maxFileSize)
 * @returns Hasil validasi, termasuk hasil validasi untuk setiap file
 */
export const validateImageFiles = (
  files: File[],
  constraints: {
    maxAddedFiles?: number; // Maksimal file yang bisa ditambahkan
    maxFileSize?: number; // Maksimal ukuran per file
  },
): {
  errors: string[]; // Array error yang ditemukan
  failedFiles?: ValidationResult[]; // Detail file yang gagal validasi
  fileResults: ValidationResult[]; // Hasil validasi per file
  valid: boolean; // Apakah semua file valid
} => {
  const errors: string[] = [];
  const fileResults: ValidationResult[] = [];
  const failedFiles: ValidationResult[] = [];

  // Validate file count
  const countResult = validateImageCount(files.length, constraints.maxAddedFiles);
  if (!countResult.valid && countResult.error) {
    errors.push(countResult.error);
  }

  // Validate each file
  files.forEach((file) => {
    const fileSizeResult = validateImageFileSize(file, constraints.maxFileSize);
    fileResults.push(fileSizeResult);

    if (!fileSizeResult.valid && fileSizeResult.error) {
      errors.push(fileSizeResult.error);
      failedFiles.push(fileSizeResult);
    }
  });

  return {
    errors: Array.from(new Set(errors)), // Remove duplicates
    failedFiles,
    fileResults,
    valid: errors.length === 0,
  };
};
