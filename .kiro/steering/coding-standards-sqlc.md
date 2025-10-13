# sqlcに関するコーディング規約

## 採用ツール
- sqlc（SQL to Go code generator）
- データベースアクセスのための型安全なGoコードを自動生成

## 基本原則

- 生SQLを書き、型安全なGoコードを生成する
- クエリは読みやすく、保守しやすい形で記述する
- パフォーマンスを考慮したクエリ設計を行う
- トランザクション管理を適切に行う
- CQRSパターンに従い、コマンドとクエリを分離する

## ディレクトリ構造

```
services/manager/backend/
├── internal/
│   └── infrastructure/
│       └── persistence/
│           ├── sqlc/              # sqlc関連ファイル
│           │   ├── db.go         # sqlc生成コード
│           │   ├── models.go     # sqlc生成コード
│           │   ├── queries.sql.go # sqlc生成コード
│           │   └── querier.go    # sqlc生成コード
│           ├── query/             # SQLクエリファイル
│           │   ├── user_command.sql  # コマンド用クエリ
│           │   └── user_query.sql    # クエリ用クエリ
│           ├── schema/            # スキーマ定義
│           │   └── schema.sql
│           ├── user_command_repository.go  # リポジトリ実装
│           └── user_query_repository.go
│
└── sqlc.yaml                      # sqlc設定ファイル
```

## sqlc.yaml 設定

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "internal/infrastructure/persistence/query"
    schema: "internal/infrastructure/persistence/schema"
    gen:
      go:
        package: "sqlc"
        out: "internal/infrastructure/persistence/sqlc"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
        emit_pointers_for_null_types: true
```

## クエリファイルの命名規則

### ファイル名
- コマンド用：`{集約名}_command.sql`
- クエリ用：`{集約名}_query.sql`
- 例：`user_command.sql`, `user_query.sql`

### クエリ名
- コマンド：`{動詞}{集約名}`（例：`CreateUser`, `UpdateUser`, `DeleteUser`）
- クエリ：`Get{集約名}`, `List{集約名}`, `Count{集約名}`
- 小文字のスネークケースは使わない（sqlcがGoの命名規則に変換）

## コマンド側のクエリ

### 基本的なCRUD操作

```sql
-- name: CreateUser :execresult
INSERT INTO users (
    id,
    email,
    password_hash,
    name,
    created_at,
    updated_at
) VALUES (
    sqlc.arg('id'),
    sqlc.arg('email'),
    sqlc.arg('password_hash'),
    sqlc.arg('name'),
    sqlc.arg('created_at'),
    sqlc.arg('updated_at')
);

-- name: GetUserByID :one
SELECT 
    id,
    email,
    password_hash,
    name,
    bio,
    is_active,
    email_verified_at,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT 
    id,
    email,
    password_hash,
    name,
    bio,
    is_active,
    email_verified_at,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE email = sqlc.arg('email') AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET
    email = sqlc.arg('email'),
    name = sqlc.arg('name'),
    bio = sqlc.arg('bio'),
    updated_at = sqlc.arg('updated_at')
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    password_hash = sqlc.arg('password_hash'),
    updated_at = sqlc.arg('updated_at')
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET
    deleted_at = sqlc.arg('deleted_at'),
    updated_at = sqlc.arg('updated_at')
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg('id');
```

### 注釈の使い分け

- `:exec` - 結果を返さない（UPDATE, DELETE）
- `:execresult` - 影響を受けた行数を返す（INSERT）
- `:one` - 1行を返す（SELECT）
- `:many` - 複数行を返す（SELECT）

## クエリ側のクエリ

### リスト取得

```sql
-- name: ListUsers :many
SELECT 
    u.id,
    u.email,
    u.name,
    u.is_active,
    u.created_at
FROM users u
WHERE 
    u.deleted_at IS NULL
    AND (sqlc.arg('keyword') = '' OR u.name LIKE CONCAT('%', sqlc.arg('keyword'), '%'))
ORDER BY u.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsers :one
SELECT COUNT(*)
FROM users u
WHERE 
    u.deleted_at IS NULL
    AND (sqlc.arg('keyword') = '' OR u.name LIKE CONCAT('%', sqlc.arg('keyword'), '%'));
