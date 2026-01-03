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

### 6. Firebase IDトークンとユーザー情報の取得

**WHEN** Firebase Authから認証成功の結果を受け取る  
**THEN** フロントエンドは以下の処理を実行する
- Firebase UserオブジェクトからIDトークンを取得する
- Firebase Userオブジェクトから以下の全ての情報を取得する
  - uid: Firebase UID
  - email: メールアドレス
  - emailVerified: メール確認済みフラグ
  - displayName: 表示名
  - photoURL: プロフィール画像URL
  - phoneNumber: 電話番号（Googleアカウントに設定されている場合）
  - metadata.creationTime: Firebase上のアカウント作成日時
  - metadata.lastSignInTime: 最終サインイン日時
  - providerData: 認証プロバイダー情報（Google）

**IF** IDトークンの取得に失敗する  
**THEN** システムは「認証に失敗しました。再度お試しください」というエラーメッセージを表示する

### 7. バックエンドへのIDトークンとユーザー情報の送信

**WHEN** フロントエンドがFirebase IDトークンとユーザー情報を取得する  
**THEN** フロントエンドは以下の処理を実行する
- バックエンドのユーザー登録/ログインAPIにIDトークンを送信する
- Firebase Userオブジェクトから取得した全ての情報を含める
  - uid
  - email
  - emailVerified
  - displayName
  - photoURL
  - phoneNumber
  - metadata（creationTime、lastSignInTime）
  - providerData

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
- Firebaseから取得した全ての情報を使用して新しいユーザーレコードを作成する
  - firebase_uid: Firebase UID
  - email: メールアドレス
  - email_verified: メール確認済みフラグ（Firebaseから取得）
  - display_name: 表示名
  - photo_url: プロフィール画像URL
  - phone_number: 電話番号
  - firebase_created_at: Firebase上のアカウント作成日時
  - firebase_last_sign_in_at: Firebase上の最終サインイン日時
  - provider_id: 認証プロバイダーID（google.com）
- created_at、updated_atをUTCの現在時刻で記録する
- email_verified_atをFirebaseのemailVerifiedフラグに基づいて設定する
- is_activeを1（有効）に設定する
- データベースに保存する

### 11. メールアドレスの重複チェック

**WHEN** バックエンドが新規ユーザーを登録しようとする  
**AND** Firebaseから取得したメールアドレスが既に別のユーザーで登録されている  
**THEN** バックエンドは以下の処理を実行する
- 既存のユーザーレコードにfirebase_uidを追加する
- ログイン処理を実行する

### 12. 認証トークンの管理

**WHEN** ユーザー登録またはログインが正常に完了する  
**THEN** システムは以下の処理を実行する
- Firebase IDトークンをそのまま認証トークンとして使用する
- 独自のJWTトークンは生成しない
- バックエンドはFirebase IDトークンを検証することでユーザーを認証する

### 13. 登録成功時の処理

**WHEN** ユーザー登録が正常に完了する  
**THEN** フロントエンドは以下の処理を実行する
- ユーザーを自動的にログイン状態にする
- Firebase IDトークンはFirebase SDKが自動的に管理する
- バックエンドへのAPIリクエスト時は、Firebase SDKから最新のIDトークンを取得してAuthorizationヘッダーに設定する
- バックエンドから返されたユーザー情報を状態管理に保存する
- ダッシュボードページにリダイレクトする
- 「登録が完了しました」という成功メッセージを表示する

### 14. ログイン成功時の処理

**WHEN** 既存ユーザーのログインが正常に完了する  
**THEN** バックエンドは以下の処理を実行する
- last_login_atをUTCの現在時刻で更新する

**THEN** フロントエンドは以下の処理を実行する
- Firebase IDトークンはFirebase SDKが自動的に管理する
- バックエンドへのAPIリクエスト時は、Firebase SDKから最新のIDトークンを取得してAuthorizationヘッダーに設定する
- バックエンドから返されたユーザー情報を状態管理に保存する
- ダッシュボードページにリダイレクトする
- 「ログインしました」という成功メッセージを表示する

### 15. 認証が必要なAPIリクエスト

**WHEN** フロントエンドが認証が必要なAPIをリクエストする  
**THEN** フロントエンドは以下の処理を実行する
- Firebase SDKのcurrentUser.getIdToken()を呼び出して最新のIDトークンを取得する
- IDトークンをAuthorizationヘッダー（Bearer <token>）に設定する
- APIリクエストを送信する

