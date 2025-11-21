# ゲストはソーシャルログイン（Google）によってユーザー登録できる

## ユーザーストーリー

**As a** ゲスト（未登録の利用者）  
**I want** Googleアカウントを使ってユーザー登録できる  
**So that** パスワードを設定せずに簡単にシステムにログインしてCVの管理機能を利用できる

## 概要

ゲストがmanagerサービスを利用するために、Googleアカウントを使用してユーザー登録を行う機能を提供します。Firebase Authenticationを使用してGoogleで認証を行い、登録が完了すると、ユーザーとしてシステムにログインし、CV管理機能にアクセスできるようになります。

## 受け入れ基準（Acceptance Criteria）

### 1. ソーシャルログインボタンの表示

**WHEN** ゲストがユーザー登録ページにアクセスする  
**THEN** システムは「Googleでログイン」ボタンを表示する

**WHEN** ゲストがログインページにアクセスする  
**THEN** システムは「Googleでログイン」ボタンを表示する

### 2. Google認証フローの開始

**WHEN** ゲストが「Googleでログイン」ボタンをクリックする  
**THEN** フロントエンドは以下の処理を実行する
- Firebase AuthのGoogle認証プロバイダーを使用して認証フローを開始する
- Googleの認証ページをポップアップまたはリダイレクトで表示する

### 3. Google認証ページの表示

**WHEN** ゲストがGoogleの認証ページにアクセスする  
**THEN** Googleは以下の情報を表示する
- アプリケーション名
- 要求される権限（メールアドレス、基本プロフィール情報）
- 許可/拒否のオプション

### 4. 認証の許可

**WHEN** ゲストがGoogleの認証ページで「許可」をクリックする  
**THEN** Firebase Authは以下の処理を実行する
- Googleから認証情報を取得する
- Firebase IDトークンを生成する
- フロントエンドに認証結果を返す

### 5. 認証の拒否

**WHEN** ゲストがGoogleの認証ページで「拒否」をクリックする  
**THEN** Firebase Authはエラーを返す

**WHEN** フロントエンドがエラーを受け取る  
**THEN** システムは「Google認証がキャンセルされました」というメッセージを表示する

### 6. Firebase IDトークンの取得

**WHEN** Firebase Authから認証成功の結果を受け取る  
**THEN** フロントエンドは以下の処理を実行する
- Firebase UserオブジェクトからIDトークンを取得する
- ユーザー情報（uid、email、displayName、photoURL）を取得する

**IF** IDトークンの取得に失敗する  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 7. バックエンドへのIDトークン送信

**WHEN** フロントエンドがFirebase IDトークンを取得する  
**THEN** フロントエンドは以下の処理を実行する
- バックエンドのユーザー登録/ログインAPIにIDトークンを送信する
- Firebase UID、メールアドレス、名前、プロフィール画像URLを含める

### 8. Firebase IDトークンの検証

**WHEN** バックエンドがIDトークンを受け取る  
**THEN** バックエンドは以下の処理を実行する
- Firebase Admin SDKを使用してIDトークンを検証する
- トークンの署名、有効期限、issuer、audienceを確認する
- トークンからFirebase UID、メールアドレス等を抽出する

**IF** IDトークンの検証に失敗する  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 9. 既存ユーザーの確認

**WHEN** バックエンドがFirebase UIDを取得する  
**THEN** バックエンドはFirebase UIDを使用して既存ユーザーを検索する

**IF** 既存ユーザーが見つかる  
**THEN** システムはログイン処理を実行する（新規登録ではなく）

### 10. 新規ユーザー登録の実行

**WHEN** 既存ユーザーが見つからない  
**THEN** バックエンドは以下の処理を実行する
- UUID v7形式でユーザーIDを生成する
- Firebaseから取得した情報を使用して新しいユーザーレコードを作成する
  - email: Firebaseから取得したメールアドレス
  - name: Firebaseから取得した名前
  - firebase_uid: Firebase UID
  - profile_image: Firebaseから取得したプロフィール画像URL（オプション）
- password_hashはNULLに設定する（ソーシャルログインのため）
- created_at、updated_atをUTCの現在時刻で記録する
- email_verified_atをUTCの現在時刻で記録する（Firebaseで検証済み）
- is_activeを1（有効）に設定する
- データベースに保存する

