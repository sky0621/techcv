# ゲストはソーシャルログイン（Google）によってユーザー登録できる

## ユーザーストーリー

**As a** ゲスト（未登録の利用者）  
**I want** Googleアカウントを使ってユーザー登録できる  
**So that** パスワードを設定せずに簡単にシステムにログインしてCVの管理機能を利用できる

## 概要

ゲストがmanagerサービスを利用するために、Googleアカウントを使用してユーザー登録を行う機能を提供します。OAuth 2.0プロトコルを使用してGoogleで認証を行い、登録が完了すると、ユーザーとしてシステムにログインし、CV管理機能にアクセスできるようになります。

## 受け入れ基準（Acceptance Criteria）

### 1. ソーシャルログインボタンの表示

**WHEN** ゲストがユーザー登録ページにアクセスする  
**THEN** システムは「Googleでログイン」ボタンを表示する

**WHEN** ゲストがログインページにアクセスする  
**THEN** システムは「Googleでログイン」ボタンを表示する

### 2. Google認証フローの開始

**WHEN** ゲストが「Googleでログイン」ボタンをクリックする  
**THEN** システムは以下の処理を実行する
- OAuth 2.0認証フローを開始する
- stateパラメータを生成してセッションに保存する（CSRF対策）
- Googleの認証ページにリダイレクトする
- 必要なスコープ（email、profile）をリクエストする

### 3. Google認証ページの表示

**WHEN** ゲストがGoogleの認証ページにリダイレクトされる  
**THEN** Googleは以下の情報を表示する
- アプリケーション名
- 要求される権限（メールアドレス、基本プロフィール情報）
- 許可/拒否のオプション

### 4. 認証の許可

**WHEN** ゲストがGoogleの認証ページで「許可」をクリックする  
**THEN** GoogleはシステムのコールバックURLにリダイレクトする
- 認証コード（authorization code）をクエリパラメータとして含める
- stateパラメータを返す

### 5. 認証の拒否

**WHEN** ゲストがGoogleの認証ページで「拒否」をクリックする  
**THEN** GoogleはシステムのコールバックURLにリダイレクトする
- エラー情報をクエリパラメータとして含める

**WHEN** システムがエラー情報を受け取る  
**THEN** システムは「Google認証がキャンセルされました」というメッセージを表示する

### 6. stateパラメータの検証

**WHEN** システムがGoogleからのコールバックを受け取る  
**THEN** システムは返されたstateパラメータがセッションに保存されたものと一致することを検証する

**IF** stateパラメータが一致しない  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 7. アクセストークンの取得

**WHEN** システムが認証コードを受け取る  
**AND** stateパラメータの検証が成功する  
**THEN** システムは以下の処理を実行する
- 認証コードを使用してGoogleのトークンエンドポイントにリクエストを送信する
- アクセストークンとIDトークンを取得する

**IF** トークンの取得に失敗する  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 8. ユーザー情報の取得

**WHEN** システムがアクセストークンを取得する  
**THEN** システムは以下の処理を実行する
- IDトークンを検証する（署名、有効期限、issuer、audience）
- IDトークンからユーザー情報を抽出する（sub、email、name、picture）

**IF** IDトークンの検証に失敗する  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 9. 既存ユーザーの確認

**WHEN** システムがGoogleからユーザー情報を取得する  
**THEN** システムはGoogleのユーザーID（sub）を使用して既存ユーザーを検索する

**IF** 既存ユーザーが見つかる  
**THEN** システムはログイン処理を実行する（新規登録ではなく）

### 10. 新規ユーザー登録の実行

**WHEN** 既存ユーザーが見つからない  
**THEN** システムは以下の処理を実行する
- UUID v7形式でユーザーIDを生成する
- Googleから取得した情報を使用して新しいユーザーレコードを作成する
  - email: Googleから取得したメールアドレス
  - name: Googleから取得した名前
  - google_id: GoogleのユーザーID（sub）
  - profile_image: Googleから取得したプロフィール画像URL（オプション）
