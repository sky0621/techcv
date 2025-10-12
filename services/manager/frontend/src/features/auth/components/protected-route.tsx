import { PropsWithChildren, useEffect } from 'react';
import { useNavigate, useRouterState } from '@tanstack/react-router';

import { LoadingScreen } from '@/components/system/loading-screen';
import { useAuth } from '@/features/auth/hooks/use-auth';

export const ProtectedRoute = ({ children }: PropsWithChildren): JSX.Element => {
  const { session, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const routerState = useRouterState();

  useEffect(() => {
    if (session.status === 'unauthenticated') {
      void navigate({
        to: '/login',
        search: {
          redirectTo: routerState.location.pathname
        },
        replace: true
      });
    }
  }, [navigate, routerState.location.pathname, session.status]);

  if (session.status === 'loading') {
    return <LoadingScreen message="サインイン状態を確認しています" />;
  }

  if (!isAuthenticated) {
    return <LoadingScreen message="リダイレクトしています" />;
  }

  return <>{children}</>;
};
