# Backendアーキテクチャ

## アーキテクチャスタイル

このプロジェクトは**クリーンアーキテクチャ**、**ドメイン駆動設計（DDD）**、**CQRS（Command Query Responsibility Segregation）**を採用します。

- クリーンアーキテクチャで技術的な関心事を分離
- DDDでビジネスドメインをモデリング
- CQRSでコマンド（書き込み）とクエリ（読み込み）を分離

これにより、保守性と拡張性の高いシステムを構築します。

## クリーンアーキテクチャの原則

### レイヤー構成

```
┌─────────────────────────────────────┐
│   Frameworks & Drivers (外側)       │
│   - HTTP Handler                    │
│   - Database                        │
│   - External APIs                   │
├─────────────────────────────────────┤
│   Interface Adapters               │
│   - Controllers                     │
│   - Presenters                      │
│   - Gateways                        │
├─────────────────────────────────────┤
│   Use Cases (Application Logic)    │
│   - Business Rules                  │
│   - Application Services            │
├─────────────────────────────────────┤
│   Entities (Domain) (内側)          │
│   - Domain Models                   │
│   - Business Logic                  │
└─────────────────────────────────────┘
```

### 依存関係のルール

- **依存の方向は内側に向かう**：外側のレイヤーは内側のレイヤーに依存できるが、内側は外側に依存してはいけない
- **ドメイン層は独立**：エンティティ（ドメイン）は他のどのレイヤーにも依存しない
- **インターフェースで依存を逆転**：内側のレイヤーがインターフェースを定義し、外側のレイヤーがそれを実装する

## ディレクトリ構造

```
internal/
├── domain/           # エンティティ層（最も内側）
│   ├── model/       # ドメインモデル
│   ├── repository/  # リポジトリインターフェース
│   └── service/     # ドメインサービス
│
├── usecase/         # ユースケース層
│   ├── input/       # 入力DTO
│   ├── output/      # 出力DTO
│   └── interactor/  # ユースケース実装
│
├── adapter/         # インターフェースアダプター層
│   ├── controller/  # HTTPコントローラー
│   ├── presenter/   # プレゼンター
│   └── gateway/     # 外部サービスゲートウェイ
│
└── infrastructure/  # フレームワーク・ドライバー層（最も外側）
    ├── persistence/ # データベース実装
    ├── router/      # ルーティング
    └── config/      # 設定
```

## 各レイヤーの責務

### Domain層（エンティティ）
- ビジネスルールの中核を表現
- 他のレイヤーに依存しない
- データベースやフレームワークの知識を持たない
- 例：User、Order、Product などのドメインモデル

### UseCase層（アプリケーションロジック）
- アプリケーション固有のビジネスルールを実装
- ドメイン層のエンティティを操作
- リポジトリインターフェースを通じてデータにアクセス
- 例：ユーザー登録、注文処理などのユースケース

### Adapter層（インターフェースアダプター）
- 外部とユースケースの間でデータを変換
- HTTPリクエスト/レスポンスの処理
- ユースケースの入出力DTOへの変換
- 例：HTTPハンドラー、JSONシリアライザー

### Infrastructure層（フレームワーク・ドライバー）
- 具体的な技術実装
- データベースアクセス
- 外部APIとの通信
- フレームワークの設定
- 例：PostgreSQL実装、HTTPサーバー設定

## 実装ガイドライン

### 依存性注入
- コンストラクタインジェクションを使用
- インターフェースを通じて依存関係を注入
```go
type UserUseCase struct {
    userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
    return &UserUseCase{userRepo: userRepo}
}
```

### インターフェースの定義場所
- インターフェースは使う側（内側のレイヤー）で定義
- 実装は外側のレイヤーで行う
```go
// domain/repository/user.go
type UserRepository interface {
    Save(ctx context.Context, user *model.User) error
    FindByID(ctx context.Context, id string) (*model.User, error)
}

// infrastructure/persistence/user_repository.go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) Save(ctx context.Context, user *model.User) error {
    // 実装
}
```

