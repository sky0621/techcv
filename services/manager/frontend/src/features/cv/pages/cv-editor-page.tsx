import { Button } from '@/components/ui/button';

export const CVEditorPage = (): JSX.Element => {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-6 py-16">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">CV編集</h1>
          <p className="text-sm text-muted-foreground">
            基本情報や職務経歴などを入力し、公開可否を設定します。
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline">下書きを保存</Button>
          <Button>更新内容を公開</Button>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <section className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-lg font-semibold">基本情報</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            氏名、連絡先情報などのプロフィールを入力します。
          </p>
          <div className="mt-4 space-y-3 text-sm text-muted-foreground">
            <p>フォーム要素を実装してください（氏名、メール、電話、住所、公開設定など）。</p>
          </div>
        </section>

        <section className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-lg font-semibold">職務経歴</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            経験した会社や役割、期間を追加します。
          </p>
          <div className="mt-4 space-y-3 text-sm text-muted-foreground">
            <p>複数の職務経歴を追加・並び替えられるフォームを実装予定です。</p>
          </div>
        </section>

        <section className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-lg font-semibold">スキル</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            スキル名、レベル、経験年数、公開可否を管理します。</p>
          <div className="mt-4 space-y-3 text-sm text-muted-foreground">
            <p>スキル用のテーブルまたはフォームレイアウトを追加してください。</p>
          </div>
        </section>

        <section className="rounded-lg border bg-card p-6 shadow-sm">
          <h2 className="text-lg font-semibold">学歴・資格・プロジェクト</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            セクションごとにカードまたはタブで管理する予定です。
          </p>
          <div className="mt-4 space-y-3 text-sm text-muted-foreground">
            <p>タブUIやアコーディオンUIを用いた構成を検討してください。</p>
          </div>
        </section>
      </div>
    </div>
  );
};
