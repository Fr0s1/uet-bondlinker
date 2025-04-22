
import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { requestNotificationPermission, onMessageListener } from '@/lib/firebase';
import { useToast } from '@/hooks/use-toast';
import { api } from '@/lib/api-client';

export const useFCM = () => {
  const { isAuthenticated } = useAuth();
  const { toast } = useToast();
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

  useEffect(() => {
    if (isAuthenticated) {
      const unsubscribe = onMessageListener()
        .then((payload: any) => {
          toast({
            title: payload.notification?.title || 'New Notification',
            description: payload.notification?.body,
          });
        })
        .catch((err) => console.error('FCM message error:', err));

      return () => {
        if (typeof unsubscribe === 'function') {
          unsubscribe();
        }
      };
    }
  }, [isAuthenticated]);

  return { fcmToken };
};
