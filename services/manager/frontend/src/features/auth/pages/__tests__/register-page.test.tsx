import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { RegisterPage } from '../register-page';

jest.mock('@/features/auth/hooks/use-register', () => ({
  useRegisterUser: () => ({
    mutateAsync: jest.fn(),
    isPending: false
  })
}));

describe('RegisterPage', () => {
  it('shows validation error when email is invalid', async () => {
    const user = userEvent.setup();
    render(<RegisterPage />);

    await user.type(screen.getByLabelText('メールアドレス'), 'invalid-email');
    await user.type(screen.getByLabelText('パスワード'), 'Passw0rd');
    await user.type(screen.getByLabelText('パスワード確認'), 'Passw0rd');

    await user.click(screen.getByRole('button', { name: '登録メールを送信' }));

    expect(await screen.findByText('メールアドレスの形式が正しくありません')).toBeInTheDocument();
  });
});
