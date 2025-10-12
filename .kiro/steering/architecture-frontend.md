# Frontendアーキテクチャ

## アーキテクチャスタイル

このプロジェクトは**レイヤードアーキテクチャ**と**プレゼンテーション・ドメイン分離**を採用します。

- UIとビジネスロジックを明確に分離
- 状態管理を集中化し、データフローを明確に
- コンポーネントの責務を明確化し、再利用性を向上

これにより、保守性と拡張性の高いフロントエンドを構築します。

## レイヤー構成

```
┌─────────────────────────────────────┐
│   Presentation Layer (外側)         │
│   - Pages                           │
│   - UI Components                   │
│   - Routing                         │
├─────────────────────────────────────┤
│   Application Layer                │
│   - Hooks                           │
│   - State Management (Jotai)       │
│   - Form Validation                 │
├─────────────────────────────────────┤
│   Domain Layer                     │
│   - Types                           │
│   - Business Logic                  │
│   - Validation Rules                │
├─────────────────────────────────────┤
│   Infrastructure Layer (内側)       │
│   - API Client                      │
│   - HTTP Communication              │
│   - External Services               │
└─────────────────────────────────────┘
```

## ディレクトリ構造

```
src/
├── pages/                    # ページコンポーネント（Presentation Layer）
│   ├── auth/
│   │   ├── LoginPage.tsx
│   │   └── RegisterPage.tsx
│   └── user/
│       ├── UserListPage.tsx
│       └── UserDetailPage.tsx
│
├── components/               # 再利用可能なコンポーネント（Presentation Layer）
│   ├── ui/                  # 基本的なUIコンポーネント（shadcn/ui）
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   └── Card.tsx
│   └── features/            # 機能別コンポーネント
│       └── user/
│           ├── UserCard.tsx
│           ├── UserList.tsx
│           └── UserForm.tsx
│
├── hooks/                    # カスタムHooks（Application Layer）
│   ├── api/                 # API通信用Hooks
│   │   ├── useUsers.ts
│   │   └── useAuth.ts
│   ├── form/                # フォーム管理用Hooks
│   │   └── useUserForm.ts
│   └── state/               # 状態管理用Hooks
│       └── useAuthState.ts
│
├── stores/                   # グローバル状態管理（Application Layer）
│   ├── authStore.ts         # Jotai atoms
│   └── userStore.ts
│
├── domain/                   # ドメイン層（Domain Layer）
│   ├── models/              # ドメインモデル・型定義
│   │   ├── user.ts
│   │   └── auth.ts
│   ├── validation/          # バリデーションルール
│   │   ├── userValidation.ts
│   │   └── authValidation.ts
│   └── services/            # ドメインサービス（ビジネスロジック）
│       └── userService.ts
│
├── api/                      # API通信（Infrastructure Layer）
│   ├── client.ts            # APIクライアント設定（ky）
│   ├── generated/           # OpenAPI Generator出力
│   │   ├── api.ts
│   │   └── models.ts
│   └── endpoints/           # APIエンドポイント定義
│       ├── userApi.ts
│       └── authApi.ts
│
├── routes/                   # ルーティング設定（TanStack Router）
│   ├── __root.tsx
│   ├── index.tsx
│   └── auth/
│       ├── login.tsx
│       └── register.tsx
│
├── utils/                    # ユーティリティ関数
│   ├── formatDate.ts
│   └── validation.ts
│
└── constants/                # 定数
    └── index.ts
```

## 各レイヤーの責務

### Presentation Layer（プレゼンテーション層）
- **責務**：UIの表示とユーザーインタラクション
- **含まれるもの**：
  - ページコンポーネント（Pages）
  - UIコンポーネント（Components）
  - ルーティング（Routes）
- **ルール**：
  - ビジネスロジックを含まない
  - Hooksを通じてデータを取得・更新
  - 見た目とユーザー操作のみに集中

```typescript
// pages/user/UserListPage.tsx
export const UserListPage = () => {
  const { data: users, isLoading } = useUsers();

  if (isLoading) return <LoadingSpinner />;

  return (
    <div>
      <h1>ユーザー一覧</h1>
      <UserList users={users} />
    </div>
  );
};
```

### Application Layer（アプリケーション層）
- **責務**：アプリケーション固有のロジックと状態管理
- **含まれるもの**：
  - カスタムHooks
  - グローバル状態管理（Jotai）
  - フォーム管理
