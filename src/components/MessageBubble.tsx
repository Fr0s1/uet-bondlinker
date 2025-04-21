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
        className={`max-w-[75%] rounded-xl p-3 ${isCurrentUser
          ? 'bg-social-blue text-white rounded-tr-none'
          : 'bg-white border border-gray-200 rounded-tl-none'
          }`}
      >
        <p className={`text-sm ${isCurrentUser ? 'text-white' : 'text-gray-800'} text-left`}>
          {message.content}
        </p>
      </div>
    </div>
  );
};

export default MessageBubble;
