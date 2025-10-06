# Implementation Plan

## Backend Implementation

- [ ] 1. プロジェクト構造とインフラ基盤の構築
  - Clean Architectureに基づいたディレクトリ構造を作成
  - Terraformでインフラ定義（VPC、Cloud SQL、Cloud Run、IAM）
  - GitHub Actionsワークフロー設定
  - _Requirements: 全体_

- [ ] 1.1 Goプロジェクトの初期化
  - `go mod init`でプロジェクト初期化
  - 必要な依存パッケージをインストール（Echo、sqlc、excelize、Connect等）
  - Clean Architectureディレクトリ構造作成（domain, usecase, interface, infrastructure）
  - _Requirements: 全体_

- [ ] 1.2 Terraformインフラ定義
  - VPCとサブネット設定
  - Cloud SQL for MySQL設定（プライベートIP）
  - VPC Connector設定
  - Cloud Run設定（Backend: internal, Frontend: public）
  - IAM・Service Account設定
  - Secret Manager設定
  - _Requirements: 全体_

- [ ] 1.3 GitHub Actionsワークフロー作成
  - Backend デプロイワークフロー
  - Frontend デプロイワークフロー
  - Terraform適用ワークフロー
  - Workload Identity Federation設定
  - _Requirements: 全体_

- [ ] 2. データベーススキーマとマイグレーション
  - schema.sqlファイル作成
  - sqldefでマイグレーション実行
  - sqlc設定とクエリ定義
  - _Requirements: 1, 2, 3_

- [ ] 2.1 schema.sql作成
  - users, cvs, basic_infos, work_experiences, skills, educations, certifications, projectsテーブル定義
  - インデックスと外部キー制約設定
  - _Requirements: 1, 2, 3_

- [ ] 2.2 sqlc設定とクエリ定義
  - sqlc.yaml設定ファイル作成
  - queries/ディレクトリにSQLクエリ定義
  - `sqlc generate`で型安全なGoコード生成
  - _Requirements: 1, 2, 3_

- [ ] 3. Domain層の実装
  - Entityとリポジトリインターフェース定義
  - Domain Serviceの実装
  - _Requirements: 1, 2, 3_

- [ ] 3.1 User Domain実装
  - User Entity定義
  - UserRole定義
  - User Repository Interface定義
  - _Requirements: 1_

- [ ] 3.2 CV Domain実装
  - CV, BasicInfo, WorkExperience, Skill, Education, Certification, Project Entity定義
  - CV Repository Interface定義
  - CV Domain Service（公開フィールドフィルタリング）実装
  - _Requirements: 2, 3_

- [ ] 4. Infrastructure層の実装
  - Repositoryの実装（sqlc使用）
  - 外部サービス連携（Google OAuth）
  - _Requirements: 1, 2_

- [ ] 4.1 User Repository実装
  - sqlcで生成されたコードを使用してUser Repository実装
  - FindByID, FindByGoogleID, Create, Update実装
  - _Requirements: 1_

- [ ] 4.2 CV Repository実装
  - sqlcで生成されたコードを使用してCV Repository実装
  - FindByID, FindByUserID, Create, Update, Delete実装
  - 関連エンティティ（BasicInfo, WorkExperience等）の操作実装
  - _Requirements: 2, 3_

- [ ] 4.3 Google OAuth Client実装
  - Google ID Token検証機能実装
  - Google User Info取得機能実装
  - _Requirements: 1_

- [ ] 5. Use Case層の実装
  - 認証、ユーザー管理、CV管理、Excel出力のUse Case実装
  - _Requirements: 1, 2, 3, 4, 5_

- [ ] 5.1 Auth Use Case実装
  - VerifyGoogleToken Use Case実装
  - GetOrCreateUser Use Case実装
  - JWT生成機能実装
  - _Requirements: 1_

- [ ] 5.2 User Use Case実装
  - GetUsers, GetUserByID, UpdateUser Use Case実装
  - 権限チェック機能実装
  - _Requirements: 1_

