import React from 'react';
import ReactDOM from 'react-dom/client';
import { AppProviders } from './providers/app-providers';
import { AppRouter } from './router/app-router';
import './styles/globals.css';

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found');
}

ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <AppProviders>
      <AppRouter />
    </AppProviders>
  </React.StrictMode>
);
