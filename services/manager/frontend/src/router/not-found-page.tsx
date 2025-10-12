import { Link } from '@tanstack/react-router';

export const NotFoundPage = (): JSX.Element => (
  <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-background text-center text-foreground">
    <div>
      <p className="text-sm font-semibold text-muted-foreground">404</p>
      <h1 className="mt-2 text-2xl font-bold">お探しのページが見つかりません</h1>
      <p className="mt-2 text-sm text-muted-foreground">
        URLが正しいか確認するか、ダッシュボードへ戻ってください。
      </p>
    </div>
    <Link
      to="/"
      className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
    >
      ダッシュボードへ戻る
    </Link>
  </div>
);
