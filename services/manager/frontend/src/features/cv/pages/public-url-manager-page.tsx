import { Button } from '@/components/ui/button';

export const PublicURLManagerPage = (): JSX.Element => {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-3xl flex-col gap-6 px-6 py-16">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold">公開URLの管理</h1>
        <p className="text-sm text-muted-foreground">
          公開URLの発行、リセット、状態確認を行います。
        </p>
      </header>

      <section className="space-y-4 rounded-lg border bg-card p-6 shadow-sm">
        <div className="space-y-1">
          <h2 className="text-lg font-semibold">現在の公開URL</h2>
          <p className="text-sm text-muted-foreground">
            APIから取得した公開URL情報を表示してください。
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button>URLを生成</Button>
          <Button variant="outline">URLを無効化</Button>
        </div>
      </section>
    </div>
  );
};
