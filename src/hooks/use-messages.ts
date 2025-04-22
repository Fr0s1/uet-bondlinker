import { useState } from 'react';
import { useQuery, useMutation, useQueryClient, QueryClient } from '@tanstack/react-query';
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

export const updateLastMessage = (queryClient: QueryClient, conversationId: string, message: Message, isRead: boolean) => {
  queryClient.setQueryData(['conversations'], (oldConversations: Conversation[]) => {
    return oldConversations.map(it => {
      if (it.id != conversationId) {
        return it
      }
      return {
        ...it,
        lastMessage: {
          content: message.content,
          createdAt: message.createdAt,
          isRead: isRead
        }
      }
    }).sort((a, b) => {
      return new Date(b.lastMessage?.createdAt || 0).getTime() - new Date(a.lastMessage?.createdAt || 0).getTime()
    })
  })
}

export const markAsReadLocal = (queryClient: QueryClient, conversationId: string) => {
  queryClient.setQueryData(['conversations'], (oldConversations: Conversation[]) => {
    return oldConversations.map(it => {
      if (it.id != conversationId) {
        return it
      }
      return {
        ...it,
        lastMessage: {
          ...it.lastMessage,
          isRead: true
        }
      }
    })
  })
}

export const appendMessage = (queryClient: QueryClient, conversationId: string, message: Message) => {
  queryClient.setQueryData(['messages', conversationId, 1], (oldMessages: Message[] | undefined) => {
    return [...(oldMessages || []), message]
  })
}

export const useConversations = (user?: User) => {
  const { data: conversations, isLoading, error, refetch } = useQuery<Conversation[]>({
    queryKey: ['conversations'],
    queryFn: () => api.get<Conversation[]>('/conversations'),
    enabled: !!user
  });

  const createConversation = useMutation<Conversation, Error, string>({
    mutationFn: (userId: string) => api.post<Conversation>('/conversations', { recipientId: userId }),
    onSuccess: (newConversation) => {
      refetch()
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
    refetch,
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

      updateLastMessage(queryClient, conversationId, newMessage, true)
    },
  });

  const markAsRead = useMutation<void, Error, void>({
    mutationFn: () => api.post<void>(`/conversations/${conversationId}/read`),
    onSuccess: () => {
      markAsReadLocal(queryClient, conversationId)
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
