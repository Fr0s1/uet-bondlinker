
import { useNavigate, useParams } from 'react-router';
import ConversationList from '@/components/ConversationList';
import ChatWindow from '@/components/ChatWindow';
import { useAuth } from '@/contexts/AuthContext';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useWebSocket } from '@/hooks/use-websocket';
import { useQueryClient } from '@tanstack/react-query';
import { appendMessage, Message, updateLastMessage, useConversation, useConversations } from '@/hooks/use-messages';

const Messages = () => {
  const queryClient = useQueryClient()
  const { conversationId } = useParams<{ conversationId?: string }>();
  const { data: conversation } = useConversation(conversationId)
  const { user } = useAuth();
  const { conversations, isLoading } = useConversations(user)
  const [recipientTyping, setRecipientTyping] = useState(false);
  const [recipientTypingTimeoutID, setRecipientTypingTimeoutID] = useState(0);
  const navigate = useNavigate()
  const { ws, sendMessage: sendWsMessage } = useWebSocket();

  const onSelectConversation = useCallback((conversationId: string | null) => {
    navigate(`/messages/c/${conversationId}`)
  }, [])

  const onIsTypingChange = useCallback((isTyping: boolean) => {
    if (!isTyping || !conversation) {
      return
    }
    sendWsMessage(conversation.recipient.id, 'typing', {
      conversationId,
      isTyping: true,
      userId: user?.id,
    })
  }, [conversation]);

  useEffect(() => {
    if (!ws || !user) return;

    const handleEvent = (event) => {
      const data = JSON.parse(event.data);
      console.log(data);

      if (data.type === 'message') {
        appendMessage(queryClient, data.payload?.conversationId, data.payload as Message)
        setRecipientTyping(false);
        const index = conversations.findIndex(it => it.id === data.payload?.conversationId)
        if (index == -1) {
          queryClient.invalidateQueries({ queryKey: ['conversations'] })
        } else {
          updateLastMessage(queryClient, data.payload?.conversationId, data.payload as Message, data.payload?.conversationId == conversationId)
        }
      }

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
    }

    ws.addEventListener('message', handleEvent)

    return () => {
      console.log('cleaning useEffect removeEventListener');
      ws.removeEventListener('message', handleEvent)
    }
  }, [ws, conversationId, user, conversations.length])

  if (!user) {
    return (
      <div className="bg-gray-50">
        <div className="container py-12">
          <div className="bg-white rounded-xl p-8 text-left card-shadow">
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Please login</h2>
            <p className="text-gray-600">You need to be logged in to access messages.</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 py-4">
      {/* Left Sidebar - Conversation List */}
      <aside className="lg:col-span-3">
        <ConversationList
          conversations={conversations}
          isLoading={isLoading}
          selectedConversationId={conversationId}
          onSelectConversation={onSelectConversation}
        />
      </aside>

      {/* Main Content - Chat Window */}
      <main className="lg:col-span-9">
        {conversationId ? (
          <ChatWindow conversationId={conversationId} recipientTyping={recipientTyping} onIsTypingChange={onIsTypingChange} />
        ) : (
          <div className="bg-white rounded-xl p-8 text-center h-[calc(100vh-120px)] flex items-center justify-center card-shadow">
            <div>
              <h2 className="text-xl font-semibold text-gray-700 mb-2">Select a conversation</h2>
              <p className="text-gray-500">Choose a conversation from the list to start chatting</p>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default Messages;
