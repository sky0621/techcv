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
  idToken?: string;
}

export const sessionAtom = atom<SessionState>({
  status: 'loading'
});

export const isAuthenticatedAtom = atom(
  (get) => get(sessionAtom).status === 'authenticated' && Boolean(get(sessionAtom).idToken)
);