```

### 詳細取得（JOIN使用）

```sql
-- name: GetUserDetail :one
SELECT 
    u.id,
    u.email,
    u.name,
    u.bio,
    u.profile_image,
    u.is_active,
    u.email_verified_at,
    u.last_login_at,
    u.created_at,
    COUNT(DISTINCT o.id) as order_count,
    COALESCE(SUM(o.total_amount), 0) as total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id AND o.deleted_at IS NULL
WHERE u.id = sqlc.arg('id') AND u.deleted_at IS NULL
GROUP BY u.id;
```

### 集計クエリ

```sql
-- name: CountUsersByStatus :one
SELECT COUNT(*)
FROM users
WHERE is_active = sqlc.arg('is_active') AND deleted_at IS NULL;

-- name: GetUserStatistics :one
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN is_active = 1 THEN 1 END) as active_users,
    COUNT(CASE WHEN email_verified_at IS NOT NULL THEN 1 END) as verified_users
FROM users
WHERE deleted_at IS NULL;
```

## NULL値の扱い

### NULLを許可するカラム

sqlcは自動的にNULL許可カラムをポインタ型に変換します。

```sql
-- name: GetUser :one
SELECT 
    id,
    email,
    bio,              -- TEXT NULL
    email_verified_at, -- DATETIME(6) NULL
    deleted_at        -- DATETIME(6) NULL
FROM users
WHERE id = sqlc.arg('id');
```

生成されるGoコード：
```go
type User struct {
    ID              []byte
    Email           string
    Bio             sql.NullString    // NULLを許可
    EmailVerifiedAt sql.NullTime      // NULLを許可
    DeletedAt       sql.NullTime      // NULLを許可
}
```

### COALESCE の使用

NULL値をデフォルト値に変換する場合：

```sql
-- name: GetUserWithDefaults :one
SELECT 
    id,
    email,
    COALESCE(bio, '') as bio,
    COALESCE(profile_image, '') as profile_image
