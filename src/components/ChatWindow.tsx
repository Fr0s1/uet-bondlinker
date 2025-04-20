import React, { useState, useRef, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loader2, Send } from 'lucide-react';
import MessageBubble from './MessageBubble';
import { Message, useConversation, useMessages } from '@/hooks/use-messages';
import { useWebSocket } from '@/hooks/use-websocket';
import { useDebounce } from '@/hooks/use-debounce';
import { da } from 'date-fns/locale';

interface ChatWindowProps {
  conversationId: string;
}

const ChatWindow = ({ conversationId }: ChatWindowProps) => {
  const { user } = useAuth();
  const [newMessage, setNewMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [isTypingTimeoutID, setIsTypingTimeoutID] = useState(0);
  const [recipientTyping, setRecipientTyping] = useState(false);
  const [recipientTypingTimeoutID, setRecipientTypingTimeoutID] = useState(0);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { sendMessage: sendWsMessage, ws } = useWebSocket();
  const debouncedTyping = useDebounce(isTyping, 1000);

  // Fetch conversation with proper typing
  const { data: conversation, isLoading: isConversationLoading, sendMessage, markAsRead } = useConversation(conversationId)

  // Fetch messages with proper typing
  const { messages, isLoading: isMessagesLoading, appendMessage } = useMessages(conversationId)

  useEffect(() => {
    if (!ws || !conversation) return;

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log(data);

      if (data.payload?.conversationId !== conversationId) {
        return
      }

      if (data.type === 'typing') {
        const typingEvent = data.payload;
        if (typingEvent.userId !== user?.id) {
          setRecipientTyping(typingEvent.isTyping);
          clearTimeout(recipientTypingTimeoutID)
          setRecipientTypingTimeoutID(window.setTimeout(() => {
            setRecipientTyping(false)
          }, 5000))
        }
        return
      }

      if (data.type === 'message') {
        setRecipientTyping(false);
        clearTimeout(recipientTypingTimeoutID)
        appendMessage(data.payload as Message)
      }
    };
  }, [ws, conversation, user]);

  // Send typing indicator
  useEffect(() => {
    if (!conversation || !debouncedTyping || !newMessage.trim()) return;

    sendWsMessage(conversation.recipient.id, 'typing', {
      conversationId,
      isTyping: debouncedTyping,
      userId: user?.id,
    });
  }, [debouncedTyping, conversation]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value);
    if (!e.target.value) {
      return
    }
    setIsTyping(true);
    clearTimeout(isTypingTimeoutID);
    setIsTypingTimeoutID(window.setTimeout(() => {
      setIsTyping(false)
    }, 5000))
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();

    if (!newMessage.trim()) return;

    sendMessage.mutate(newMessage, {
      onSuccess: () => {
        setIsTyping(false);
        clearTimeout(isTypingTimeoutID);
        scrollToBottom();
      }
    });
    setNewMessage('');
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, recipientTyping]);

  useEffect(() => {
    if (conversation) {
      markAsRead.mutate()
    }
  }, [conversation])

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
            src={conversation.recipient.avatar}
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
            {recipientTyping && (
              <div className="flex items-center text-gray-500 text-sm">
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                {conversation?.recipient.name} is typing...
              </div>
            )}
            <div className="!mt-0" ref={messagesEndRef} />
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
          onChange={handleInputChange}
          className="flex-grow"
        />
        <Button
          type="submit"
          disabled={sendMessage.isPending || !newMessage.trim()}
        >
          {sendMessage.isPending ? (
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