### 11. メールアドレスの重複チェック

**WHEN** バックエンドが新規ユーザーを登録しようとする  
**AND** Firebaseから取得したメールアドレスが既に別のユーザーで登録されている  
**THEN** バックエンドは以下の処理を実行する
- 既存のユーザーレコードにfirebase_uidを追加する
- ログイン処理を実行する

### 12. カスタムトークンの生成（オプション）

**WHEN** ユーザー登録またはログインが正常に完了する  
**THEN** バックエンドは以下のいずれかの処理を実行する
- オプションA: Firebase IDトークンをそのまま使用する
- オプションB: 独自のJWTトークンを生成してクライアントに返す
  - トークンにユーザーID、メールアドレス、発行日時、有効期限を含める

### 13. 登録成功時の処理

**WHEN** ユーザー登録が正常に完了する  
**THEN** フロントエンドは以下の処理を実行する
- ユーザーを自動的にログイン状態にする
- Firebase IDトークンをメモリまたはFirebase SDKに保持させる
- バックエンドから返されたユーザー情報を状態管理に保存する
- ダッシュボードページにリダイレクトする
- 「登録が完了しました」という成功メッセージを表示する

### 14. ログイン成功時の処理

**WHEN** 既存ユーザーのログインが正常に完了する  
**THEN** バックエンドは以下の処理を実行する
- last_login_atをUTCの現在時刻で更新する

**THEN** フロントエンドは以下の処理を実行する
- Firebase IDトークンをメモリまたはFirebase SDKに保持させる
- バックエンドから返されたユーザー情報を状態管理に保存する
- ダッシュボードページにリダイレクトする
- 「ログインしました」という成功メッセージを表示する

### 15. エラーハンドリング

**IF** Firebase Authとの通信中にネットワークエラーが発生する  
**THEN** システムは「ネットワークエラーが発生しました。再度お試しください」というエラーメッセージを表示する

**IF** データベースエラーやその他の予期しないエラーが発生する  
**THEN** システムは「登録処理中にエラーが発生しました。しばらくしてから再度お試しください」というエラーメッセージを表示する

### 16. セキュリティ要件

**WHEN** フロントエンドがFirebase Authを使用する  
**THEN** システムは以下のセキュリティ対策を実装する
- HTTPSを使用して通信を暗号化する
- Firebase SDKが自動的にCSRF対策を実施する
- IDトークンは自動的にリフレッシュされる

**WHEN** バックエンドがFirebase IDトークンを検証する  
**THEN** システムは以下のセキュリティ対策を実装する
- Firebase Admin SDKを使用してIDトークンの署名を検証する
- トークンの有効期限を検証する
- トークンのissuerとaudienceを検証する
- 検証済みのFirebase UIDのみを信頼する

### 17. レスポンシブデザイン

**WHEN** ゲストがモバイルデバイスから登録ページにアクセスする  
**THEN** システムはモバイル画面に最適化されたレイアウトで「Googleでログイン」ボタンを表示する

## 技術的な制約

### バックエンド技術スタック
- Firebase Admin SDK for Goを使用する（firebase.google.com/go/v4）
- Firebase IDトークンの検証にはAdmin SDKのauth.VerifyIDToken()を使用する
- Firebase Project IDは環境変数で管理する
- Firebase Service Account認証情報（JSON）は環境変数またはファイルで管理する
- 日時はすべてUTCで保存し、DATETIME(6)型を使用する（マイクロ秒精度）

### データベーススキーマ要件

**usersテーブルの拡張**:
- `firebase_uid` VARCHAR(128) UNIQUE - Firebase UID
- `profile_image` VARCHAR(500) - プロフィール画像URL
- `password_hash` VARCHAR(255) - ソーシャルログインの場合はNULL許可に変更
- INDEX on `firebase_uid` for fast lookup

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
- Firebase JavaScript SDK v10+（firebase/auth）
- TanStack Router for routing
- TanStack Query for data fetching
- Jotai for global state management
- ky for HTTP client
- shadcn/ui for UI components
- Firebase Authenticationの設定（Web API Key、Auth Domain）は環境変数で管理