FROM users
WHERE id = sqlc.arg('id');
```

## パラメータの扱い

### 名前付きパラメータ

sqlcは `sqlc.arg()` を使用して名前付きパラメータを定義します。

```sql
-- name: CreateUser :execresult
INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
VALUES (sqlc.arg('id'), sqlc.arg('email'), sqlc.arg('password_hash'), sqlc.arg('name'), sqlc.arg('created_at'), sqlc.arg('updated_at'));
```

生成されるGoコード：
```go
type CreateUserParams struct {
    ID           []byte
    Email        string
    PasswordHash string
    Name         string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
```

### オプショナルパラメータ

条件付き検索の実装：

```sql
-- name: SearchUsers :many
SELECT id, email, name, created_at
FROM users
WHERE 
    deleted_at IS NULL
    AND (sqlc.arg('keyword') = '' OR email LIKE CONCAT('%', sqlc.arg('keyword'), '%'))
    AND (sqlc.arg('is_active') = 0 OR is_active = sqlc.arg('is_active'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
```

## トランザクション管理

### トランザクション用のインターフェース

sqlcは `Querier` インターフェースを生成します。

```go
// infrastructure/persistence/user_command_repository.go
type userCommandRepository struct {
    db *sql.DB
}

func (r *userCommandRepository) Save(ctx context.Context, user *model.User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    queries := sqlc.New(tx)
    
    // sqlc生成メソッドを使用
    _, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
        ID:           user.ID().Bytes(),
        Email:        user.Email().String(),
        PasswordHash: user.Password().String(),
        Name:         user.Name(),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    })
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

### WithTx パターン

トランザクションを引数で受け取るパターン：

```go
func (r *userCommandRepository) SaveWithTx(
    ctx context.Context,
    tx *sql.Tx,
    user *model.User,
) error {
    queries := sqlc.New(tx)
    
    _, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
        ID:           user.ID().Bytes(),
        Email:        user.Email().String(),
        PasswordHash: user.Password().String(),
        Name:         user.Name(),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    })
    
    return err
}
```

## 型変換

### UUID（BINARY(16)）の扱い

```go
// ドメインモデル → sqlcパラメータ
func toSQLCUserID(id model.UserID) []byte {
    return id.Bytes()
}

// sqlc結果 → ドメインモデル
func fromSQLCUserID(b []byte) (model.UserID, error) {
    return model.NewUserIDFromBytes(b)
}
```

### 値オブジェクトの変換

```go
// Email値オブジェクト → string
func toSQLCEmail(email model.Email) string {
    return email.String()
}

// string → Email値オブジェクト
func fromSQLCEmail(s string) (model.Email, error) {
    return model.NewEmail(s)
}
```

### 日時の扱い

```go
// time.Time → DATETIME(6)
createdAt := time.Now().UTC()

// sql.NullTime → *time.Time
func fromSQLCNullTime(nt sql.NullTime) *time.Time {
    if !nt.Valid {
        return nil
    }
    return &nt.Time
}
```

## リポジトリ実装パターン

### コマンド側リポジトリ

```go
// infrastructure/persistence/user_command_repository.go
type userCommandRepository struct {
    db *sql.DB
}

func NewUserCommandRepository(db *sql.DB) domain.UserCommandRepository {
    return &userCommandRepository{db: db}
}

func (r *userCommandRepository) Save(ctx context.Context, user *model.User) error {
    queries := sqlc.New(r.db)
    
    // 既存チェック
    existing, err := queries.GetUserByID(ctx, user.ID().Bytes())
    if err != nil && err != sql.ErrNoRows {
        return fmt.Errorf("failed to check existing user: %w", err)
    }
    
    if existing.ID != nil {
        // 更新
        return queries.UpdateUser(ctx, sqlc.UpdateUserParams{
            Email:     user.Email().String(),
            Name:      user.Name(),
            Bio:       toSQLCNullString(user.Bio()),
            UpdatedAt: time.Now().UTC(),
            ID:        user.ID().Bytes(),
        })
    }
    
    // 新規作成
    _, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
        ID:           user.ID().Bytes(),
        Email:        user.Email().String(),
        PasswordHash: user.Password().String(),
        Name:         user.Name(),
        CreatedAt:    time.Now().UTC(),
        UpdatedAt:    time.Now().UTC(),
    })
    
    return err
}

func (r *userCommandRepository) FindByID(ctx context.Context, id model.UserID) (*model.User, error) {
    queries := sqlc.New(r.db)
    
    row, err := queries.GetUserByID(ctx, id.Bytes())
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return r.toUserModel(row)
}

func (r *userCommandRepository) toUserModel(row sqlc.User) (*model.User, error) {
    id, err := model.NewUserIDFromBytes(row.ID)
    if err != nil {
        return nil, err
    }
    
    email, err := model.NewEmail(row.Email)
    if err != nil {
        return nil, err
    }
    
    password, err := model.NewHashedPasswordFromHash(row.PasswordHash)
    if err != nil {
        return nil, err
    }
    
    return model.ReconstructUser(id, email, password, row.Name, fromSQLCNullString(row.Bio))
}
```

### クエリ側リポジトリ

```go
// infrastructure/persistence/user_query_repository.go
type userQueryRepository struct {
    db *sql.DB
}

func NewUserQueryRepository(db *sql.DB) usecase.UserQueryRepository {
    return &userQueryRepository{db: db}
}

func (r *userQueryRepository) FindList(
    ctx context.Context,
    query usecase.UserListQuery,
) ([]usecase.UserListItem, error) {
    queries := sqlc.New(r.db)
    
    rows, err := queries.ListUsers(ctx, sqlc.ListUsersParams{
        Keyword: query.Keyword,
        Limit:   int32(query.PageSize),
        Offset:  int32(query.Page * query.PageSize),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    
    items := make([]usecase.UserListItem, len(rows))
    for i, row := range rows {
        items[i] = usecase.UserListItem{
            ID:        uuid.Must(uuid.FromBytes(row.ID)).String(),
            Email:     row.Email,
            Name:      row.Name,
            IsActive:  row.IsActive == 1,
            CreatedAt: row.CreatedAt,
        }
    }
    
    return items, nil
}

func (r *userQueryRepository) FindDetail(
    ctx context.Context,
    userID string,
) (*usecase.UserDetail, error) {
    queries := sqlc.New(r.db)
    
    id, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user id: %w", err)
    }
    
    row, err := queries.GetUserDetail(ctx, id[:])
    if err == sql.ErrNoRows {
        return nil, usecase.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user detail: %w", err)
    }
    
    return &usecase.UserDetail{
        ID:            uuid.Must(uuid.FromBytes(row.ID)).String(),
        Email:         row.Email,
        Name:          row.Name,
        Bio:           fromSQLCNullString(row.Bio),
        ProfileImage:  fromSQLCNullString(row.ProfileImage),
        IsActive:      row.IsActive == 1,
        OrderCount:    int(row.OrderCount),
        TotalSpent:    int(row.TotalSpent),
        LastLoginAt:   fromSQLCNullTime(row.LastLoginAt),
        CreatedAt:     row.CreatedAt,
    }, nil
}
```

## ヘルパー関数

### NULL値変換

```go
// string → sql.NullString
func toSQLCNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}

// sql.NullString → string
func fromSQLCNullString(ns sql.NullString) string {
    if !ns.Valid {
        return ""
    }
    return ns.String
}

// *time.Time → sql.NullTime
func toSQLCNullTime(t *time.Time) sql.NullTime {
    if t == nil {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: *t, Valid: true}
}

// sql.NullTime → *time.Time
func fromSQLCNullTime(nt sql.NullTime) *time.Time {
    if !nt.Valid {
        return nil
    }
    return &nt.Time
}
```

## エラーハンドリング

### 標準的なエラー処理

```go
func (r *userCommandRepository) FindByEmail(
    ctx context.Context,
    email model.Email,
) (*model.User, error) {
    queries := sqlc.New(r.db)
    
    row, err := queries.GetUserByEmail(ctx, email.String())
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    
    return r.toUserModel(row)
}
```

### 重複エラーの検出

```go
import (
    "github.com/go-sql-driver/mysql"
)

func (r *userCommandRepository) Save(ctx context.Context, user *model.User) error {
    queries := sqlc.New(r.db)
    
    _, err := queries.CreateUser(ctx, /* params */)
    if err != nil {
        // MySQL重複エラー（1062）の検出
        if mysqlErr, ok := err.(*mysql.MySQLError); ok {
            if mysqlErr.Number == 1062 {
                return domain.ErrUserEmailDuplicate
            }
        }
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

## ベストプラクティス

### クエリの最適化
- 必要なカラムのみをSELECT
- インデックスを活用したWHERE句
- LIMITとOFFSETでページネーション
- JOINは必要な場合のみ使用

### 型安全性
- sqlcの生成コードを信頼する
- 手動でSQLを書かない
- パラメータは構造体で渡す

### パフォーマンス
- N+1問題を避ける（JOINまたは一括取得）
- プリペアドステートメントを活用（sqlcが自動で行う）
- 適切なインデックスを設定

### 保守性
- クエリファイルを集約ごとに分割
- コマンドとクエリを明確に分離
- 複雑なクエリにはコメントを追加

## コード生成

### sqlcの実行

```bash
# Makefile
.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: sqlc-verify
sqlc-verify:
	sqlc verify
```

### 生成されるファイル
- `db.go` - DBTX インターフェース
- `models.go` - テーブルに対応する構造体
- `queries.sql.go` - クエリメソッド
- `querier.go` - Querier インターフェース

### 生成コードの扱い
- 生成されたコードは編集しない
- バージョン管理にコミットする
- CIで生成コードの整合性をチェック

## テスト

### リポジトリのテスト

```go
func TestUserCommandRepository_Save(t *testing.T) {
    // テスト用DBの準備
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserCommandRepository(db)
    
    // テストデータ
    user, err := model.NewUser(
        model.MustNewEmail("test@example.com"),
        model.MustNewPassword("password123"),
        "Test User",
    )
    require.NoError(t, err)
    
    // 保存
    err = repo.Save(context.Background(), user)
    require.NoError(t, err)
    
    // 取得して確認
    found, err := repo.FindByID(context.Background(), user.ID())
    require.NoError(t, err)
    assert.Equal(t, user.Email().String(), found.Email().String())
}
```

## 避けるべきパターン

- 生成されたコードを手動で編集する
- 動的SQLを使用する（sqlcの範囲外）
- 複雑すぎるクエリ（可読性が低下）
- トランザクション管理の漏れ
- エラーハンドリングの省略
- NULL値の不適切な扱い
