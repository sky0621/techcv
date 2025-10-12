import { Button } from '@/components/ui/button';
import { appConfig } from '@/lib/env';

export const LoginPage = (): JSX.Element => {
  const handleLogin = () => {
    const redirectUrl = `${appConfig.apiBaseUrl}/api/v1/auth/google/login`;
    window.location.href = redirectUrl;
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md rounded-lg border bg-card p-8 shadow-sm">
        <div className="space-y-4 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">CV管理システムにサインイン</h1>
          <p className="text-sm text-muted-foreground">
            Googleアカウントでサインインし、CVを作成・管理できます。
          </p>
        </div>
        <Button className="mt-8 w-full" onClick={handleLogin}>
          Googleでサインイン
        </Button>
      </div>
    </div>
  );
};