**WHEN** バックエンドが認証が必要なAPIリクエストを受け取る  
**THEN** バックエンドは以下の処理を実行する
- AuthorizationヘッダーからFirebase IDトークンを抽出する
- Firebase Admin SDKでIDトークンを検証する
- 検証成功後、トークンからFirebase UIDを取得する
- Firebase UIDを使用してユーザー情報をデータベースから取得する
- リクエスト処理を続行する

**IF** IDトークンの検証に失敗する  
**THEN** バックエンドは401 Unauthorizedエラーを返す

### 16. IDトークンのリフレッシュ

**WHEN** Firebase IDトークンの有効期限が近づく  
**THEN** Firebase SDKは以下の処理を自動的に実行する
- バックグラウンドでIDトークンをリフレッシュする
- 新しいIDトークンを取得する
- アプリケーションコードは何もする必要がない

**WHEN** フロントエンドがgetIdToken()を呼び出す  
**THEN** Firebase SDKは以下の処理を実行する
- 現在のIDトークンが有効であればそれを返す
- 有効期限が切れている場合は自動的にリフレッシュしてから返す

### 17. エラーハンドリング

**IF** Firebase Authとの通信中にネットワークエラーが発生する  
**THEN** システムは「ネットワークエラーが発生しました。再度お試しください」というエラーメッセージを表示する

**IF** データベースエラーやその他の予期しないエラーが発生する  
**THEN** システムは「登録処理中にエラーが発生しました。しばらくしてから再度お試しください」というエラーメッセージを表示する

**IF** Firebase IDトークンの有効期限が切れている、または無効である  
**THEN** バックエンドは401 Unauthorizedエラーを返す

**WHEN** フロントエンドが401エラーを受け取る  
**THEN** フロントエンドは以下の処理を実行する
- ユーザーをログアウトさせる
- ログインページにリダイレクトする
- 「セッションの有効期限が切れました。再度ログインしてください」というメッセージを表示する

### 18. セキュリティ要件

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

### 19. レスポンシブデザイン

**WHEN** ゲストがモバイルデバイスから登録ページにアクセスする  
**THEN** システムはモバイル画面に最適化されたレイアウトで「Googleでログイン」ボタンを表示する

## 技術的な制約

### バックエンド技術スタック
- プログラミング言語: Golang 1.25
- Webフレームワーク: Echo
- Firebase Admin SDK for Goを使用する（firebase.google.com/go/v4）
- Firebase IDトークンの検証にはAdmin SDKのauth.VerifyIDToken()を使用する
- Firebase Project IDは環境変数で管理する（envconfigを使用）
- Firebase Service Account認証情報（JSON）は環境変数またはファイルで管理する
- ローカル環境変数の読み込みにはgodotenvを使用する
- ログ出力にはslog（Go標準ライブラリ）を使用する
- ID生成にはgoogle/uuidのUUID v7を使用する
- データベースアクセスにはsqlcを使用する
- データベースマイグレーションにはsqldefを使用する
- OpenAPI仕様はOpenAPI 3.0.3を使用する
- OpenAPIコード生成にはoapi-codegenを使用する
- 日時はすべてUTCで保存し、DATETIME(6)型を使用する（マイクロ秒精度）
- 時刻の入出力はISO8601拡張形式、UTC、ミリ秒精度を使用する

### データベーススキーマ要件

**usersテーブル**:
```sql
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY COMMENT 'UUID v7（アプリケーション内部ID）',
    firebase_uid VARCHAR(128) NOT NULL UNIQUE COMMENT 'Firebase UID',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'メールアドレス',
    email_verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'メール確認済みフラグ（Firebaseから取得）',
    display_name VARCHAR(255) COMMENT '表示名（Firebaseから取得）',
    photo_url VARCHAR(500) COMMENT 'プロフィール画像URL（Firebaseから取得）',
    phone_number VARCHAR(50) COMMENT '電話番号（Firebaseから取得、オプション）',
    provider_id VARCHAR(50) NOT NULL COMMENT '認証プロバイダーID（例: google.com）',
    firebase_created_at DATETIME(6) COMMENT 'Firebase上のアカウント作成日時',
    firebase_last_sign_in_at DATETIME(6) COMMENT 'Firebase上の最終サインイン日時',
    bio TEXT COMMENT '自己紹介（アプリケーション独自項目）',
    is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'アクティブ状態',
    email_verified_at DATETIME(6) COMMENT 'メール確認日時（アプリケーション側で管理）',
    last_login_at DATETIME(6) COMMENT '最終ログイン日時（アプリケーション側で管理）',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
    deleted_at DATETIME(6) COMMENT '削除日時（論理削除用）',
    
    UNIQUE KEY uq_users_email (email),
    UNIQUE KEY uq_users_firebase_uid (firebase_uid),
    INDEX idx_users_firebase_uid (firebase_uid),
    INDEX idx_users_provider_id (provider_id),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー情報';
```