- password_hashはNULLに設定する（ソーシャルログインのため）
- created_at、updated_atをUTCの現在時刻で記録する
- email_verified_atをUTCの現在時刻で記録する（Googleで検証済み）
- is_activeを1（有効）に設定する
- データベースに保存する

### 11. メールアドレスの重複チェック

**WHEN** システムが新規ユーザーを登録しようとする  
**AND** Googleから取得したメールアドレスが既に別のユーザーで登録されている  
**THEN** システムは以下の処理を実行する
- 既存のユーザーレコードにgoogle_idを追加する
- ログイン処理を実行する

### 12. 認証トークンの生成

**WHEN** ユーザー登録またはログインが正常に完了する  
**THEN** システムは以下の処理を実行する
- JWTトークンを生成する
- トークンにユーザーID、メールアドレス、発行日時、有効期限を含める
- トークンをクライアントに返す

### 13. 登録成功時の処理

**WHEN** ユーザー登録が正常に完了する  
**THEN** システムは以下の処理を実行する
- ユーザーを自動的にログイン状態にする
- 認証トークンをローカルストレージまたはクッキーに保存する
- ダッシュボードページにリダイレクトする
- 「登録が完了しました」という成功メッセージを表示する

### 14. ログイン成功時の処理

**WHEN** 既存ユーザーのログインが正常に完了する  
**THEN** システムは以下の処理を実行する
- last_login_atをUTCの現在時刻で更新する
- 認証トークンをローカルストレージまたはクッキーに保存する
- ダッシュボードページにリダイレクトする
- 「ログインしました」という成功メッセージを表示する

### 15. エラーハンドリング

**IF** Google APIとの通信中にネットワークエラーが発生する  
**THEN** システムは「ネットワークエラーが発生しました。再度お試しください」というエラーメッセージを表示する

**IF** データベースエラーやその他の予期しないエラーが発生する  
**THEN** システムは「登録処理中にエラーが発生しました。しばらくしてから再度お試しください」というエラーメッセージを表示する

### 16. セキュリティ要件

**WHEN** システムがOAuth 2.0フローを実装する  
**THEN** システムは以下のセキュリティ対策を実装する
- stateパラメータを使用してCSRF攻撃を防ぐ
- HTTPSを使用して通信を暗号化する
- IDトークンの署名を検証する
- トークンの有効期限を検証する
- トークンのissuerとaudienceを検証する

**WHEN** システムがアクセストークンを保存する  
**THEN** システムはトークンを安全に保存し、適切な有効期限を設定する

### 17. レスポンシブデザイン

**WHEN** ゲストがモバイルデバイスから登録ページにアクセスする  
**THEN** システムはモバイル画面に最適化されたレイアウトで「Googleでログイン」ボタンを表示する

## 技術的な制約

### バックエンド技術スタック
- OAuth 2.0クライアントにはgolang.org/x/oauth2を使用する
- Google OAuth 2.0エンドポイント:
  - 認証エンドポイント: https://accounts.google.com/o/oauth2/v2/auth
  - トークンエンドポイント: https://oauth2.googleapis.com/token
  - ユーザー情報エンドポイント: IDトークンから取得
- 必要なスコープ: openid, email, profile
- リダイレクトURIは環境変数で設定可能にする
- Google Client IDとClient Secretは環境変数で管理する
- stateパラメータはセッションまたはクッキーで管理する（有効期限10分）
- JWTトークンの有効期限は24時間とする
- 日時はすべてUTCで保存し、DATETIME(6)型を使用する（マイクロ秒精度）

### データベーススキーマ要件

**usersテーブルの拡張**:
- `google_id` VARCHAR(255) UNIQUE - GoogleのユーザーID（sub）
- `profile_image` VARCHAR(500) - プロフィール画像URL
- `password_hash` VARCHAR(255) - ソーシャルログインの場合はNULL許可に変更
- INDEX on `google_id` for fast lookup

