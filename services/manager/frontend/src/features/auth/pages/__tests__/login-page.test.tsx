import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { LoginPage } from '../login-page';

describe('LoginPage', () => {
  it('renders login headline', () => {
    render(<LoginPage />);

    expect(
      screen.getByRole('heading', { name: 'CV管理システムにサインイン' })
    ).toBeInTheDocument();
  });

  it('redirects to backend login endpoint when clicking the button', async () => {
    const user = userEvent.setup();
    const originalLocation = window.location;
    delete (window as { location?: Location }).location;
    (window as { location: { href: string } }).location = { href: '' };

    render(<LoginPage />);

    await user.click(screen.getByRole('button', { name: 'Googleでサインイン' }));

    expect(window.location.href).toBe('http://localhost:8080/techcv/api/v1/auth/google/login');

    window.location = originalLocation;
  });
});
