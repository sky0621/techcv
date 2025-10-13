import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useRouterState } from '@tanstack/react-router';

import { LoadingScreen } from '@/components/system/loading-screen';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/features/auth/hooks/use-auth';
import { useVerifyRegistration } from '@/features/auth/hooks/use-verify-registration';
import { APIError } from '@/lib/api-error';

const AUTH_TOKEN_STORAGE_KEY = 'techcv_manager_auth_token';

export const VerifyRegistrationPage = (): JSX.Element => {
  const navigate = useNavigate();
  const routerState = useRouterState();
  const { signIn } = useAuth();
  const verifyRegistration = useVerifyRegistration();

  const [hasTriggered, setHasTriggered] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const token = useMemo(() => {
    const params = new URLSearchParams(routerState.location.searchStr);
    return params.get('token');
  }, [routerState.location.searchStr]);

  useEffect(() => {
    if (hasTriggered || verifyRegistration.isPending) {
      return;
    }

    if (!token) {
      setHasTriggered(true);
      setErrorMessage('確認トークンが見つかりません。再度ユーザー登録をお試しください。');
      return;
    }

    setHasTriggered(true);
    verifyRegistration.mutate(
      { token },
      {
        onSuccess: (result) => {
          localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, result.authToken);
          signIn({
            id: result.user.id,
            email: result.user.email,
            name: result.user.name ?? '',
            avatarUrl: undefined
          });

          void navigate({ to: '/' });
        },
        onError: (error) => {
          if (error instanceof APIError) {
            setErrorMessage(error.message);
          } else {
            setErrorMessage('確認に失敗しました。もう一度お試しください。');
          }
        }
      }
    );
  }, [hasTriggered, navigate, signIn, token, verifyRegistration]);

  if (verifyRegistration.isPending) {
    return <LoadingScreen message="メールアドレスを確認しています" />;
  }

  if (errorMessage) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background px-4 py-8">
        <div className="w-full max-w-md space-y-4 rounded-lg border bg-card p-8 text-center shadow-sm">
          <h1 className="text-xl font-semibold">確認に失敗しました</h1>
          <p className="text-sm text-muted-foreground">{errorMessage}</p>
          <div className="pt-2">
            <Button onClick={() => navigate({ to: '/register' })}>登録ページへ戻る</Button>
          </div>
        </div>
      </div>
    );
  }

  return <LoadingScreen message="ダッシュボードへリダイレクトしています" />;
};