- **ルール**：
  - UIとドメインロジックを橋渡し
  - データフェッチングとキャッシング
  - 状態の管理と更新

```typescript
// hooks/api/useUsers.ts
export const useUsers = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => userApi.fetchUsers(),
  });
};

// stores/authStore.ts
export const authAtom = atom<AuthState>({
  user: null,
  isAuthenticated: false,
});
```

### Domain Layer（ドメイン層）
- **責務**：ビジネスルールとドメインモデル
- **含まれるもの**：
  - 型定義（Models）
  - バリデーションルール
  - ドメインサービス
- **ルール**：
  - フレームワークに依存しない
  - 純粋なTypeScriptコード
  - ビジネスロジックの中核

```typescript
// domain/models/user.ts
export type User = {
  id: string;
  email: string;
  name: string;
  createdAt: Date;
};

export type CreateUserInput = {
  email: string;
  password: string;
  name: string;
};

// domain/validation/userValidation.ts
export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

export const validatePassword = (password: string): boolean => {
  return password.length >= 8;
};

// domain/services/userService.ts
export const userService = {
  canRegister: (input: CreateUserInput): boolean => {
    return validateEmail(input.email) && validatePassword(input.password);
  },
};
```

### Infrastructure Layer（インフラストラクチャ層）
- **責務**：外部システムとの通信
- **含まれるもの**：
  - APIクライアント（ky）
  - HTTP通信
  - OpenAPI生成コード
- **ルール**：
  - 技術的な実装の詳細を隠蔽
  - エラーハンドリング
  - レスポンスの変換

```typescript
// api/client.ts
import ky from 'ky';

export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem('token');
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
  },
});

// api/endpoints/userApi.ts
export const userApi = {
  fetchUsers: async (): Promise<User[]> => {
    return apiClient.get('users').json<User[]>();
  },
  
  createUser: async (input: CreateUserInput): Promise<User> => {
    return apiClient.post('users', { json: input }).json<User>();
  },
};
```

## 状態管理戦略（Jotai）

### Atomの設計原則

- **小さく分割**：関心事ごとにAtomを分ける
- **派生状態**：計算可能な状態は派生Atomで表現
- **不変性**：状態の更新は新しいオブジェクトを作成

```typescript
// stores/authStore.ts
import { atom } from 'jotai';
import type { User } from '@/domain/models/user';

// プリミティブAtom
export const userAtom = atom<User | null>(null);
export const tokenAtom = atom<string | null>(null);

// 派生Atom（読み取り専用）
export const isAuthenticatedAtom = atom((get) => {
  return get(userAtom) !== null && get(tokenAtom) !== null;
});

// 書き込み可能な派生Atom
export const loginAtom = atom(
  null,
  (get, set, { user, token }: { user: User; token: string }) => {
    set(userAtom, user);
    set(tokenAtom, token);
    localStorage.setItem('token', token);
  }
);

export const logoutAtom = atom(null, (get, set) => {
  set(userAtom, null);
  set(tokenAtom, null);
  localStorage.removeItem('token');
});
```

### Atomの使用

```typescript
// hooks/state/useAuthState.ts
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { userAtom, isAuthenticatedAtom, loginAtom, logoutAtom } from '@/stores/authStore';

export const useAuthState = () => {
  const user = useAtomValue(userAtom);
  const isAuthenticated = useAtomValue(isAuthenticatedAtom);
  const login = useSetAtom(loginAtom);
  const logout = useSetAtom(logoutAtom);

  return { user, isAuthenticated, login, logout };
};
```

## データフェッチング戦略

### TanStack Query（React Query）の活用

- **サーバーステート管理**：APIから取得したデータのキャッシュと同期
- **自動リフェッチ**：データの鮮度を保つ
- **楽観的更新**：UXの向上

```typescript
// hooks/api/useUsers.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userApi } from '@/api/endpoints/userApi';
import type { CreateUserInput } from '@/domain/models/user';

export const useUsers = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: userApi.fetchUsers,
    staleTime: 5 * 60 * 1000, // 5分間はキャッシュを使用
  });
};

export const useUser = (userId: string) => {
  return useQuery({
    queryKey: ['users', userId],
    queryFn: () => userApi.fetchUser(userId),
    enabled: !!userId, // userIdがある場合のみ実行
  });
};

export const useCreateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateUserInput) => userApi.createUser(input),
    onSuccess: () => {
      // ユーザー一覧を再取得
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateUserInput }) =>
      userApi.updateUser(id, input),
    onSuccess: (data, variables) => {
      // 特定のユーザーのキャッシュを更新
      queryClient.setQueryData(['users', variables.id], data);
      // ユーザー一覧も無効化
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};
```