### フロントエンド技術スタック
- フレームワーク: React 18+ with TypeScript
- Firebase JavaScript SDK v10+（firebase/auth）
- UIライブラリ: shadcn/ui
- 状態管理: Jotai
- ルーティング: TanStack Router
- HTTPクライアント: TanStack Query
- OpenAPIコード生成: OpenAPI Generator
- Firebase Authenticationの設定（Web API Key、Auth Domain）は環境変数で管理

## 非機能要件

- Firebase認証フロー全体は30秒以内に完了する
- バックエンドのユーザー登録処理は3秒以内に完了する
- Firebase IDトークンの検証は100ms以内に完了する
- 「Googleでログイン」ボタンはGoogleのブランドガイドラインに準拠する
- すべてのエラーは明確で理解しやすいメッセージで表示する
- APIエラーレスポンスは統一されたフォーマットで返す
  ```json
  {
    "requestId": "88374925",
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "email",
        "code": "INVALID_EMAIL_FORMAT"
      }
    ]
  }
  ```
  - requestId: 必須 - リクエストをユニークに特定するためのランダムID
  - code: omitempty - エラーコード（大文字のスネークケース）
  - details: omitempty - 詳細エラー情報の配列（主にバリデーションエラー用）
- ログはCloudLogging形式のJSON構造化ログとして出力する
- ログにはcontext由来のrequest_idを自動付与する
- slogでのログ出力はContext付きログ関数を使用する
- Firebase認証フローの主要なステップをログに記録する（デバッグ用）
- Firebase IDトークンのリフレッシュは自動的に行われる（Firebase SDKが管理）
- 認証が必要な全てのAPIエンドポイントでFirebase IDトークンを検証する
- IDトークンの検証失敗時は適切なエラーレスポンスを返す（401 Unauthorized）

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
  - TanStack Query（HTTPクライアント）でバックエンドと通信
  - OpenAPI Generatorで生成されたコード
- **状態管理**: 
  - Firebase認証状態はFirebase SDKが管理
  - アプリケーション固有のユーザー情報はJotaiで管理

## 備考

- この機能はmanagerサービスのフロントエンドとバックエンドの両方で実装が必要
- 認証方式はソーシャルログイン（Firebase Authentication）のみを使用
- パスワードによる認証は実装しない（password_hashカラムは不要）
- Firebaseで認証されたメールアドレスは自動的に検証済みとみなす（email_verified_at設定）

### Firebase Authenticationの設定
- Firebase Consoleで事前に設定が必要
  - Google認証プロバイダーを有効化
  - 承認済みドメインを登録
- Firebase Admin SDKの初期化にはService Account認証情報が必要
  - 本番環境ではGoogle Cloud Secret Managerの使用を推奨

### 認証トークンの管理
- Firebase IDトークンを認証トークンとして使用する
- 独自のJWTトークンは発行しない
- Firebase IDトークンはHTTPヘッダー（Authorization: Bearer <token>）で送信する
- Firebase IDトークンの有効期限は1時間（Firebase SDKが自動リフレッシュ）
- バックエンドでは認証が必要な全てのAPIリクエストでFirebase IDトークンを検証する
- フロントエンドはFirebase SDKのgetIdToken()を使用して常に最新のトークンを取得する

### データベース設計
- データベース: MySQL 8.0+
- データベーススキーマはsqldefでマイグレーション管理する
- SQLクエリはsqlcで型安全なGoコードを生成する
- sqlcによって生成されたファイルは直接修正しない
- ユーザーの一意性はfirebase_uidで保証される
- 文字コード: utf8mb4、照合順序: utf8mb4_unicode_ci

### APIエンドポイント
- エンドポイントベース: `/techcv/api/v1`
- `/techcv/api/v1/auth/firebase/register` - Firebase認証後のユーザー登録
- `/techcv/api/v1/auth/firebase/login` - Firebase認証後のログイン
- 上記以外の認証が必要なAPIは全てAuthorizationヘッダーでFirebase IDトークンを要求する
- REST APIを採用する（リソースベースで扱いが難しい場合は原則を崩すことも認める）

### 将来の拡張性
- 他のソーシャルログイン（GitHub、Microsoft等）も追加可能な設計
  - Firebase Authは複数の認証プロバイダーをサポート
  - データベーススキーマは変更不要（firebase_uidで統一）
  - provider_idカラムで認証プロバイダーを識別
