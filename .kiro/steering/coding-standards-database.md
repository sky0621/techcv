# Databaseに関するコーディング規約

## 採用Database
MySQL v8.0+

## 基本原則

- データの整合性を最優先する
- パフォーマンスを考慮した設計を行う
- 正規化を基本とするが、パフォーマンスのために非正規化も検討する
- マイグレーションで全ての変更を管理する
- インデックスを適切に設計する

## 命名規則

### テーブル名
- 小文字のスネークケースを使用
- 複数形を使用（例: `users`, `orders`, `order_items`）
- 中間テーブルは関連する2つのテーブル名を結合（例: `user_roles`）
- プレフィックスは使用しない

```sql
-- 良い例
users
orders
order_items
user_roles

-- 悪い例
User
tbl_users
user
```

### カラム名
- 小文字のスネークケースを使用
- 明確で説明的な名前を使用
- 省略形は避ける（一般的なものを除く）
- boolean型は `is_`, `has_`, `can_` などのプレフィックスを使用

```sql
-- 良い例
user_id
email_address
is_active
has_verified_email
created_at

-- 悪い例
userId
email_addr
active
verified
create_date
```

### 主キー
- `id` を使用（テーブル名のプレフィックスは不要）
- UUID v7を使用（BINARY(16)で格納）
- AUTO_INCREMENTは使用しない

```sql
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);
```

### 外部キー
- `{参照先テーブル名}_id` の形式を使用
- 外部キー制約を必ず設定

```sql
CREATE TABLE orders (
    id BINARY(16) PRIMARY KEY,
    user_id BINARY(16) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### インデックス名
- `idx_{テーブル名}_{カラム名}` の形式を使用
- 複合インデックスは `idx_{テーブル名}_{カラム1}_{カラム2}` の形式

```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id_created_at ON orders(user_id, created_at);
```

### ユニーク制約名
- `uq_{テーブル名}_{カラム名}` の形式を使用

```sql
ALTER TABLE users ADD CONSTRAINT uq_users_email UNIQUE (email);
```

## データ型

### 文字列
- `VARCHAR(n)` を使用（固定長が必要な場合のみ `CHAR(n)`）
- メールアドレス: `VARCHAR(255)`
- 名前: `VARCHAR(100)`
- 説明文: `TEXT`
- 長い文章: `MEDIUMTEXT` または `LONGTEXT`

```sql
email VARCHAR(255) NOT NULL,
name VARCHAR(100) NOT NULL,
bio TEXT,
content MEDIUMTEXT
```

### 数値
- 整数: `INT`, `BIGINT`
- 小数: `DECIMAL(p, s)` （金額など正確な計算が必要な場合）
- 浮動小数点: `DOUBLE` （近似値で問題ない場合）

```sql
age INT,
price DECIMAL(10, 2),  -- 最大99,999,999.99
latitude DOUBLE,
longitude DOUBLE
```

### 日時
- `DATETIME(6)` を使用（マイクロ秒まで記録）
- タイムゾーンはアプリケーション層で管理
- `TIMESTAMP` は使用しない（2038年問題）

```sql
created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
deleted_at DATETIME(6)  -- 論理削除用
```

### 真偽値
- `TINYINT(1)` を使用（0 = false, 1 = true）
- `BOOLEAN` は `TINYINT(1)` のエイリアス

```sql
is_active TINYINT(1) NOT NULL DEFAULT 1,
has_verified_email TINYINT(1) NOT NULL DEFAULT 0
```

### UUID
- `BINARY(16)` で格納
- アプリケーション層でUUID v7を生成

```sql
id BINARY(16) PRIMARY KEY,
user_id BINARY(16) NOT NULL
```

### JSON
- `JSON` 型を使用
- 構造化されたデータで検索が必要な場合に使用
- 頻繁に検索するフィールドは別カラムに抽出を検討

```sql
metadata JSON,
settings JSON
```

### ENUM
- 使用を避ける（変更が困難）
- 代わりに参照テーブルを作成

```sql
-- 避ける
status ENUM('draft', 'published', 'archived')
```

## テーブル設計

### 必須カラム
全てのテーブルに以下のカラムを含める：

```sql
CREATE TABLE table_name (
    id BINARY(16) PRIMARY KEY,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
);
```

### 論理削除
論理削除が必要な場合は `deleted_at` を追加：

```sql
deleted_at DATETIME(6),
INDEX idx_table_name_deleted_at (deleted_at)
```

### 楽観的ロック
同時更新制御が必要な場合は `version` を追加：

```sql
version INT NOT NULL DEFAULT 1
```

### テーブル設計例

```sql
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    bio TEXT,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    email_verified_at DATETIME(6),
    last_login_at DATETIME(6),
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at DATETIME(6),
    
    UNIQUE KEY uq_users_email (email),
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## インデックス設計

