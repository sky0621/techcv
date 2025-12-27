import { getApp, getApps, initializeApp } from 'firebase/app';
import {
  GoogleAuthProvider,
  getAuth,
  onIdTokenChanged,
  signInWithPopup,
  signOut,
  type Unsubscribe,
  type User
} from 'firebase/auth';

import { appConfig } from '@/lib/env';

const app = getApps().length ? getApp() : initializeApp(appConfig.firebase);
const auth = getAuth(app);

const googleProvider = new GoogleAuthProvider();
googleProvider.setCustomParameters({ prompt: 'select_account' });

export const signInWithGooglePopup = () => signInWithPopup(auth, googleProvider);

export const signOutFromFirebase = () => signOut(auth);

export const onFirebaseIdTokenChanged = (handler: (user: User | null) => void): Unsubscribe =>
  onIdTokenChanged(auth, handler);
