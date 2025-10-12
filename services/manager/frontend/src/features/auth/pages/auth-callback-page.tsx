import { useEffect, useState } from 'react';
import { useNavigate } from '@tanstack/react-router';

import { LoadingScreen } from '@/components/system/loading-screen';
import { useAuth } from '@/features/auth/hooks/use-auth';
import { apiClient } from '@/lib/api-client';

interface CurrentUserResponse {
  user: {
    id: string;
    email: string;
    name: string;
    avatarUrl?: string;
  };
}

export const AuthCallbackPage = (): JSX.Element => {
  const navigate = useNavigate();
  const { signIn, setLoading, signOut } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const state = params.get('state') ?? undefined;
    const redirectTo = params.get('redirectTo') ?? '/';

    if (!code) {
      setError('認証コードが取得できませんでした。');
      return;
    }

    const finalizeAuth = async () => {
      try {
        setLoading();
        await apiClient.get('api/v1/auth/google/callback', {
          searchParams: {
            code,
            state
          }
        });

        const response = await apiClient.get('api/v1/auth/me').json<CurrentUserResponse>();
        signIn(response.user);
        await navigate({ to: redirectTo, replace: true });
      } catch (callbackError) {
        console.error(callbackError);
        setError('サインイン処理に失敗しました。もう一度お試しください。');
        signOut();
      }
    };

    void finalizeAuth();
  }, [navigate, setLoading, signIn, signOut]);

  if (error) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-background px-4 text-center">
        <div>
          <h1 className="text-xl font-semibold">サインインに失敗しました</h1>
          <p className="mt-2 text-sm text-muted-foreground">{error}</p>
        </div>
        <button
          type="button"
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          onClick={() => navigate({ to: '/login', replace: true })}
        >
          ログイン画面へ戻る
        </button>
      </div>
    );
  }

  return <LoadingScreen message="アカウント情報を取得しています" />;
};
