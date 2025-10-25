# フロントエンド実装計画

## 1. 画面と UX
- `services/manager/frontend/src/features/auth/pages/login-page.tsx` と `register-page.tsx` に Google サインインボタンを追加し、Google ブランドガイドラインに沿ったスタイルとローディング状態を実装する。
- クリック時にバックエンドの `/techcv/api/v1/auth/google/login` を呼び出し、付与された `redirectTo` パラメータを引き継ぐ。処理は共通フック（例: `useGoogleLogin`）に切り出す。
- `AuthCallbackPage` で `code`、`state`、`error`、`redirectTo` のクエリを解析し、異常系（ユーザー拒否、state 欠如、トークン失敗）ごとに日本語メッセージを表示する。
- 成功時は `/api/v1/auth/google/callback` を叩いてクッキーを受け取り、続けて `/api/v1/auth/me` を取得して `useAuth` の `signIn` を更新し、`redirectTo` へ遷移する。

## 2. 状態管理と API 連携
- `apiClient` に対しては現在の `credentials: 'include'` を継続し、Google フロー専用のメソッドを utilities として整理する。
- エラーオブジェクトのパースを共通化し、API からのエラーコードを UI に反映できるようにする（CSRF 失敗、トークン期限切れなど）。
- 必要に応じて `sessionAtom` の状態に「Google 認証中」ステータスを追加し、ローディングスピナーなどでユーザーにフィードバックする。

## 3. ルーティングと設定
- 既存の `AppRouter` にエラー表示用のルートやクエリハンドリングを追加し、リダイレクト先を指定できるようにする。
- `.env.example` に `VITE_GOOGLE_CLIENT_ID` を追記し、README で Google OAuth 設定方法とローカル開発時のリダイレクト URI を説明する。
- 必要であれば `providers` や `features/auth/state` にクッキー更新イベントをフックするロジックを追加する。

## 4. テスト
- Login/Register ページがボタンを表示し、クリックでハンドラが呼び出されることを React Testing Library で確認する。
- `AuthCallbackPage` のユニットテストで、成功・state 不一致・エラー時のレンダリングが期待どおりになるか検証する。
- `vitest` と `MSW` などを使って API 通信をモックし、`/auth/me` 取得後に Jotai 状態が更新されることを確認する。