- [ ] 5.3 CV Use Case実装
  - GetCVList, GetCVByID, CreateCV, UpdateCV, DeleteCV Use Case実装
  - 権限に基づくアクセス制御実装
  - 公開フィールドフィルタリング適用
  - _Requirements: 2, 3, 4_

- [ ] 5.4 Excel Export Use Case実装
  - ExportSingleCV, ExportMultipleCVs Use Case実装
  - excelize使用したExcel生成機能実装
  - 公開フィールドのみをExcelに含める処理実装
  - _Requirements: 5_

- [ ] 6. Interface層の実装（Handlers & Middleware）
  - Echo HTTPハンドラー実装
  - Connect gRPCハンドラー実装
  - 認証ミドルウェア実装
  - _Requirements: 1, 2, 3, 4, 5_

- [ ] 6.1 Protocol Buffers定義
  - auth.proto, user.proto, cv.proto, export.proto定義
  - `protoc`でGoコード生成
  - _Requirements: 1, 2, 3, 5_

- [ ] 6.2 Auth Handler実装
  - Google OAuth コールバックハンドラー（Echo）
  - VerifyGoogleToken gRPCハンドラー（Connect）
  - GetMe gRPCハンドラー（Connect）
  - _Requirements: 1_

- [ ] 6.3 User Handler実装
  - ListUsers, GetUser, UpdateUser gRPCハンドラー（Connect）
  - 権限チェック実装
  - _Requirements: 1_

- [ ] 6.4 CV Handler実装
  - ListCVs, GetCV, CreateCV, UpdateCV, DeleteCV gRPCハンドラー（Connect）
  - 権限チェック実装
  - _Requirements: 2, 3, 4_

- [ ] 6.5 Excel Export Handler実装
  - ExportCV, ExportMultipleCVs RESTハンドラー（Echo）
  - ファイルダウンロードレスポンス実装
  - _Requirements: 5_

- [ ] 6.6 認証ミドルウェア実装
  - JWT検証ミドルウェア実装
  - ユーザー情報をコンテキストに設定
  - _Requirements: 1_

- [ ] 7. サーバー起動とルーティング設定
  - Echoサーバー初期化
  - Connect handlerをEchoに統合
  - ルーティング設定
  - _Requirements: 全体_

- [ ] 7.1 メインサーバー実装
  - cmd/api/main.go作成
  - 依存性注入（DI）設定
  - Echoサーバー起動
  - Connect handlerマウント
  - ヘルスチェックエンドポイント実装
  - _Requirements: 全体_

- [ ] 8. Dockerfile作成とコンテナ化
  - マルチステージビルドDockerfile作成
  - イメージサイズ最適化
  - _Requirements: 全体_

## Frontend Implementation

- [ ] 9. Next.jsプロジェクトセットアップ
  - Next.js 14+ (App Router) プロジェクト初期化
  - Material-UI (MUI) セットアップ
  - TanStack Query セットアップ
  - NextAuth.js セットアップ
  - _Requirements: 全体_

- [ ] 9.1 プロジェクト初期化と依存関係インストール
  - `create-next-app`でプロジェクト作成
  - MUI, TanStack Query, NextAuth.js, React Hook Form, Zodインストール
  - Clean Architectureディレクトリ構造作成
  - _Requirements: 全体_

- [ ] 9.2 NextAuth.js設定
  - Google OAuth Provider設定
  - セッション管理設定
  - 認証コールバック実装
  - _Requirements: 1_

- [ ] 9.3 Connect Client設定
  - Protocol Buffers定義（Backendと共有）
  - Connect client生成
  - TanStack Queryとの統合
  - _Requirements: 全体_

- [ ] 10. 共通コンポーネントとレイアウト
  - レイアウトコンポーネント作成
  - ナビゲーション実装
  - 認証状態管理
  - _Requirements: 全体_

- [ ] 10.1 レイアウトコンポーネント実装
  - AppBar, Drawer, Footer実装（MUI）
  - レスポンシブデザイン対応
  - _Requirements: 全体_

