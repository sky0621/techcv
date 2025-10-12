interface LoadingScreenProps {
  message?: string;
}

export const LoadingScreen = ({ message = 'Loading...' }: LoadingScreenProps): JSX.Element => (
  <div className="flex min-h-screen items-center justify-center bg-background text-foreground">
    <span className="text-sm font-medium text-muted-foreground">{message}</span>
  </div>
);
