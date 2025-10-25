# 全体概要

## 目的
- `.kiro/specs/manager/requirements/user_story/ゲストはソーシャルログイン（Google）によってユーザー登録できる.md` に基づき、ゲストが Google OAuth 2.0 を通じて manager サービスに登録・サインインできるようにする。
- `.kiro/specs/manager/design/google-social-login.md` のレイヤー構成方針を踏まえ、既存のクリーンアーキテクチャを崩さずに機能を追加する。

## スコープ
- バックエンド（Go）: DB スキーマ、ドメイン・ユースケース、インフラ実装、HTTP ハンドラ、設定とログ。
- フロントエンド（React）: ログイン UI、OAuth コールバック処理、API クライアント統合、エラーハンドリング。
- OpenAPI 仕様更新とコード生成、開発者向けドキュメント、インフラ設定（環境変数・マニフェスト）。

## 前提条件と依存関係
- Google Cloud Console で取得した OAuth クライアント ID/Secret とリダイレクト URI が用意されていること。
- OAuth state の永続化はインターフェース経由で抽象化し、まずはインメモリ実装、将来的には Redis 等の共有ストアを使用する。
- API は HTTP Only なセッションクッキーを発行し、`credentials: 'include'` を利用する現行フロントエンドと整合させる。
- 既存のメールアドレス・パスワード登録フローは維持し、Google ログインが追加の経路として共存する。
