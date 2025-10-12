export const CVPreviewPage = (): JSX.Element => {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-4xl flex-col gap-6 px-6 py-16">
      <header className="space-y-2">
        <h1 className="text-3xl font-bold">CVプレビュー</h1>
        <p className="text-sm text-muted-foreground">
          公開されるCVがどのように表示されるかをリアルタイムに確認します。
        </p>
      </header>

      <div className="rounded-lg border bg-card p-8 shadow-sm">
        <p className="text-sm text-muted-foreground">
          バックエンドから取得したCVデータをここで整形して表示します。セクション毎に公開設定を反映し、PDF出力用のレイアウトも検討してください。
        </p>
      </div>
    </div>
  );
};
