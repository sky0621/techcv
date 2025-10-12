# CV Management Frontend

This is the frontend scaffold for the CV management system described in `.kiro/specs/cv-management-system/design.md`. The stack follows the design document: React 18 with TypeScript, TanStack Router, Jotai, shadcn/ui (Tailwind CSS), ky, and Jest.

## Getting Started

```bash
npm install
npm run dev
```

The development server runs on <http://localhost:5173>. Override the backend API endpoint with `VITE_API_BASE_URL` (default: `http://localhost:8080`).

## Available Scripts

- `npm run dev`: start the Vite dev server
- `npm run build`: type-check and build for production
- `npm run preview`: preview the production build locally
- `npm run lint`: run ESLint
- `npm run test`: run Jest with React Testing Library

## Project Layout

```
src/
  components/   reusable UI building blocks and layouts
  features/     feature-oriented modules (auth, cv, dashboard)
  providers/    application-level providers (state, react-query)
  router/       routing configuration with TanStack Router
  styles/       Tailwind globals and tokens
  lib/          shared helpers such as the ky client
```

Page-level components for authentication, dashboard, CV editing, previews, and public URL management are included as placeholders to guide future implementation.
