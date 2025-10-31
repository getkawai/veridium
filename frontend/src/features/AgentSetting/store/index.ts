'use client';

import { StoreApiWithSelector } from '@/types';
import { createContext } from 'zustand-utils';
import { shallow } from 'zustand/shallow';
import { createWithEqualityFn } from 'zustand/traditional';

import { Store, store } from './action';

export type { State } from './initialState';

export const createStore = () => createWithEqualityFn(store, shallow);

export const { useStore, useStoreApi, Provider } = createContext<StoreApiWithSelector<Store>>();

export { selectors } from './selectors';
