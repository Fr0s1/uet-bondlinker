
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { formatDistanceToNow } from 'date-fns';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

interface Conversation {
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

interface ConversationListProps {
  selectedConversationId: string | null;
  onSelectConversation: (id: string) => void;
}

const ConversationList = ({ selectedConversationId, onSelectConversation }: ConversationListProps) => {
  const { user } = useAuth();
  
  // Fetch conversations
  const { data: conversations, isLoading } = useQuery({
    queryKey: ['conversations'],
    queryFn: async () => {
      return api.get<Conversation[]>('/conversations');
    },
    enabled: !!user
  });
  
  if (isLoading) {
    return (
      <div className="bg-white rounded-xl p-4 h-[calc(100vh-120px)] flex items-center justify-center card-shadow">
        <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
        <span className="ml-2 text-gray-500">Loading conversations...</span>
      </div>
    );
  }
  
  if (!conversations || conversations.length === 0) {
    return (
      <div className="bg-white rounded-xl p-4 h-[calc(100vh-120px)] flex items-center justify-center card-shadow">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-700">No conversations</h3>
          <p className="text-gray-500 text-sm mt-1">Start a new conversation by searching for users</p>
        </div>
      </div>
    );
  }
  
  return (
    <div className="bg-white rounded-xl overflow-hidden card-shadow h-[calc(100vh-120px)] flex flex-col">
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold text-gray-800">Messages</h2>
      </div>
      
      <div className="overflow-y-auto flex-grow">
        {conversations.map((conversation) => (
          <div 
            key={conversation.id}
            className={`p-3 border-b cursor-pointer hover:bg-gray-50 transition-colors ${
              selectedConversationId === conversation.id ? 'bg-gray-100' : ''
            }`}
            onClick={() => onSelectConversation(conversation.id)}
          >
            <div className="flex items-center">
              <Avatar className="h-12 w-12 mr-3">
                <AvatarImage src={conversation.recipient.avatar || "/placeholder.svg"} alt={conversation.recipient.name} />
                <AvatarFallback>{conversation.recipient.name.substring(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>
              
              <div className="flex-grow min-w-0">
                <div className="flex justify-between items-center">
                  <h3 className="font-medium text-gray-900 truncate">{conversation.recipient.name}</h3>
                  {conversation.lastMessage && (
                    <span className="text-xs text-gray-500">
                      {formatDistanceToNow(new Date(conversation.lastMessage.createdAt), { addSuffix: true })}
                    </span>
                  )}
                </div>
                
                {conversation.lastMessage && (
                  <p className={`text-sm truncate ${
                    conversation.lastMessage.isRead ? 'text-gray-500' : 'text-gray-900 font-medium'
                  }`}>
                    {conversation.lastMessage.content}
                  </p>
                )}
              </div>
              
              {conversation.lastMessage && !conversation.lastMessage.isRead && (
                <div className="ml-2 w-2 h-2 bg-social-blue rounded-full"></div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ConversationList;
