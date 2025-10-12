import { useParams } from '@tanstack/react-router';

export const PublicCVViewPage = (): JSX.Element => {
  const { publicId } = useParams({
    from: '/public/$publicId'
  });

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-3xl flex-col gap-6 px-6 py-16">
      <header className="space-y-2 text-center">
        <h1 className="text-3xl font-bold">公開CV</h1>
        <p className="text-sm text-muted-foreground">
          公開URL（ID: {publicId}）に紐づくCV情報を表示します。
        </p>
      </header>

      <div className="space-y-4 rounded-lg border bg-card p-8 shadow-sm">
        <p className="text-sm text-muted-foreground">
          公開用のCV表示コンポーネントを実装してください。バックエンドの`GET /api/v1/public/cv/:public_url`を利用します。
        </p>
      </div>
    </div>
  );
};
