import { useCallback, useState } from 'react';

export interface UseAsyncActionOptions<T> {
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
}

/**
 * Hook for handling async actions with loading state
 * Replacement for useActionSWR
 */
export function useAsyncAction<T extends (...args: any[]) => Promise<any>>(
  action: T,
  options?: UseAsyncActionOptions<Awaited<ReturnType<T>>>,
): {
  mutate: T;
  isValidating: boolean;
  error: Error | undefined;
} {
  const [isValidating, setIsValidating] = useState(false);
  const [error, setError] = useState<Error | undefined>();

  const mutate = useCallback(
    (async (...args: Parameters<T>) => {
      setIsValidating(true);
      setError(undefined);
      try {
        const result = await action(...args);
        options?.onSuccess?.(result);
        return result;
      } catch (err) {
        const error = err as Error;
        setError(error);
        options?.onError?.(error);
        throw err;
      } finally {
        setIsValidating(false);
      }
    }) as T,
    [action, options],
  );

  return { mutate, isValidating, error };
}

