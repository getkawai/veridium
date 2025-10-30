/**
 * Vector similarity search utilities for SQLite
 * Since SQLite doesn't have native vector support like pgvector,
 * we implement cosine similarity in JavaScript/SQL
 */

/**
 * Convert a Float32Array or number array to a Buffer for SQLite blob storage
 */
export function vectorToBuffer(vector: number[] | Float32Array): Buffer {
  const float32Array = vector instanceof Float32Array ? vector : new Float32Array(vector);
  return Buffer.from(float32Array.buffer);
}

/**
 * Convert a Buffer from SQLite back to a number array
 */
export function bufferToVector(buffer: Buffer): number[] {
  const float32Array = new Float32Array(
    buffer.buffer,
    buffer.byteOffset,
    buffer.length / Float32Array.BYTES_PER_ELEMENT,
  );
  return Array.from(float32Array);
}

/**
 * Calculate cosine similarity between two vectors
 * Returns a value between -1 and 1, where 1 means identical vectors
 */
export function cosineSimilarity(vecA: number[], vecB: number[]): number {
  if (vecA.length !== vecB.length) {
    throw new Error('Vectors must have the same length');
  }

  let dotProduct = 0;
  let normA = 0;
  let normB = 0;

  for (let i = 0; i < vecA.length; i++) {
    dotProduct += vecA[i] * vecB[i];
    normA += vecA[i] * vecA[i];
    normB += vecB[i] * vecB[i];
  }

  normA = Math.sqrt(normA);
  normB = Math.sqrt(normB);

  if (normA === 0 || normB === 0) {
    return 0;
  }

  return dotProduct / (normA * normB);
}

/**
 * Calculate Euclidean distance between two vectors
 * Lower values indicate more similar vectors
 */
export function euclideanDistance(vecA: number[], vecB: number[]): number {
  if (vecA.length !== vecB.length) {
    throw new Error('Vectors must have the same length');
  }

  let sum = 0;
  for (let i = 0; i < vecA.length; i++) {
    const diff = vecA[i] - vecB[i];
    sum += diff * diff;
  }

  return Math.sqrt(sum);
}

/**
 * Find the top K most similar vectors from a list
 * @param queryVector The vector to compare against
 * @param vectors Array of objects containing vectors and associated data
 * @param k Number of top results to return
 * @param similarityFn Similarity function to use (default: cosine similarity)
 * @returns Array of top K most similar items with their similarity scores
 */
export function findTopKSimilar<T extends { vector: number[] | Buffer; [key: string]: any }>(
  queryVector: number[],
  vectors: T[],
  k: number = 10,
  similarityFn: (a: number[], b: number[]) => number = cosineSimilarity,
): Array<T & { similarity: number }> {
  // Calculate similarities
  const withSimilarities = vectors.map((item) => {
    const vector = item.vector instanceof Buffer ? bufferToVector(item.vector) : item.vector;
    const similarity = similarityFn(queryVector, vector);
    return {
      ...item,
      similarity,
    };
  });

  // Sort by similarity (descending for cosine, ascending for distance metrics)
  withSimilarities.sort((a, b) => b.similarity - a.similarity);

  // Return top K
  return withSimilarities.slice(0, k);
}

/**
 * Batch process vector similarity searches
 * Useful for large datasets where you want to process in chunks
 */
export async function batchVectorSearch<T extends { vector: number[] | Buffer; [key: string]: any }>(
  queryVector: number[],
  vectors: T[],
  options: {
    k?: number;
    batchSize?: number;
    similarityFn?: (a: number[], b: number[]) => number;
    onProgress?: (processed: number, total: number) => void;
  } = {},
): Promise<Array<T & { similarity: number }>> {
  const {
    k = 10,
    batchSize = 1000,
    similarityFn = cosineSimilarity,
    onProgress,
  } = options;

  const results: Array<T & { similarity: number }> = [];

  for (let i = 0; i < vectors.length; i += batchSize) {
    const batch = vectors.slice(i, i + batchSize);
    const batchResults = findTopKSimilar(queryVector, batch, k, similarityFn);
    results.push(...batchResults);

    onProgress?.(Math.min(i + batchSize, vectors.length), vectors.length);

    // Allow event loop to breathe
    await new Promise((resolve) => setTimeout(resolve, 0));
  }

  // Sort all results and return top K
  results.sort((a, b) => b.similarity - a.similarity);
  return results.slice(0, k);
}

/**
 * Generate SQL for creating a functional index for vector search optimization
 * Note: This won't make search faster for all cases, but can help with filtering
 */
export function generateVectorIndexSQL(tableName: string, vectorColumn: string): string {
  return `
    -- Create an index on the vector column length for filtering
    CREATE INDEX IF NOT EXISTS idx_${tableName}_${vectorColumn}_length 
    ON ${tableName}(length(${vectorColumn}));
  `.trim();
}

