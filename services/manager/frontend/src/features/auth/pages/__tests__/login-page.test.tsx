import { render, screen } from '@testing-library/react';

import { LoginPage } from '../login-page';

describe('LoginPage', () => {
  it('renders login headline', () => {
    render(<LoginPage />);

    expect(
      screen.getByRole('heading', { name: 'CV管理システムにサインイン' })
    ).toBeInTheDocument();
  });
});
