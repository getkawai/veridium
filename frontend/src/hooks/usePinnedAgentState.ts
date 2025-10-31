import { useState, useEffect, useMemo, useCallback } from 'react';

export const usePinnedAgentState = () => {
  const [isPinned, setIsPinnedState] = useState(false);

  // Read initial state from URL on mount
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const pinnedParam = urlParams.get('pinned');
    if (pinnedParam === 'true') {
      setIsPinnedState(true);
    }
  }, []);

  // Update URL when state changes
  const setIsPinned = useCallback((value: boolean | ((prev: boolean) => boolean)) => {
    const newValue = typeof value === 'function' ? value(isPinned) : value;
    setIsPinnedState(newValue);

    const url = new URL(window.location.href);
    if (newValue) {
      url.searchParams.set('pinned', 'true');
    } else {
      url.searchParams.delete('pinned');
    }

    // Update URL without page reload
    window.history.replaceState({}, '', url.toString());
  }, [isPinned]);

  const actions = useMemo(
    () => ({
      pinAgent: () => setIsPinned(true),
      setIsPinned,
      togglePinAgent: () => setIsPinned((prev) => !prev),
      unpinAgent: () => setIsPinned(false),
    }),
    [setIsPinned],
  );

  return [isPinned, actions] as const;
};