- [ ] 10.2 認証コンポーネント実装
  - GoogleSignInButton実装
  - AuthProvider実装
  - ProtectedRoute実装
  - _Requirements: 1_

- [ ] 11. Auth Feature実装
  - Domain, Application, Infrastructure, Presentation層実装
  - _Requirements: 1_

- [ ] 11.1 Auth Domain & Application層
  - User Entity定義
  - Auth Use Case Interface定義
  - _Requirements: 1_

- [ ] 11.2 Auth Infrastructure層
  - Connect client使用したAPI通信実装
  - TanStack Query hooks実装（useAuth, useLogin等）
  - _Requirements: 1_

- [ ] 11.3 Auth Presentation層
  - ログインページ実装
  - ログアウト機能実装
  - _Requirements: 1_

- [ ] 12. CV Feature実装（一覧・詳細表示）
  - Domain, Application, Infrastructure, Presentation層実装
  - _Requirements: 2, 3, 4_

- [ ] 12.1 CV Domain & Application層
  - CV Entity定義
  - CV Use Case Interface定義
  - _Requirements: 2, 3, 4_

- [ ] 12.2 CV Infrastructure層
  - Connect client使用したAPI通信実装
  - TanStack Query hooks実装（useCVList, useCV等）
  - _Requirements: 2, 3, 4_

- [ ] 12.3 CV一覧画面実装
  - CV一覧表示コンポーネント実装（MUI Table）
  - 権限に基づく表示制御
  - _Requirements: 4_

- [ ] 12.4 CV詳細画面実装
  - CV詳細表示コンポーネント実装（MUI Card, Typography等）
  - 公開設定に基づく条件付き表示
  - レスポンシブデザイン対応
  - _Requirements: 3, 4_

- [ ] 13. CV Feature実装（登録・編集）
  - CV登録・編集フォーム実装
  - _Requirements: 2, 3_

- [ ] 13.1 CV登録・編集フォーム実装
  - React Hook Form + Zod使用したフォーム実装
  - 基本情報セクション実装
  - 職務経歴セクション実装（動的追加・削除）
  - スキルセクション実装（動的追加・削除）
  - 学歴セクション実装（動的追加・削除）
  - 資格セクション実装（動的追加・削除）
  - プロジェクトセクション実装（動的追加・削除）
  - _Requirements: 2_

- [ ] 13.2 公開/非公開切り替え機能実装
  - 各項目に公開/非公開チェックボックス追加
  - 公開設定の保存機能実装
  - _Requirements: 3_

- [ ] 13.3 バリデーションとエラーハンドリング
  - Zodスキーマ定義
  - バリデーションエラー表示
  - API エラーハンドリング
  - _Requirements: 2_

- [ ] 14. Excel出力機能実装
  - Excel出力ボタンとダウンロード機能実装
  - _Requirements: 5_

- [ ] 14.1 Excel出力機能実装
  - CV詳細画面にExcel出力ボタン追加
  - Backend APIを呼び出してExcelファイルダウンロード
  - 複数CV一括出力機能実装（システムオーナーのみ）
  - _Requirements: 5_

- [ ] 15. Dockerfile作成とコンテナ化
  - Next.js用Dockerfile作成
  - 環境変数設定
  - イメージサイズ最適化
  - _Requirements: 全体_

## Integration & Deployment

- [ ] 16. ローカル開発環境構築
  - Docker Composeでローカル環境構築
  - Backend, Frontend, MySQLコンテナ設定
  - _Requirements: 全体_

- [ ] 17. Terraformでインフラ構築
  - `terraform apply`でGCPリソース作成
  - Cloud SQL, Cloud Run, VPC Connector等のプロビジョニング
  - _Requirements: 全体_

- [ ] 18. GitHub Actionsでデプロイ
  - Backend デプロイ実行
  - Frontend デプロイ実行
  - 動作確認
  - _Requirements: 全体_

- [ ] 19. エンドツーエンドテスト
  - 主要なユーザーフローのテスト
  - Google OAuth認証フロー確認
  - CV登録・編集・表示フロー確認
  - Excel出力フロー確認
  - _Requirements: 全体_