## ルーティング（TanStack Router）

### ファイルベースルーティング

```typescript
// routes/__root.tsx
import { Outlet, createRootRoute } from '@tanstack/react-router';

export const Route = createRootRoute({
  component: () => (
    <div>
      <nav>{/* ナビゲーション */}</nav>
      <Outlet />
    </div>
  ),
});

// routes/index.tsx
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: HomePage,
});

// routes/users/index.tsx
export const Route = createFileRoute('/users/')({
  component: UserListPage,
  loader: async ({ context }) => {
    // ページ遷移前にデータをプリフェッチ
    return context.queryClient.ensureQueryData({
      queryKey: ['users'],
      queryFn: userApi.fetchUsers,
    });
  },
});

// routes/users/$userId.tsx
export const Route = createFileRoute('/users/$userId')({
  component: UserDetailPage,
  loader: async ({ params, context }) => {
    return context.queryClient.ensureQueryData({
      queryKey: ['users', params.userId],
      queryFn: () => userApi.fetchUser(params.userId),
    });
  },
});
```

### 型安全なナビゲーション

```typescript
import { useNavigate } from '@tanstack/react-router';

const navigate = useNavigate();

// 型安全なナビゲーション
navigate({ to: '/users/$userId', params: { userId: '123' } });
```

## フォーム管理

### React Hook Formとの統合

```typescript
// hooks/form/useUserForm.ts
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { validateEmail, validatePassword } from '@/domain/validation/userValidation';

const userSchema = z.object({
  email: z.string().refine(validateEmail, 'メールアドレスの形式が正しくありません'),
  password: z.string().refine(validatePassword, 'パスワードは8文字以上である必要があります'),
  name: z.string().min(1, '名前は必須です'),
});

type UserFormData = z.infer<typeof userSchema>;

export const useUserForm = () => {
  return useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      email: '',
      password: '',
      name: '',
    },
  });
};
```

### フォームコンポーネント

```typescript
// components/features/user/UserForm.tsx
import { useUserForm } from '@/hooks/form/useUserForm';
import { useCreateUser } from '@/hooks/api/useUsers';

export const UserForm = () => {
  const { register, handleSubmit, formState: { errors } } = useUserForm();
  const { mutate: createUser, isPending } = useCreateUser();

  const onSubmit = (data: UserFormData) => {
    createUser(data, {
      onSuccess: () => {
        toast.success('ユーザーを作成しました');
      },
      onError: (error) => {
        toast.error('ユーザーの作成に失敗しました');
      },
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input {...register('email')} error={errors.email?.message} />
      <Input {...register('password')} type="password" error={errors.password?.message} />
      <Input {...register('name')} error={errors.name?.message} />
      <Button type="submit" disabled={isPending}>
        登録
      </Button>
    </form>
  );
};
```

## コンポーネント設計原則

### プレゼンテーショナルコンポーネントとコンテナコンポーネント

#### プレゼンテーショナルコンポーネント
- **責務**：見た目の表示のみ
- **特徴**：Propsを受け取り、UIを描画
- **状態**：ローカルなUI状態のみ（開閉状態など）

```typescript
// components/features/user/UserCard.tsx
type UserCardProps = {
  user: User;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
};

export const UserCard = ({ user, onEdit, onDelete }: UserCardProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{user.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>{user.email}</p>
      </CardContent>
      <CardFooter>
        <Button onClick={() => onEdit(user.id)}>編集</Button>
        <Button onClick={() => onDelete(user.id)} variant="destructive">
          削除
        </Button>
      </CardFooter>
    </Card>
  );
};
```

#### コンテナコンポーネント
- **責務**：データの取得とビジネスロジック
- **特徴**：Hooksを使用してデータを管理
- **状態**：グローバル状態やサーバーステート

```typescript
// components/features/user/UserListContainer.tsx
export const UserListContainer = () => {
  const { data: users, isLoading } = useUsers();
  const { mutate: deleteUser } = useDeleteUser();
  const navigate = useNavigate();

  const handleEdit = (id: string) => {
    navigate({ to: '/users/$userId/edit', params: { userId: id } });
  };

  const handleDelete = (id: string) => {
    if (confirm('本当に削除しますか？')) {
      deleteUser(id);
    }
  };

  if (isLoading) return <LoadingSpinner />;

  return (
    <div>
      {users?.map((user) => (
        <UserCard
          key={user.id}
          user={user}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ))}
    </div>
  );
};
```

