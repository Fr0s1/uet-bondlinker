import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { toast } from '@/components/ui/use-toast';
import { User } from '@/contexts/AuthContext';

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

export const useConversations = (user?: User) => {
  const queryClient = useQueryClient();

  const { data: conversations, isLoading, error } = useQuery<Conversation[]>({
    queryKey: ['conversations'],
    queryFn: () => api.get<Conversation[]>('/conversations'),
    enabled: !!user
  });

  const createConversation = useMutation<Conversation, Error, string>({
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

const messagesLimit = 50

export const useConversation = (conversationId: string) => {
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery<Conversation>({
    queryKey: ['conversation', conversationId],
    queryFn: () => api.get<Conversation>(`/conversations/${conversationId}`),
    enabled: !!conversationId,
  });

  const sendMessage = useMutation<Message, Error, string>({
    mutationFn: (content: string) =>
      api.post<Message>(`/conversations/${conversationId}/messages`, { content }),
    onSuccess: (newMessage) => {
      queryClient.setQueryData(['messages', conversationId, 1], (oldData: Message[] | undefined) => {
        return oldData ? [...oldData, newMessage] : [newMessage];
      });
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
    },
  });

  const markAsRead = useMutation<void, Error, void>({
    mutationFn: () => api.post<void>(`/conversations/${conversationId}/read`),
    onSuccess: () => {
    },
  });

  return {
    data,
    isLoading,
    error,
    markAsRead,
    sendMessage
  }
};

export const useMessages = (conversationId: string) => {
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();

  const { data: messages, isLoading, error } = useQuery<Message[]>({
    queryKey: ['messages', conversationId, page],
    queryFn: () => api.get<Message[]>(
      `/conversations/${conversationId}/messages?limit=${messagesLimit}&offset=${(page - 1) * messagesLimit}`
    ),
    enabled: !!conversationId,
  });


  return {
    messages: messages || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useCreateDirectMessage = () => {
  const queryClient = useQueryClient();

  interface CreateDirectMessageParams {
    username: string;
    content: string;
  }

  return useMutation<Conversation, Error, CreateDirectMessageParams>({
    mutationFn: async ({
      username,
      content
    }: CreateDirectMessageParams) => {
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

      if (!!content.trim()) {
        // Then send the message
        await api.post<Message>(`/conversations/${conversation.id}/messages`, {
          content: content.trim()
        });
      }

      return conversation;
    },
    onSuccess: (_, { content }) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      if (!!content.trim()) {
        toast({
          title: "Message sent",
          description: "Your direct message has been sent successfully.",
        });
      }
    },
  });
};
