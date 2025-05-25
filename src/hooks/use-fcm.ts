
import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { requestNotificationPermission, onMessageListener } from '@/lib/firebase';
import { api } from '@/lib/api-client';

export const useFCM = () => {
  const { isAuthenticated } = useAuth();
  const [fcmToken, setFcmToken] = useState<string | null>(null);

  const saveFCMToken = async (token: string) => {
    try {
      await api.post('/users/fcm-token', {
        token,
        device: navigator.userAgent
      });
    } catch (error) {
      console.error('Error saving FCM token:', error);
    }
  };

  useEffect(() => {
    if (isAuthenticated && 'serviceWorker' in navigator) {
      requestNotificationPermission().then((token) => {
        if (token) {
          setFcmToken(token);
          saveFCMToken(token);
        }
      });
    }
  }, [isAuthenticated]);

  return { fcmToken };
};
