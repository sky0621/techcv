import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react';

import { LoginPage } from '../login-page';

jest.mock('@tanstack/react-router', () => ({
  useNavigate: () => jest.fn(),
  useRouterState: () => ({ location: { search: {} } })
}));

jest.mock('@/features/auth/hooks/use-auth', () => ({
  useAuth: jest.fn()
}));

describe('LoginPage', () => {
  const signInWithGoogle = jest.fn();

  beforeEach(() => {
    jest.resetAllMocks();
    signInWithGoogle.mockResolvedValue(undefined);
    const { useAuth } = jest.requireMock('@/features/auth/hooks/use-auth');
    useAuth.mockReturnValue({
      session: { status: 'unauthenticated' },
      signInWithGoogle
    });
  });

  it('renders login headline', () => {
    render(<LoginPage />);

    expect(
      screen.getByRole('heading', { name: 'CV管理システムにサインイン' })
    ).toBeInTheDocument();
  });

  it('invokes Firebase Google sign-in when clicking the button', async () => {
    const user = userEvent.setup();
    render(<LoginPage />);

    await act(async () => {
      await user.click(screen.getByRole('button', { name: 'Googleでサインイン' }));
    });

    expect(signInWithGoogle).toHaveBeenCalledTimes(1);
  });
});