## 非機能要件

- Firebase認証フロー全体は30秒以内に完了する
- バックエンドのユーザー登録処理は3秒以内に完了する
- 「Googleでログイン」ボタンはGoogleのブランドガイドラインに準拠する
- すべてのエラーは明確で理解しやすいメッセージで表示する
- APIエラーレスポンスは統一されたフォーマット（requestId、code、details）で返す
- ログはCloudLogging形式のJSON構造化ログとして出力する
- ログにはcontext由来のrequest_idを自動付与する
- Firebase認証フローの主要なステップをログに記録する（デバッグ用）
- Firebase IDトークンのリフレッシュは自動的に行われる（Firebase SDKが管理）

## 関連するユビキタス言語

- **ゲスト（guest）**: まだユーザー登録していない利用者
- **ユーザー（user）**: ユーザー登録が済んだ利用者
- **manager**: WebエンジニアのCVを管理するサービス
- **ソーシャルログイン（social login）**: 外部サービス（Google等）のアカウントを使用した認証方式
- **Firebase Authentication**: Googleが提供する認証サービス
- **Firebase UID**: Firebaseが発行するユーザーの一意識別子
- **Firebase IDトークン（ID token）**: Firebase Authが発行するJWT形式の認証トークン

## アーキテクチャ上の考慮事項

### バックエンド（Clean Architecture + DDD + CQRS）

- **Domain層**: User集約、Email値オブジェクト、FirebaseUID値オブジェクトを定義
- **UseCase層**: RegisterUserWithFirebaseコマンド、LoginWithFirebaseコマンドを実装
- **Adapter層**: Firebase認証ハンドラー、HTTPハンドラー
- **Infrastructure層**: 
  - sqlcを使用したリポジトリ実装
  - Firebase Admin SDK実装（IDトークン検証）
- **CQRS**: ユーザー登録はコマンド側で実装（集約を使用）
- **トランザクション**: User集約の保存は単一トランザクションで実行
- **バリデーション**: Firebase IDトークンの検証、ユーザー情報の検証

### フロントエンド（レイヤードアーキテクチャ）

- **Presentation層**: 
  - 登録/ログインページコンポーネント
  - Googleログインボタンコンポーネント
- **Application層**: 
  - useFirebaseAuth Hook（Firebase認証処理）
  - useRegisterWithFirebase Hook（バックエンド連携）
- **Domain層**: ユーザー型定義
- **Infrastructure層**: 
  - Firebase JavaScript SDK（認証処理）
  - APIクライアント（ky）でバックエンドと通信
  - OpenAPI生成コード
- **状態管理**: 
  - Firebase認証状態はFirebase SDKが管理
  - アプリケーション固有のユーザー情報はJotaiで管理

## 備考

- この機能はmanagerサービスのフロントエンドとバックエンドの両方で実装が必要
- メールアドレス/パスワードによる登録機能と併用可能
- 同じメールアドレスで両方の認証方式を使用できる（firebase_uidを追加）
- Firebaseで認証されたメールアドレスは自動的に検証済みとみなす（email_verified_at設定）
- Firebase Authenticationの設定はFirebase Consoleで事前に行う必要がある
  - Google認証プロバイダーを有効化
  - 承認済みドメインを登録
- 将来的には他のソーシャルログイン（GitHub、Microsoft等）も追加可能な設計とする
  - Firebase Authは複数の認証プロバイダーをサポート
- データベーススキーマはsqldefでマイグレーション管理する
- SQLクエリはsqlcで型安全なGoコードを生成する
- APIエンドポイントは以下とする:
  - `/techcv/api/v1/auth/firebase/register` - Firebase認証後のユーザー登録
  - `/techcv/api/v1/auth/firebase/login` - Firebase認証後のログイン
- Firebase IDトークンはHTTPヘッダー（Authorization: Bearer <token>）で送信する
- Firebase Admin SDKの初期化にはService Account認証情報が必要
  - 本番環境ではGoogle Cloud Secret Managerの使用を推奨
- Firebase IDトークンの有効期限は1時間（自動リフレッシュ）
- バックエンドでは毎回Firebase IDトークンを検証する必要がある
