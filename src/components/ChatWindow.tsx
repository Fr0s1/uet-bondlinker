
import React, { useState, useRef, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/contexts/AuthContext';
import { api } from '@/lib/api-client';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loader2, Send } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import MessageBubble from './MessageBubble';

interface Message {
  id: string;
  content: string;
  createdAt: string;
  senderId: string;
  recipientId: string;
}

interface ChatWindowProps {
  conversationId: string;
}

const ChatWindow = ({ conversationId }: ChatWindowProps) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [newMessage, setNewMessage] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // Fetch conversation
  const { data: conversation, isLoading: isConversationLoading } = useQuery({
    queryKey: ['conversation', conversationId],
    queryFn: async () => {
      try {
        // Mock data - in reality, you'd have a backend endpoint
        // const response = await api.get(`/conversations/${conversationId}`);
        // return response;
        
        // Mock data
        return {
          id: conversationId,
          recipient: {
            id: '2',
            name: 'Jane Smith',
            username: 'janesmith',
            avatar: null
          }
        };
      } catch (error) {
        console.error('Error fetching conversation:', error);
        throw error;
      }
    },
    enabled: !!conversationId
  });
  
  // Fetch messages
  const { data: messages, isLoading: isMessagesLoading } = useQuery({
    queryKey: ['messages', conversationId],
    queryFn: async () => {
      try {
        // Mock data - in reality, you'd have a backend endpoint
        // const response = await api.get(`/conversations/${conversationId}/messages`);
        // return response;
        
        // Mock messages
        const mockMessages = [
          {
            id: '1',
            content: 'Hey, how are you doing?',
            createdAt: new Date(Date.now() - 1000 * 60 * 60).toISOString(), // 1 hour ago
            senderId: '2', // from recipient
            recipientId: user?.id || '',
          },
          {
            id: '2',
            content: 'I\'m good, thanks for asking. How about you?',
            createdAt: new Date(Date.now() - 1000 * 60 * 55).toISOString(), // 55 minutes ago
            senderId: user?.id || '',
            recipientId: '2',
          },
          {
            id: '3',
            content: 'Not bad! Just working on some new features for the app.',
            createdAt: new Date(Date.now() - 1000 * 60 * 50).toISOString(), // 50 minutes ago
            senderId: '2',
            recipientId: user?.id || '',
          },
          {
            id: '4',
            content: 'That sounds interesting! What kind of features?',
            createdAt: new Date(Date.now() - 1000 * 60 * 45).toISOString(), // 45 minutes ago
            senderId: user?.id || '',
            recipientId: '2',
          },
          {
            id: '5',
            content: 'I\'m working on a messaging system, actually!',
            createdAt: new Date(Date.now() - 1000 * 60 * 40).toISOString(), // 40 minutes ago
            senderId: '2',
            recipientId: user?.id || '',
          }
        ];
        
        return mockMessages;
      } catch (error) {
        console.error('Error fetching messages:', error);
        throw error;
      }
    },
    enabled: !!conversationId && !!user
  });
  
  // Send message mutation
  const sendMessageMutation = useMutation({
    mutationFn: async (content: string) => {
      try {
        // In reality, you'd have a backend endpoint
        // return api.post(`/conversations/${conversationId}/messages`, { content });
        
        // Mock sending a message
        console.log('Sending message:', content);
        return {
          id: Date.now().toString(),
          content,
          createdAt: new Date().toISOString(),
          senderId: user?.id || '',
          recipientId: conversation?.recipient.id || '',
        };
      } catch (error) {
        console.error('Error sending message:', error);
        throw error;
      }
    },
    onSuccess: (newMessage) => {
      // Update the messages cache
      queryClient.setQueryData(['messages', conversationId], (oldData: Message[] | undefined) => {
        return oldData ? [...oldData, newMessage] : [newMessage];
      });
      
      // Scroll to bottom
      scrollToBottom();
    }
  });
  
  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newMessage.trim()) return;
    
    sendMessageMutation.mutate(newMessage);
    setNewMessage('');
  };
  
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };
  
  useEffect(() => {
    scrollToBottom();
  }, [messages]);
  
  if (isConversationLoading || isMessagesLoading) {
    return (
      <div className="bg-white rounded-xl p-4 h-[calc(100vh-120px)] flex items-center justify-center card-shadow">
        <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
        <span className="ml-2 text-gray-500">Loading conversation...</span>
      </div>
    );
  }
  
  if (!conversation) {
    return (
      <div className="bg-white rounded-xl p-4 h-[calc(100vh-120px)] flex items-center justify-center card-shadow">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-700">Conversation not found</h3>
          <p className="text-gray-500 text-sm mt-1">The conversation you're looking for doesn't exist</p>
        </div>
      </div>
    );
  }
  
  return (
    <div className="bg-white rounded-xl overflow-hidden card-shadow h-[calc(100vh-120px)] flex flex-col">
      {/* Chat header */}
      <div className="p-4 border-b flex items-center">
        <Avatar className="h-10 w-10 mr-3">
          <AvatarImage 
            src={conversation.recipient.avatar || "/placeholder.svg"} 
            alt={conversation.recipient.name} 
          />
          <AvatarFallback>
            {conversation.recipient.name.substring(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        <div>
          <h2 className="font-semibold text-gray-800">{conversation.recipient.name}</h2>
          <p className="text-xs text-gray-500">@{conversation.recipient.username}</p>
        </div>
      </div>
      
      {/* Messages area */}
      <div className="flex-grow overflow-y-auto p-4 bg-gray-50">
        {messages && messages.length > 0 ? (
          <div className="space-y-3">
            {messages.map((message) => (
              <MessageBubble 
                key={message.id}
                message={message}
                isCurrentUser={message.senderId === user?.id}
              />
            ))}
            <div ref={messagesEndRef} />
          </div>
        ) : (
          <div className="h-full flex items-center justify-center">
            <p className="text-gray-500">No messages yet. Start the conversation!</p>
          </div>
        )}
      </div>
      
      {/* Message input */}
      <form onSubmit={handleSendMessage} className="p-3 border-t flex gap-2">
        <Input
          type="text"
          placeholder="Type a message..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          className="flex-grow"
        />
        <Button 
          type="submit" 
          disabled={sendMessageMutation.isPending || !newMessage.trim()}
        >
          {sendMessageMutation.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Send className="h-4 w-4" />
          )}
          <span className="sr-only">Send</span>
        </Button>
      </form>
    </div>
  );
};

export default ChatWindow;
