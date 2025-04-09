
import React from 'react';
import { formatDistanceToNow } from 'date-fns';

interface Message {
  id: string;
  content: string;
  createdAt: string;
  senderId: string;
  recipientId: string;
}

interface MessageBubbleProps {
  message: Message;
  isCurrentUser: boolean;
}

const MessageBubble = ({ message, isCurrentUser }: MessageBubbleProps) => {
  return (
    <div className={`flex ${isCurrentUser ? 'justify-end' : 'justify-start'}`}>
      <div 
        className={`max-w-[75%] rounded-xl p-3 ${
          isCurrentUser 
            ? 'bg-social-blue text-white rounded-tr-none' 
            : 'bg-white border border-gray-200 rounded-tl-none'
        }`}
      >
        <p className={`text-sm ${isCurrentUser ? 'text-white' : 'text-gray-800'} text-left`}>
          {message.content}
        </p>
        <div className={`text-xs mt-1 ${isCurrentUser ? 'text-blue-100' : 'text-gray-500'} text-left`}>
          {formatDistanceToNow(new Date(message.createdAt), { addSuffix: true })}
        </div>
      </div>
    </div>
  );
};

export default MessageBubble;