### データの流れ
1. HTTPリクエスト → Controller（Adapter層）
2. Controller → UseCase（入力DTOに変換）
3. UseCase → Domain（ビジネスロジック実行）
4. Domain → Repository Interface（データ永続化）
5. Repository Implementation（Infrastructure層）→ Database
6. 結果を逆方向に返す

## ドメイン駆動設計（DDD）

### DDDの戦術的パターン

#### エンティティ（Entity）
- 一意の識別子を持つドメインオブジェクト
- ライフサイクルを通じて同一性を保つ
- ビジネスロジックをカプセル化
```go
type User struct {
    id       UserID
    email    Email
    password HashedPassword
}

func (u *User) ChangeEmail(newEmail Email) error {
    // ビジネスルールの検証
    if err := u.validateEmailChange(newEmail); err != nil {
        return err
    }
    u.email = newEmail
    return nil
}
```

#### 値オブジェクト（Value Object）
- 識別子を持たず、属性の値で区別される
- 不変（イミュータブル）
- 等価性は値で判断
```go
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: value}, nil
}

func (e Email) String() string {
    return e.value
}
```

#### 集約（Aggregate）
- 関連するエンティティと値オブジェクトの集まり
- 集約ルート（Aggregate Root）を通じてのみアクセス
- トランザクション境界を定義
```go
type Order struct {
    id         OrderID
    customerID CustomerID
    items      []OrderItem  // 集約内のエンティティ
    status     OrderStatus
}

func (o *Order) AddItem(product Product, quantity int) error {
    // 集約ルートがビジネスルールを保証
    if o.status != OrderStatusDraft {
        return errors.New("cannot add items to non-draft order")
    }
    item := NewOrderItem(product, quantity)
    o.items = append(o.items, item)
    return nil
}
```

#### リポジトリ（Repository）
- 集約の永続化と取得を抽象化
- コレクションのようなインターフェース
- ドメイン層でインターフェースを定義
```go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Delete(ctx context.Context, id UserID) error
}
```

#### ドメインサービス（Domain Service）
- 単一のエンティティに属さないビジネスロジック
- 複数の集約にまたがる操作
- ステートレス
```go
type UserDomainService struct{}

func (s *UserDomainService) IsDuplicateEmail(
    ctx context.Context,
    email Email,
    repo UserRepository,
) (bool, error) {
    user, err := repo.FindByEmail(ctx, email)
    if err != nil {
        return false, err
    }
    return user != nil, nil
}
```

#### ドメインイベント（Domain Event）
- ドメイン内で発生した重要な出来事
- 過去形で命名
- 他の集約や外部システムへの通知に使用
```go
type UserRegistered struct {
    UserID      UserID
    Email       Email
    OccurredAt  time.Time
}
```

### DDDの戦略的パターン

#### 境界づけられたコンテキスト（Bounded Context）
- ドメインモデルが有効な範囲を明確化
- コンテキストごとに独立したモデルを持つ
- 例：ユーザー管理コンテキスト、注文管理コンテキスト

#### ユビキタス言語（Ubiquitous Language）
- ドメインエキスパートと開発者が共有する共通言語
- コード、ドキュメント、会話で一貫して使用
- `.kiro/steering/ubiquitous_language.md` で定義・管理

### DDDとクリーンアーキテクチャの統合

```
internal/
├── domain/                    # DDD: ドメイン層
│   ├── model/                # エンティティ、値オブジェクト、集約
│   │   ├── user/
│   │   │   ├── user.go      # 集約ルート
│   │   │   ├── email.go     # 値オブジェクト
│   │   │   └── user_id.go   # 値オブジェクト
│   │   └── order/
│   ├── repository/           # リポジトリインターフェース
│   │   ├── user.go
│   │   └── order.go
│   ├── service/              # ドメインサービス
│   │   └── user_service.go
│   └── event/                # ドメインイベント
│       └── user_registered.go
│
├── usecase/                   # DDD: アプリケーション層
│   └── user/
│       └── register_user.go  # ユースケース（アプリケーションサービス）
│
├── adapter/                   # インターフェースアダプター
└── infrastructure/            # インフラストラクチャ層
    └── persistence/
        └── user_repository.go # リポジトリ実装
```