### インデックスを作成すべき場合
- 主キー（自動的に作成される）
- 外部キー
- WHERE句で頻繁に使用されるカラム
- JOIN条件で使用されるカラム
- ORDER BY句で使用されるカラム
- ユニーク制約が必要なカラム

### 複合インデックス
- カーディナリティの高いカラムを先に配置
- WHERE句とORDER BY句の両方で使用される場合を考慮
- インデックスの順序が重要（左端一致の原則）

```sql
-- user_idで検索し、created_atでソートする場合
CREATE INDEX idx_orders_user_id_created_at ON orders(user_id, created_at);

-- この順序では効率的
SELECT * FROM orders WHERE user_id = ? ORDER BY created_at DESC;

-- この順序では非効率（user_idがない）
SELECT * FROM orders WHERE created_at > ? ORDER BY created_at DESC;
```

### インデックスを避けるべき場合
- カーディナリティが低いカラム（性別など）
- 頻繁に更新されるカラム
- 小さなテーブル（数千行以下）

### カバリングインデックス
- SELECT句のカラムを全てインデックスに含める（パフォーマンス最適化）

```sql
-- user_idとstatusで検索し、created_atも取得する場合
CREATE INDEX idx_orders_user_id_status_created_at 
ON orders(user_id, status, created_at);
```

## 制約

### NOT NULL制約
- 必須項目には必ず `NOT NULL` を指定
- NULLを許可する場合は明示的に設計判断を行う

```sql
email VARCHAR(255) NOT NULL,
bio TEXT  -- NULLを許可
```

### UNIQUE制約
- 一意性が必要なカラムには `UNIQUE` を指定
- 複合ユニーク制約も活用

```sql
email VARCHAR(255) NOT NULL UNIQUE,

-- 複合ユニーク制約
UNIQUE KEY uq_user_roles_user_id_role_id (user_id, role_id)
```

### 外部キー制約
- 参照整合性を保証するために必ず設定
- ON DELETE / ON UPDATE の動作を明示的に指定

```sql
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
```

### CHECK制約（MySQL 8.0.16+）
- 値の範囲を制限する場合に使用

```sql
age INT CHECK (age >= 0 AND age <= 150),
price DECIMAL(10, 2) CHECK (price >= 0)
```

## マイグレーション

### 基本ルール
- 全てのスキーマ変更はマイグレーションで管理
- マイグレーションファイルは決して変更しない
- 本番環境に適用済みのマイグレーションは削除しない
- ロールバック可能な設計を心がける

### マイグレーションファイル名
- `YYYYMMDDHHMMSS_description.sql` の形式
- 説明は英語で簡潔に

```
20240115120000_create_users_table.sql
20240115130000_add_email_verified_at_to_users.sql
20240115140000_create_orders_table.sql
```

### テーブル作成

```sql
-- 20240115120000_create_users_table.sql
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    
    UNIQUE KEY uq_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### カラム追加

```sql
-- 20240115130000_add_bio_to_users.sql
ALTER TABLE users ADD COLUMN bio TEXT AFTER name;
```

### インデックス追加

```sql
-- 20240115140000_add_index_users_created_at.sql
CREATE INDEX idx_users_created_at ON users(created_at);
```

### カラム変更（注意が必要）

```sql
-- 20240115150000_change_users_name_length.sql
-- 本番環境では大きなテーブルの変更に時間がかかる可能性がある
ALTER TABLE users MODIFY COLUMN name VARCHAR(200) NOT NULL;
```

### データ移行を伴う変更

```sql
-- 20240115160000_split_name_column.sql
-- 1. 新しいカラムを追加
ALTER TABLE users 
ADD COLUMN first_name VARCHAR(50),
ADD COLUMN last_name VARCHAR(50);

-- 2. データを移行（アプリケーション層で行うことも検討）
UPDATE users SET 
    first_name = SUBSTRING_INDEX(name, ' ', 1),
    last_name = SUBSTRING_INDEX(name, ' ', -1);

-- 3. NOT NULL制約を追加
ALTER TABLE users 
MODIFY COLUMN first_name VARCHAR(50) NOT NULL,
MODIFY COLUMN last_name VARCHAR(50) NOT NULL;