既存のカラム:
- `id` BINARY(16) PRIMARY KEY - UUID v7
- `email` VARCHAR(255) NOT NULL UNIQUE
- `name` VARCHAR(100)
- `bio` TEXT
- `is_active` TINYINT(1) NOT NULL DEFAULT 1
- `email_verified_at` DATETIME(6)
- `last_login_at` DATETIME(6)
- `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
- `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
- `deleted_at` DATETIME(6)

### フロントエンド技術スタック
- React 18+ with TypeScript
- TanStack Router for routing
- TanStack Query for data fetching
- Jotai for global state management
- ky for HTTP client
- shadcn/ui for UI components
- OAuth 2.0フローはバックエンドで処理し、フロントエンドはリダイレクトのみを行う

## 非機能要件

- Google認証フロー全体は30秒以内に完了する
- ユーザー登録処理は3秒以内に完了する
- 「Googleでログイン」ボタンはGoogleのブランドガイドラインに準拠する
- すべてのエラーは明確で理解しやすいメッセージで表示する
- APIエラーレスポンスは統一されたフォーマット（requestId、code、details）で返す
- ログはCloudLogging形式のJSON構造化ログとして出力する
- ログにはcontext由来のrequest_idを自動付与する
- OAuth 2.0フローのすべてのステップをログに記録する（デバッグ用）

## 関連するユビキタス言語

- **ゲスト（guest）**: まだユーザー登録していない利用者
- **ユーザー（user）**: ユーザー登録が済んだ利用者
- **manager**: WebエンジニアのCVを管理するサービス
- **ソーシャルログイン（social login）**: 外部サービス（Google等）のアカウントを使用した認証方式
- **OAuth 2.0**: 認可のための業界標準プロトコル
- **IDトークン（ID token）**: ユーザーの認証情報を含むJWTトークン
- **アクセストークン（access token）**: リソースへのアクセス権を表すトークン

## アーキテクチャ上の考慮事項

### バックエンド（Clean Architecture + DDD + CQRS）

- **Domain層**: User集約、Email値オブジェクト、GoogleID値オブジェクトを定義
- **UseCase層**: RegisterUserWithGoogleコマンド、LoginWithGoogleコマンドを実装
- **Adapter層**: OAuth 2.0コールバックハンドラー、HTTPハンドラー
- **Infrastructure層**: 
  - sqlcを使用したリポジトリ実装
  - Google OAuth 2.0クライアント実装
  - JWTトークン生成サービス
- **CQRS**: ユーザー登録はコマンド側で実装（集約を使用）
- **トランザクション**: User集約の保存は単一トランザクションで実行
- **バリデーション**: IDトークンの検証、ユーザー情報の検証

### フロントエンド（レイヤードアーキテクチャ）

- **Presentation層**: 
  - 登録/ログインページコンポーネント
  - Googleログインボタンコンポーネント
  - OAuth 2.0コールバックページ
- **Application層**: 
  - useGoogleLogin Hook（リダイレクト処理）
  - useAuthCallback Hook（コールバック処理）
- **Domain層**: ユーザー型定義
- **Infrastructure層**: APIクライアント（ky）、OpenAPI生成コード
- **状態管理**: 認証状態はJotaiで管理

## 備考

- この機能はmanagerサービスのフロントエンドとバックエンドの両方で実装が必要
- メールアドレス/パスワードによる登録機能と併用可能
- 同じメールアドレスで両方の認証方式を使用できる（google_idを追加）
- Googleで認証されたメールアドレスは自動的に検証済みとみなす（email_verified_at設定）
- Google OAuth 2.0の設定（Client ID、Client Secret）はGoogle Cloud Consoleで事前に取得する必要がある
- リダイレクトURIはGoogle Cloud Consoleで登録する必要がある
- 将来的には他のソーシャルログイン（GitHub、Microsoft等）も追加可能な設計とする
- データベーススキーマはsqldefでマイグレーション管理する
- SQLクエリはsqlcで型安全なGoコードを生成する
- APIエンドポイントは以下とする:
  - `/techcv/api/v1/auth/google/login` - Google認証開始
  - `/techcv/api/v1/auth/google/callback` - OAuth 2.0コールバック
- stateパラメータの管理にはセッションストア（Redis等）の使用を検討する