## CQRS（Command Query Responsibility Segregation）

### CQRSの基本原則

コマンド（書き込み）とクエリ（読み込み）の責務を分離し、それぞれに最適化されたモデルを使用します。

#### コマンド側（Command Side）
- **目的**：データの登録・更新・削除
- **モデル**：DDDの集約を使用
- **特徴**：ビジネスルールの検証、トランザクション管理、ドメインイベントの発行

#### クエリ側（Query Side）
- **目的**：データの取得・参照
- **モデル**：要件ごとに自由な構造体を定義
- **特徴**：パフォーマンス最適化、柔軟なデータ構造、結合クエリの活用

### 実装方針

#### コマンド側の実装

コマンド側では集約を使用し、ビジネスルールを保証します。

```go
// domain/model/user/user.go - 集約ルート
type User struct {
    id       UserID
    email    Email
    password HashedPassword
    profile  Profile
}

func (u *User) ChangeEmail(newEmail Email) error {
    // ビジネスルールの検証
    if err := u.validateEmailChange(newEmail); err != nil {
        return err
    }
    u.email = newEmail
    return nil
}

// domain/repository/user.go - コマンド用リポジトリ
type UserCommandRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    Delete(ctx context.Context, id UserID) error
}

// usecase/user/register_user.go - コマンドユースケース
type RegisterUserUseCase struct {
    userRepo domain.UserCommandRepository
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) error {
    // 集約を使用してビジネスロジックを実行
    user, err := domain.NewUser(input.Email, input.Password)
    if err != nil {
        return err
    }
    return uc.userRepo.Save(ctx, user)
}
```

#### クエリ側の実装

クエリ側では要件に応じた専用の構造体を定義し、パフォーマンスを最適化します。

```go
// usecase/user/query/user_list.go - クエリ用の構造体
type UserListItem struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type UserListQuery struct {
    Page     int
    PageSize int
    Keyword  string
}

// usecase/user/query/user_detail.go - 別の要件用の構造体
type UserDetail struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    Name          string    `json:"name"`
    Bio           string    `json:"bio"`
    ProfileImage  string    `json:"profile_image"`
    OrderCount    int       `json:"order_count"`
    TotalSpent    int       `json:"total_spent"`
    LastLoginAt   time.Time `json:"last_login_at"`
    CreatedAt     time.Time `json:"created_at"`
}

// usecase/user/query/repository.go - クエリ用リポジトリ
type UserQueryRepository interface {
    FindList(ctx context.Context, query UserListQuery) ([]UserListItem, error)
    FindDetail(ctx context.Context, userID string) (*UserDetail, error)
    CountByStatus(ctx context.Context, status string) (int, error)
}

// infrastructure/persistence/user_query_repository.go - クエリ実装
type userQueryRepository struct {
    db *sql.DB
}

func (r *userQueryRepository) FindList(ctx context.Context, query UserListQuery) ([]UserListItem, error) {
    // SQLで直接必要なデータのみを取得（JOIN可）
    rows, err := r.db.QueryContext(ctx, `
        SELECT u.id, u.email, u.name, u.created_at
        FROM users u
        WHERE u.name LIKE ?
        ORDER BY u.created_at DESC
        LIMIT ? OFFSET ?
    `, "%"+query.Keyword+"%", query.PageSize, query.Page*query.PageSize)
    // ...
}

func (r *userQueryRepository) FindDetail(ctx context.Context, userID string) (*UserDetail, error) {
    // 複数テーブルをJOINして一度に取得
    var detail UserDetail
    err := r.db.QueryRowContext(ctx, `
        SELECT 
            u.id, u.email, u.name, u.bio, u.profile_image,
            COUNT(o.id) as order_count,
            COALESCE(SUM(o.total_amount), 0) as total_spent,
            u.last_login_at, u.created_at
        FROM users u
        LEFT JOIN orders o ON u.id = o.user_id
        WHERE u.id = ?
        GROUP BY u.id
    `, userID).Scan(&detail.ID, &detail.Email, /* ... */)
    // ...
}
```

