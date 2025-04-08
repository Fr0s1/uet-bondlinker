
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from '@/components/Navbar';
import ConversationList from '@/components/ConversationList';
import ChatWindow from '@/components/ChatWindow';
import { useAuth } from '@/contexts/AuthContext';

const Messages = () => {
  const { conversationId } = useParams<{ conversationId?: string }>();
  const { user } = useAuth();
  const [selectedConversation, setSelectedConversation] = useState<string | null>(conversationId || null);
  
  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="container py-12">
          <div className="bg-white rounded-xl p-8 text-center card-shadow">
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Please login</h2>
            <p className="text-gray-600">You need to be logged in to access messages.</p>
          </div>
        </div>
      </div>
    );
  }
  
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="container py-4">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
          {/* Left Sidebar - Conversation List */}
          <aside className="lg:col-span-3">
            <ConversationList 
              selectedConversationId={selectedConversation} 
              onSelectConversation={setSelectedConversation} 
            />
          </aside>
          
          {/* Main Content - Chat Window */}
          <main className="lg:col-span-9">
            {selectedConversation ? (
              <ChatWindow conversationId={selectedConversation} />
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
      </div>
    </div>
  );
};

export default Messages;
