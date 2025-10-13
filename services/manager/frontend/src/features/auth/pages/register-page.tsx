import { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { APIError } from '@/lib/api-error';
import { useRegisterUser } from '@/features/auth/hooks/use-register';
import { RegisterPayload, RegisterResponse } from '@/features/auth/types';

const registerSchema = z
  .object({
    email: z.string().email('メールアドレスの形式が正しくありません'),
    password: z
      .string()
      .min(8, 'パスワードは8文字以上で、英字と数字を含む必要があります')
      .regex(/^(?=.*[a-zA-Z])(?=.*\d).+$/, 'パスワードは8文字以上で、英字と数字を含む必要があります'),
    passwordConfirmation: z.string()
  })
  .refine((data) => data.password === data.passwordConfirmation, {
    message: 'パスワードが一致しません',
    path: ['passwordConfirmation']
  });

type RegisterFormValues = z.infer<typeof registerSchema>;

type FieldKey = keyof RegisterFormValues | 'root';

const fieldMap: Record<string, FieldKey> = {
  email: 'email',
  password: 'password',
  password_confirmation: 'passwordConfirmation'
};

const toPayload = (values: RegisterFormValues): RegisterPayload => ({
  email: values.email,
  password: values.password,
  passwordConfirmation: values.passwordConfirmation
});

export const RegisterPage = (): JSX.Element => {
  const registerMutation = useRegisterUser();
  const [success, setSuccess] = useState<RegisterResponse | null>(null);

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: '',
      password: '',
      passwordConfirmation: ''
    }
  });

  const setServerErrors = (error: APIError) => {
    if (error.details.length === 0) {
      form.setError('root', { type: 'server', message: error.message });
      return;
    }

    error.details.forEach((detail) => {
      const key = detail.field ? fieldMap[detail.field] ?? 'root' : 'root';
      if (key === 'root') {
        form.setError('root', { type: 'server', message: detail.message ?? error.message });
      } else {
        form.setError(key, { type: 'server', message: detail.message ?? error.message });
      }
    });
  };

  const onSubmit = form.handleSubmit(async (values) => {
    setSuccess(null);
    form.clearErrors();

    try {
      const result = await registerMutation.mutateAsync(toPayload(values));
      setSuccess(result);
      form.reset({ email: '', password: '', passwordConfirmation: '' });
    } catch (error) {
      if (error instanceof APIError) {
        setServerErrors(error);
        return;
      }

      form.setError('root', { type: 'server', message: '予期せぬエラーが発生しました' });
    }
  });

  const { errors } = form.formState;

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 py-8">
      <div className="w-full max-w-md space-y-6 rounded-lg border bg-card p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">ユーザー登録</h1>
          <p className="text-sm text-muted-foreground">
            メールアドレスとパスワードを登録し、CV管理をはじめましょう。
          </p>
        </div>

        <form className="space-y-5" onSubmit={onSubmit} noValidate>
          <div className="space-y-2">
            <Label htmlFor="email">メールアドレス</Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              {...form.register('email')}
            />
            {errors.email && <p className="text-sm text-destructive">{errors.email.message}</p>}
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">パスワード</Label>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              {...form.register('password')}
            />
            <p className="text-xs text-muted-foreground">8文字以上で、英字と数字を含めてください。</p>
            {errors.password && <p className="text-sm text-destructive">{errors.password.message}</p>}
          </div>

          <div className="space-y-2">
            <Label htmlFor="passwordConfirmation">パスワード確認</Label>
            <Input
              id="passwordConfirmation"
              type="password"
              autoComplete="new-password"
              {...form.register('passwordConfirmation')}
            />
            {errors.passwordConfirmation && (
              <p className="text-sm text-destructive">{errors.passwordConfirmation.message}</p>
            )}
          </div>

          {errors.root?.message && <p className="text-sm text-destructive">{errors.root.message}</p>}

          <Button className="w-full" disabled={registerMutation.isPending} type="submit">
            {registerMutation.isPending ? '送信中…' : '登録メールを送信'}
          </Button>
        </form>

        {success && (
          <div className="rounded-md border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
            <p className="font-medium">{success.message}</p>
            <p className="mt-1">
              トークンの有効期限: {new Date(success.expiresAt).toLocaleString('ja-JP')}
            </p>
          </div>
        )}
      </div>
    </div>
  );
};
