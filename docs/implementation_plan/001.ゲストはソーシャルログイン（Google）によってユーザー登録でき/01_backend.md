# バックエンド実装計画

## 1. データ層と設定の整備
- `services/manager/backend/db/schema.sql` に `users` テーブルを追加し、UUID v7 主キー、`password_hash` の NULL 許可、`google_id` のユニーク制約、プロフィール関連カラム、監査カラム、必要なインデックスを定義する。
- `services/manager/backend/db/queries` に Google ID 検索・作成・最終ログイン更新などの SQL を追加し、`make tidy`（または `sqlc generate`）でコードを再生成する。
- OAuth state 管理のためのストレージ（MySQL テーブル or Redis）を設計し、sqldef マイグレーションと SQLc バインディング、もしくは Redis 用の設定を追加する。
- 新規環境変数（`GOOGLE_CLIENT_ID`、`GOOGLE_CLIENT_SECRET`、`GOOGLE_REDIRECT_URI`、`SESSION_SECRET`、`COOKIE_DOMAIN`、`COOKIE_SECURE`、Redis 接続情報など）を README と `.env.example` に明記し、`cmd/api` 起動時にバリデーションする。

## 2. ドメイン・ユースケースの拡張
- `services/manager/backend/internal/domain/user` に `GoogleID` 値オブジェクトを追加し、フォーマット検証を実装する。
- `user.go` を拡張して任意の `passwordHash` と `googleID` を扱えるようにし、`NewUserWithGoogle`、`LinkGoogleAccount`、`WithLastLogin` などの補助メソッドを追加する。
- `repository.go` に `FindByGoogleID`、`UpdateGoogleID`、`UpdateLastLogin` などのメソッドを定義し、ユースケース層から利用できるようにする。
- `services/manager/backend/internal/usecase/auth`（もしくは `usecase/user/command`）に以下のユースケースを追加する:
  - `StartGoogleLogin`: state 生成と認証 URL 作成。
  - `CompleteGoogleCallback`: state 検証、トークン交換、ID トークン検証、既存ユーザー判定、新規登録、JWT 発行、クッキー設定までを内包。
- 既存テストを参考にユースケースの table-driven テストを作成し、state 不一致・既存ユーザー・新規作成・トークン検証失敗などのケースを網羅する。

## 3. インフラ層とアダプタの実装
- `internal/infrastructure/mysql` に SQLc で生成したコードを利用する `UserRepository` を追加し、MySQL のエラーをドメインエラーへ正しくマッピングする。
- トランザクションマネージャーを `*sql.DB`/`*sql.Tx` ベースに実装し、ユースケースから `WithinTransaction` が利用できるようにする。
- `internal/infrastructure/oauth/google_client.go` を新設し、`golang.org/x/oauth2` を使った stateful OAuth クライアント、トークン交換、`google.golang.org/api/idtoken` による ID トークン検証（JWKS キャッシュ含む）を実装する。
- `internal/infrastructure/auth/jwt_service.go` にセッショントークン発行機能を追加し、HTTP Only / Secure / SameSite=None を制御できる設定を用意する。
- `SessionRepository` インターフェースを定義し、開発用のインメモリ実装と本番用 Redis 実装（`internal/infrastructure/session/redis_repository.go` など）を準備する。
- 重要イベント（state 生成、トークン交換、ユーザー作成/更新、エラー）を構造化ログで記録する。

## 4. HTTP インターフェースと OpenAPI
- `services/manager/openapi/spec` に以下を追加し、生成コードを更新する:
  - `GET /auth/google/login`: 認証 URL を返すか 302 リダイレクトを行う仕様。
  - `GET /auth/google/callback`: 成功時は 204（クッキー設定のみ）または JSON を返し、失敗時は統一されたエラーレスポンスを返却。
  - 必要なリクエスト・レスポンススキーマとエラーコード。
- `internal/interface/http/handler` に Google 認証ハンドラを実装し、ユースケースを呼び出してレスポンス生成・クッキー設定・リダイレクト処理を行う。
- `/techcv/api/v1` グループに新規エンドポイントを登録し、既存ミドルウェア（リクエスト ID、Timeout、ログ）と統合する。

## 5. アプリケーション初期化
- `cmd/api/main.go` を更新し、MySQL リポジトリ、Redis セッションストア、Google OAuth クライアント、JWT サービス、追加ユースケースを組み立ててハンドラに依存性注入する。
- 起動時に必須環境変数が無い場合はエラーで終了させる。シャットダウン時に DB / Redis 接続をクリーンに閉じる。
- `Makefile` に SQLc・OpenAPI の再生成コマンドや Redis 開発環境の起動ターゲットを追加する（必要に応じて）。