## エラーハンドリング

### APIエラーの処理

```typescript
// api/client.ts
export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_BASE_URL,
  hooks: {
    afterResponse: [
      async (request, options, response) => {
        if (!response.ok) {
          const error = await response.json();
          throw new ApiError(error.message, response.status);
        }
      },
    ],
  },
});

// domain/models/error.ts
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}
```

### エラーバウンダリ

```typescript
// components/ErrorBoundary.tsx
export class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return (
        <div>
          <h1>エラーが発生しました</h1>
          <p>{this.state.error?.message}</p>
        </div>
      );
    }
    return this.props.children;
  }
}
```

## テスト戦略

### コンポーネントのテスト

```typescript
// components/features/user/UserCard.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UserCard } from './UserCard';

describe('UserCard', () => {
  const mockUser = {
    id: '1',
    email: 'test@example.com',
    name: 'Test User',
    createdAt: new Date(),
  };

  it('ユーザー情報を表示する', () => {
    render(<UserCard user={mockUser} onEdit={jest.fn()} onDelete={jest.fn()} />);
    
    expect(screen.getByText('Test User')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('編集ボタンをクリックするとonEditが呼ばれる', async () => {
    const onEdit = jest.fn();
    render(<UserCard user={mockUser} onEdit={onEdit} onDelete={jest.fn()} />);
    
    await userEvent.click(screen.getByText('編集'));
    expect(onEdit).toHaveBeenCalledWith('1');
  });
});
```

### Hooksのテスト

```typescript
// hooks/api/useUsers.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useUsers } from './useUsers';

const createWrapper = () => {
  const queryClient = new QueryClient();
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('useUsers', () => {
  it('ユーザー一覧を取得する', async () => {
    const { result } = renderHook(() => useUsers(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(2);
  });
});
```

## パフォーマンス最適化

### コンポーネントのメモ化

```typescript
// 重い計算のメモ化
const expensiveValue = useMemo(() => {
  return computeExpensiveValue(data);
}, [data]);

// コールバックのメモ化
const handleClick = useCallback(() => {
  doSomething(id);
}, [id]);

// コンポーネントのメモ化
export const UserCard = React.memo(({ user }: UserCardProps) => {
  return <div>{user.name}</div>;
});
```

### コード分割

```typescript
// ページの遅延ロード
const UserListPage = lazy(() => import('./pages/user/UserListPage'));
const UserDetailPage = lazy(() => import('./pages/user/UserDetailPage'));

// ルートでの使用
<Suspense fallback={<LoadingSpinner />}>
  <UserListPage />
</Suspense>
```

## ベストプラクティス

### 依存関係の方向
- **外側から内側へ**：Presentation → Application → Domain → Infrastructure
- **ドメイン層は独立**：他のレイヤーに依存しない
- **インフラ層は差し替え可能**：APIクライアントの実装を変更しても影響が少ない

### 状態管理の原則
- **サーバーステート**：TanStack Queryで管理
- **グローバルクライアントステート**：Jotaiで管理
- **ローカルUIステート**：useStateで管理
- **フォームステート**：React Hook Formで管理

### コンポーネント設計
- **単一責任の原則**：1つのコンポーネントは1つの責務
- **Props Drillingの回避**：深いネストはJotaiで解決
- **再利用性**：汎用的なコンポーネントはui/に配置
- **テスタビリティ**：ビジネスロジックはHooksに分離

## メリット

### レイヤードアーキテクチャのメリット
- **保守性**：関心事の分離により変更の影響範囲が限定的
- **テスタビリティ**：各レイヤーを独立してテスト可能
- **再利用性**：ドメインロジックを複数のコンポーネントで共有
- **拡張性**：新機能の追加が容易

### 状態管理の明確化
- **データフローの可視化**：どこで状態が管理されているか明確
- **パフォーマンス**：必要な部分のみ再レンダリング
- **デバッグ**：状態の変更を追跡しやすい

### 型安全性
- **TypeScript**：コンパイル時にエラーを検出
- **OpenAPI Generator**：APIの型を自動生成
- **Zod**：実行時のバリデーションと型推論
