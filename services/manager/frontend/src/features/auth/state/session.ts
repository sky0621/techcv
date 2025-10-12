import { atom } from 'jotai';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
}

export type SessionStatus = 'authenticated' | 'unauthenticated' | 'loading';

export interface SessionState {
  status: SessionStatus;
  user?: AuthUser;
}

export const sessionAtom = atom<SessionState>({
  status: 'unauthenticated'
});

export const isAuthenticatedAtom = atom((get) => get(sessionAtom).status === 'authenticated');