-- 4. 古いカラムを削除（段階的に行うことを推奨）
-- ALTER TABLE users DROP COLUMN name;
```

## パフォーマンス最適化

### クエリ最適化
- `EXPLAIN` を使用してクエリプランを確認
- N+1問題を避ける（JOINまたは一括取得）
- 必要なカラムのみをSELECT
- LIMITを使用してデータ量を制限

```sql
-- 悪い例
SELECT * FROM users;

-- 良い例
SELECT id, email, name FROM users WHERE is_active = 1 LIMIT 100;
```

### インデックスの活用
- WHERE句のカラムにインデックスを作成
- 複合インデックスの順序を最適化
- カバリングインデックスを検討

### パーティショニング
- 大量のデータを扱う場合に検討
- 日付や範囲でパーティション分割

```sql
CREATE TABLE logs (
    id BINARY(16),
    user_id BINARY(16),
    action VARCHAR(50),
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id, created_at)
)
PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2023 VALUES LESS THAN (2024),
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026)
);
```

## トランザクション

### 基本原則
- データの整合性が必要な操作は必ずトランザクション内で実行
- トランザクションは短く保つ
- デッドロックに注意

### 分離レベル
- デフォルトは `REPEATABLE READ`
- 必要に応じて `READ COMMITTED` を使用
- `READ UNCOMMITTED` は使用しない

```sql
-- アプリケーション層で設定
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
START TRANSACTION;
-- クエリ実行
COMMIT;
```

## セキュリティ

### SQLインジェクション対策
- プレースホルダーを必ず使用
- 動的SQLは避ける
- ユーザー入力を直接クエリに埋め込まない

```go
// 良い例（プレースホルダー使用）
db.Query("SELECT * FROM users WHERE email = ?", email)

// 悪い例（SQLインジェクションの危険）
db.Query("SELECT * FROM users WHERE email = '" + email + "'")
```

### 権限管理
- アプリケーション用のDBユーザーは最小限の権限のみ付与
- 本番環境では `DROP`, `TRUNCATE` 権限を付与しない
- 読み取り専用ユーザーを別途作成

### パスワード保存
- パスワードは必ずハッシュ化して保存
- bcryptなどの適切なハッシュアルゴリズムを使用
- ソルトは自動的に生成される

```sql
-- パスワードハッシュを保存
password_hash VARCHAR(255) NOT NULL
```

## ベストプラクティス

### 正規化
- 第3正規形を基本とする
- データの重複を避ける
- パフォーマンスのために意図的に非正規化する場合は文書化

### 論理削除 vs 物理削除
- ユーザーデータなど重要なデータは論理削除
- ログデータなど履歴が不要なデータは物理削除
- 論理削除の場合は `deleted_at` カラムを使用

```sql
-- 論理削除
UPDATE users SET deleted_at = NOW() WHERE id = ?;

-- 論理削除されていないデータのみ取得
SELECT * FROM users WHERE deleted_at IS NULL;
```

### 文字コードと照合順序
- `utf8mb4` を使用（絵文字対応）
- 照合順序は `utf8mb4_unicode_ci` を使用

```sql
CREATE TABLE table_name (
    ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### ストレージエンジン
- `InnoDB` を使用（トランザクション、外部キー対応）
- `MyISAM` は使用しない

### コメント
- テーブルやカラムに説明を追加

```sql
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY COMMENT 'ユーザーID（UUID v7）',
    email VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
    name VARCHAR(100) NOT NULL COMMENT 'ユーザー名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー情報';
```

## 避けるべきパターン

### アンチパターン
- EAV（Entity-Attribute-Value）パターン
- ポリモーフィック関連（1つの外部キーで複数のテーブルを参照）
- カンマ区切りの値を1つのカラムに格納
- 予約語をテーブル名・カラム名に使用
- 過度な正規化（パフォーマンス低下）
- 過度な非正規化（データ不整合のリスク）

```sql
-- 悪い例：カンマ区切りの値
tags VARCHAR(255)  -- 'tag1,tag2,tag3'

-- 良い例：中間テーブル
CREATE TABLE post_tags (
    post_id BINARY(16),
    tag_id BINARY(16),
    PRIMARY KEY (post_id, tag_id)
);
```

## テスト

### テストデータ
- テスト用のシードデータを用意
- 本番データは使用しない
- 個人情報を含まない

### マイグレーションのテスト
- 開発環境でマイグレーションを実行してテスト
- ロールバックが可能か確認
- 大量データでのパフォーマンステスト

## ドキュメント

### スキーマドキュメント
- ER図を作成・更新
- テーブル定義書を維持
- マイグレーション履歴を記録

### 命名規則の文書化
- プロジェクト固有の命名規則を文書化
- 新しいテーブル追加時に参照

