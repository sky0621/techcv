import { useAtomValue, useSetAtom } from 'jotai';
import { useCallback } from 'react';
import type { User } from 'firebase/auth';

import { ensureAuthSession } from '@/features/auth/api/auth-client';
import {
  AuthUser,
  SessionState,
  isAuthenticatedAtom,
  sessionAtom
} from '@/features/auth/state/session';
import { setApiAuthToken } from '@/lib/api-client';
import { signInWithGooglePopup, signOutFromFirebase } from '@/lib/firebase';

export const useAuth = () => {
  const session = useAtomValue(sessionAtom);
  const isAuthenticated = useAtomValue(isAuthenticatedAtom);
  const setSession = useSetAtom(sessionAtom);

  const setLoading = useCallback(() => setSession({ status: 'loading' }), [setSession]);

  const clearSession = useCallback(() => {
    setApiAuthToken(null);
    setSession({ status: 'unauthenticated' });
  }, [setSession]);

  const signIn = useCallback(
    (user: AuthUser, idToken: string) => {
      setApiAuthToken(idToken);
      setSession({
        status: 'authenticated',
        user,
        idToken
      });
    },
    [setSession]
  );

  const signOut = useCallback(async () => {
    await signOutFromFirebase();
    clearSession();
  }, [clearSession]);

  const signInWithGoogle = useCallback(async () => {
    await signInWithGooglePopup();
  }, []);

  const syncSessionWithFirebaseUser = useCallback(
    async (firebaseUser: User, idToken: string) => {
      const authResult = await ensureAuthSession(idToken);
      const name = authResult.displayName ?? firebaseUser.displayName ?? authResult.email;
      const avatarUrl = firebaseUser.photoURL ?? undefined;

      signIn(
        {
          id: authResult.userId,
          email: authResult.email,
          name,
          avatarUrl
        },
        idToken
      );
    },
    [signIn]
  );

  const updateSession = useCallback(
    (next: SessionState) => {
      setSession(next);
      setApiAuthToken(next.status === 'authenticated' ? next.idToken ?? null : null);
    },
    [setSession]
  );

  return {
    session,
    isAuthenticated,
    setLoading,
    signInWithGoogle,
    signOut,
    updateSession,
    syncSessionWithFirebaseUser,
    clearSession
  };
};
