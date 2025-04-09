
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { useAuth } from '@/contexts/AuthContext';

export interface Notification {
  id: string;
  user_id: string;
  sender_id?: string;
  type: 'follow' | 'like' | 'comment' | 'share' | 'message' | 'system_alert';
  message: string;
  related_entity_id?: string;
  entity_type?: string;
  is_read: boolean;
  created_at: string;
  sender?: {
    id: string;
    name: string;
    username: string;
    avatar?: string;
  };
}

export interface NotificationFilter {
  limit?: number;
  offset?: number;
  is_read?: boolean;
}

export const useNotifications = (filter: NotificationFilter = { limit: 20, offset: 0 }) => {
  const { isAuthenticated } = useAuth();
  const queryClient = useQueryClient();

  const fetchNotifications = async (): Promise<Notification[]> => {
    const params = new URLSearchParams();
    if (filter.limit) params.append('limit', filter.limit.toString());
    if (filter.offset) params.append('offset', filter.offset.toString());
    if (filter.is_read !== undefined) params.append('is_read', filter.is_read.toString());
    
    return api.get<Notification[]>(`/notifications?${params.toString()}`);
  };

  const fetchUnreadCount = async (): Promise<{ count: number }> => {
    return api.get<{ count: number }>('/notifications/unread-count');
  };

  const markAsRead = async (id: string): Promise<void> => {
    return api.put(`/notifications/${id}/read`);
  };

  const markAllAsRead = async (): Promise<void> => {
    return api.put('/notifications/read-all');
  };

  const { data: notifications, isLoading, error, refetch } = useQuery({
    queryKey: ['notifications', filter],
    queryFn: fetchNotifications,
    enabled: isAuthenticated,
  });

  const { data: unreadCount } = useQuery({
    queryKey: ['notifications-unread-count'],
    queryFn: fetchUnreadCount,
    enabled: isAuthenticated,
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  const markAsReadMutation = useMutation({
    mutationFn: markAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      queryClient.invalidateQueries({ queryKey: ['notifications-unread-count'] });
    },
  });

  const markAllAsReadMutation = useMutation({
    mutationFn: markAllAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      queryClient.invalidateQueries({ queryKey: ['notifications-unread-count'] });
    },
  });

  return {
    notifications,
    unreadCount: unreadCount?.count || 0,
    isLoading,
    error,
    refetch,
    markAsRead: markAsReadMutation.mutate,
    markAllAsRead: markAllAsReadMutation.mutate,
  };
};
