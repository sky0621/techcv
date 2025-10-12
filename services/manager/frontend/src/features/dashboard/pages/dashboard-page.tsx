import { Link } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';

export const DashboardPage = (): JSX.Element => {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-5xl flex-col gap-8 px-6 py-16">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">ダッシュボード</h1>
        <p className="text-sm text-muted-foreground">
          CVの管理、公開設定、共有URLの管理をここから行えます。
        </p>
      </header>

      <section className="grid gap-6 md:grid-cols-2">
        <div className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-xl font-semibold">CVを編集</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            基本情報、職務経歴、スキルなどのCV内容を編集します。
          </p>
          <Button asChild className="mt-4">
            <Link to="/cv/edit">編集画面へ</Link>
          </Button>
        </div>

        <div className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-xl font-semibold">CVプレビュー</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            公開前にCVの表示内容をプレビューで確認できます。
          </p>
          <Button asChild className="mt-4" variant="secondary">
            <Link to="/cv/preview">プレビューを見る</Link>
          </Button>
        </div>

        <div className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-xl font-semibold">公開URLの管理</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            公開URLの発行・無効化を管理し、共有設定をコントロールします。
          </p>
          <Button asChild className="mt-4" variant="outline">
            <Link to="/settings/public-url">公開URLを管理</Link>
          </Button>
        </div>
      </section>
    </div>
  );
};
