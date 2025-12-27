import { useEffect, useState } from 'react';
import { useNavigate, useRouterState } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';
import { useAuth } from '@/features/auth/hooks/use-auth';

export const LoginPage = (): JSX.Element => {
  const { signInWithGoogle, session } = useAuth();
  const [isSigningIn, setIsSigningIn] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const routerState = useRouterState();

  useEffect(() => {
    if (session.status === 'authenticated') {
      const search = routerState.location.search as Record<string, unknown>;
      const redirectTo =
        typeof search.redirectTo === 'string' && search.redirectTo.startsWith('/')
          ? search.redirectTo
          : '/';
      void navigate({ to: redirectTo as '/' });
    }
  }, [navigate, routerState.location.search, session.status]);

  const handleLogin = async () => {
    setError(null);
    setIsSigningIn(true);
    try {
      await signInWithGoogle();
    } catch (err) {
      console.error('Failed to sign in with Google', err);
      setError('サインインに失敗しました。時間をおいて再試行してください。');
    } finally {
      setIsSigningIn(false);
    }
  };

  const isButtonDisabled = isSigningIn || session.status === 'loading';

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md rounded-lg border bg-card p-8 shadow-sm">
        <div className="space-y-4 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">CV管理システムにサインイン</h1>
          <p className="text-sm text-muted-foreground">
            Googleアカウントでサインインし、CVを作成・管理できます。
          </p>
        </div>
        <Button className="mt-8 w-full" onClick={handleLogin} disabled={isButtonDisabled}>
          {isSigningIn ? 'サインイン中...' : 'Googleでサインイン'}
        </Button>
        {error ? (
          <p className="mt-4 text-center text-sm text-destructive">{error}</p>
        ) : null}
        <p className="mt-6 text-center text-sm text-muted-foreground">
          サインイン後にCVの作成・管理ができます。
        </p>
      </div>
    </div>
  );
};
