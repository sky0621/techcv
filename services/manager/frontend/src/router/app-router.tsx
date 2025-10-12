import {
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter
} from '@tanstack/react-router';
import { Outlet } from '@tanstack/react-router';

import { RootLayout } from '@/components/layout/root-layout';
import { ProtectedRoute } from '@/features/auth/components/protected-route';
import { AuthCallbackPage } from '@/features/auth/pages/auth-callback-page';
import { LoginPage } from '@/features/auth/pages/login-page';
import { DashboardPage } from '@/features/dashboard/pages/dashboard-page';
import { CVEditorPage } from '@/features/cv/pages/cv-editor-page';
import { CVPreviewPage } from '@/features/cv/pages/cv-preview-page';
import { PublicCVViewPage } from '@/features/cv/pages/public-cv-view-page';
import { PublicURLManagerPage } from '@/features/cv/pages/public-url-manager-page';
import { NotFoundPage } from '@/router/not-found-page';

const rootRoute = createRootRoute({
  component: RootLayout
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'login',
  component: LoginPage
});

const authCallbackRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'auth/callback',
  component: AuthCallbackPage
});

const protectedLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'protected-layout',
  component: () => (
    <ProtectedRoute>
      <Outlet />
    </ProtectedRoute>
  )
});

const dashboardRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/',
  component: DashboardPage
});

const cvEditorRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: 'cv/edit',
  component: CVEditorPage
});

const cvPreviewRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: 'cv/preview',
  component: CVPreviewPage
});

const publicUrlManagerRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: 'settings/public-url',
  component: PublicURLManagerPage
});

const publicCvRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'public/$publicId',
  component: PublicCVViewPage
});

const notFoundRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '*',
  component: NotFoundPage
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  authCallbackRoute,
  protectedLayoutRoute.addChildren([
    dashboardRoute,
    cvEditorRoute,
    cvPreviewRoute,
    publicUrlManagerRoute
  ]),
  publicCvRoute,
  notFoundRoute
]);

const router = createRouter({ routeTree });

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

export const AppRouter = (): JSX.Element => {
  return <RouterProvider router={router} />;
};
