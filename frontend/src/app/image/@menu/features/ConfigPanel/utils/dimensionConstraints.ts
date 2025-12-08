/**
 * Constraint untuk dimensi width dan height
 */
interface DimensionConstraints {
  height: { max: number; min: number };
  width: { max: number; min: number };
}

/**
 * Menyesuaikan dimensi gambar agar sesuai dengan constraint model sambil mempertahankan aspect ratio
 * Algoritma:
 * 1. Scale down jika melebihi maksimal
 * 2. Scale up jika di bawah minimal
 * 3. Round ke kelipatan 8 (requirement umum model AI)
 * 4. Final bounds check
 * 
 * @param originalWidth - Lebar gambar asli
 * @param originalHeight - Tinggi gambar asli
 * @param constraints - Constraint width dan height dari schema model
 * @returns Dimensi yang sudah disesuaikan dalam constraint
 */
export const constrainDimensions = (
  originalWidth: number,
  originalHeight: number,
  constraints: DimensionConstraints,
): { height: number; width: number } => {
  let width = originalWidth;
  let height = originalHeight;

  // First, scale down if exceeding maximum values
  if (width > constraints.width.max || height > constraints.height.max) {
    const scaleX = constraints.width.max / width;
    const scaleY = constraints.height.max / height;
    const scale = Math.min(scaleX, scaleY);

    width = Math.round(width * scale);
    height = Math.round(height * scale);
  }

  // Then, scale up if below minimum values
  if (width < constraints.width.min || height < constraints.height.min) {
    const scaleX = constraints.width.min / width;
    const scaleY = constraints.height.min / height;
    const scale = Math.max(scaleX, scaleY);

    width = Math.round(width * scale);
    height = Math.round(height * scale);
  }

  // Ensure final values are within bounds (may need adjustment due to rounding)
  width = Math.max(constraints.width.min, Math.min(constraints.width.max, width));
  height = Math.max(constraints.height.min, Math.min(constraints.height.max, height));

  // Round to nearest multiple of 8 (common model requirement)
  width = Math.round(width / 8) * 8;
  height = Math.round(height / 8) * 8;

  // Final bounds check after rounding
  width = Math.max(constraints.width.min, Math.min(constraints.width.max, width));
  height = Math.max(constraints.height.min, Math.min(constraints.height.max, height));

  return { height, width };
};
