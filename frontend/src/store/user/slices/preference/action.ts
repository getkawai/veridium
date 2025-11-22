import type { StateCreator } from 'zustand/vanilla';

import type { UserStore } from '@/store/user';
import { UserGuide, UserPreference } from '@/types/user';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';
import { AsyncLocalStorage } from '@/utils/localStorage';

const preferenceStorage = new AsyncLocalStorage<UserPreference>('LOBE_PREFERENCE');

const n = setNamespace('preference');

export interface PreferenceAction {
  updateGuideState: (guide: Partial<UserGuide>) => Promise<void>;
  updatePreference: (preference: Partial<UserPreference>, action?: any) => Promise<void>;
}

export const createPreferenceSlice: StateCreator<
  UserStore,
  [['zustand/devtools', never]],
  [],
  PreferenceAction
> = (set, get) => ({
  updateGuideState: async (guide) => {
    const { updatePreference } = get();
    const nextGuide = merge(get().preference.guide, guide);
    await updatePreference({ guide: nextGuide });
  },

  updatePreference: async (preference, action) => {
    const nextPreference = merge(get().preference, preference);

    set({ preference: nextPreference }, false, action || n('updatePreference'));

    // Note: Preference is stored in LocalStorage, not DB
    await preferenceStorage.saveToLocalStorage(nextPreference);
    
    console.log('[User] Updated preference via LocalStorage');
  },
});
