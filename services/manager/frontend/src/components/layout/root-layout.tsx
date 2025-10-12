import { Outlet } from '@tanstack/react-router';
import { Suspense } from 'react';

import { LoadingScreen } from '@/components/system/loading-screen';

export const RootLayout = (): JSX.Element => (
  <div className="min-h-screen bg-background text-foreground">
    <Suspense fallback={<LoadingScreen message="初期化中..." />}>
      <Outlet />
    </Suspense>
  </div>
);
