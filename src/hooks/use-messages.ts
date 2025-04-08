
import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { toast } from '@/components/ui/use-toast';

export interface Message {
  id: string;
  content: string;
  createdAt: string;
  senderId: string;
  recipientId: string;
}

export interface Conversation {
  id: string;
  recipient: {
    id: string;
    name: string;
    username: string;
    avatar: string | null;
  };
  lastMessage: {
    content: string;
    createdAt: string;
    isRead: boolean;
  } | null;
}

export const useConversations = () => {
  const queryClient = useQueryClient();
  
  const { data: conversations, isLoading, error } = useQuery({
    queryKey: ['conversations'],
    queryFn: () => api.get<Conversation[]>('/conversations'),
  });
  
  const createConversation = useMutation({
    mutationFn: (userId: string) => api.post<Conversation>('/conversations', { recipientId: userId }),
    onSuccess: (newConversation) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      toast({
        title: "Conversation created",
        description: "You can now start messaging this user.",
      });
      return newConversation;
    },
  });
  
  return {
    conversations: conversations || [],
    isLoading,
    error,
    createConversation,
  };
};

export const useConversation = (conversationId: string) => {
  return useQuery({
    queryKey: ['conversation', conversationId],
    queryFn: () => api.get<Conversation>(`/conversations/${conversationId}`),
    enabled: !!conversationId,
  });
};

export const useMessages = (conversationId: string, limit = 50) => {
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();
  
  const { data: messages, isLoading, error } = useQuery({
    queryKey: ['messages', conversationId, page, limit],
    queryFn: () => api.get<Message[]>(
      `/conversations/${conversationId}/messages?limit=${limit}&offset=${(page - 1) * limit}`
    ),
    enabled: !!conversationId,
  });
  
  const sendMessage = useMutation({
    mutationFn: (content: string) => 
      api.post<Message>(`/conversations/${conversationId}/messages`, { content }),
    onSuccess: (newMessage) => {
      queryClient.invalidateQueries({ queryKey: ['messages', conversationId] });
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
    },
  });
  
  const markAsRead = useMutation({
    mutationFn: () => api.post<void>(`/conversations/${conversationId}/read`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
    },
  });
  
  return {
    messages: messages || [],
    isLoading,
    error,
    page,
    setPage,
    sendMessage,
    markAsRead,
  };
};

export const useCreateDirectMessage = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      username, 
      content 
    }: { 
      username: string; 
      content: string;
    }) => {
      // First get the user by username
      const user = await api.get<{ id: string }>(`/users/username/${username}`);
      
      if (!user || !user.id) {
        throw new Error("User not found");
      }
      
      // Then create a conversation with that user
      const conversation = await api.post<Conversation>('/conversations', { 
        recipientId: user.id 
      });
      
      if (!conversation || !conversation.id) {
        throw new Error("Failed to create conversation");
      }
      
      // Then send the message
      await api.post<Message>(`/conversations/${conversation.id}/messages`, { 
        content 
      });
      
      return conversation;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      toast({
        title: "Message sent",
        description: "Your direct message has been sent successfully.",
      });
    },
  });
};
