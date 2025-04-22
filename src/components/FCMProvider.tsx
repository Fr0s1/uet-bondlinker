
import { ReactNode } from 'react';
import { useFCM } from '@/hooks/use-fcm';

interface FCMProviderProps {
  children: ReactNode;
}

export const FCMProvider = ({ children }: FCMProviderProps) => {
  useFCM(); // Initialize FCM
  return <>{children}</>;
};
