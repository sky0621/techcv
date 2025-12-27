import { PropsWithChildren, useEffect } from 'react';

import { useAuth } from '@/features/auth/hooks/use-auth';
import { onFirebaseIdTokenChanged } from '@/lib/firebase';

export const AuthInitializer = ({ children }: PropsWithChildren): JSX.Element => {
  const { setLoading, clearSession, syncSessionWithFirebaseUser } = useAuth();

  useEffect(() => {
    setLoading();

    const unsubscribe = onFirebaseIdTokenChanged(async (firebaseUser) => {
      if (!firebaseUser) {
        clearSession();
        return;
      }

      try {
        const idToken = await firebaseUser.getIdToken();
        await syncSessionWithFirebaseUser(firebaseUser, idToken);
      } catch (error) {
        console.error('Failed to synchronize auth session', error);
        clearSession();
      }
    });

    return unsubscribe;
  }, [clearSession, setLoading, syncSessionWithFirebaseUser]);

  return <>{children}</>;
};