### ディレクトリ構造（CQRS対応）

```
internal/
├── domain/                          # コマンド側のドメイン層
│   ├── model/                      # 集約（コマンド用）
│   │   └── user/
│   │       ├── user.go            # 集約ルート
│   │       ├── email.go           # 値オブジェクト
│   │       └── user_id.go
│   └── repository/
│       └── user_command.go        # コマンド用リポジトリIF
│
├── usecase/
│   └── user/
│       ├── command/                # コマンド側ユースケース
│       │   ├── register_user.go
│       │   ├── update_user.go
│       │   └── delete_user.go
│       └── query/                  # クエリ側ユースケース
│           ├── user_list.go       # リスト取得用の構造体
│           ├── user_detail.go     # 詳細取得用の構造体
│           └── repository.go      # クエリ用リポジトリIF
│
└── infrastructure/
    └── persistence/
        ├── user_command_repository.go  # コマンド実装
        └── user_query_repository.go    # クエリ実装
```

### CQRSのメリット

- **パフォーマンス最適化**：クエリ側で必要なデータのみを効率的に取得
- **柔軟性**：画面や要件ごとに最適なデータ構造を定義可能
- **スケーラビリティ**：読み込みと書き込みを独立してスケール可能
- **シンプルさ**：クエリ側でビジネスルールを考慮する必要がない
- **保守性**：コマンドとクエリの変更が互いに影響しにくい

### 実装時の注意点

#### コマンド側
- 必ず集約を使用する
- ビジネスルールは集約内で検証
- トランザクション境界は集約単位
- リポジトリは集約単位で操作

#### クエリ側
- 集約を使用しない（自由な構造体）
- ビジネスルールの検証は不要
- パフォーマンスを優先
- JOIN、集計関数を自由に使用
- 画面や要件ごとに専用の構造体を定義
- 読み取り専用のため、トランザクションは不要な場合が多い

#### データ整合性
- コマンド側で更新したデータは、クエリ側でも即座に反映される（同じDB）
- 将来的にRead/Writeを分離する場合は結果整合性を検討

## 実装の指針

### ドメインモデルの設計
1. ユビキタス言語を使用してモデリング
2. ビジネスルールはドメイン層に集約
3. エンティティと値オブジェクトを適切に使い分ける
4. 集約の境界を慎重に設計（小さく保つ）

### トランザクション管理
- 1つのトランザクションで1つの集約のみを変更
- 複数の集約にまたがる場合は結果整合性を検討
- ドメインイベントを活用して疎結合に

### テスト戦略
- ドメインロジックは単体テストで徹底的にテスト
- リポジトリはモックを使用してテスト
- ユースケースは統合テストでテスト

## メリット

### クリーンアーキテクチャのメリット
- **テスタビリティ**：各レイヤーを独立してテスト可能
- **保守性**：関心事の分離により変更の影響範囲が限定的
- **独立性**：フレームワークやデータベースの変更が容易
- **ビジネスロジックの保護**：ドメイン層が外部の変更から守られる

### DDDのメリット
- **ビジネスとの整合性**：ドメインモデルがビジネスを正確に表現
- **複雑性の管理**：境界づけられたコンテキストで複雑さを分割
- **共通言語**：チーム全体で一貫した用語を使用
- **変更への強さ**：ビジネスルールの変更に柔軟に対応

### CQRSのメリット
- **最適化**：コマンドとクエリをそれぞれの目的に最適化
- **柔軟性**：クエリ側で要件に応じた自由なデータ構造
- **パフォーマンス**：読み込み処理の高速化
- **保守性**：書き込みと読み込みの変更が独立
