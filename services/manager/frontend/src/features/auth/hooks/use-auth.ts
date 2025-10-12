import { useAtomValue, useSetAtom } from 'jotai';

import {
  AuthUser,
  SessionState,
  isAuthenticatedAtom,
  sessionAtom
} from '@/features/auth/state/session';

export const useAuth = () => {
  const session = useAtomValue(sessionAtom);
  const isAuthenticated = useAtomValue(isAuthenticatedAtom);
  const setSession = useSetAtom(sessionAtom);

  const setLoading = () => setSession({ status: 'loading' });

  const signIn = (user: AuthUser) =>
    setSession({
      status: 'authenticated',
      user
    });

  const signOut = () => setSession({ status: 'unauthenticated' });

  const updateSession = (next: SessionState) => setSession(next);

  return {
    session,
    isAuthenticated,
    setLoading,
    signIn,
    signOut,
    updateSession
  };
};
