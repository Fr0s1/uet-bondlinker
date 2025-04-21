import { formatDate } from "date-fns";

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
  continiousTop?: boolean;
  continiousBottom?: boolean;
  showTimeSeparator?: boolean;
}

const MessageBubble = ({ message, isCurrentUser, continiousTop, continiousBottom, showTimeSeparator }: MessageBubbleProps) => {
  return <>
    {showTimeSeparator && <div className="text-xs text-gray-500 text-center my-2">{formatDate(message.createdAt, 'HH:mm')}</div>}
    <div className={`flex ${isCurrentUser ? 'justify-end' : 'justify-start'} ${continiousTop ? '!mt-[1px]' : ''}`} title={message.createdAt}>
      <div
        className={`max-w-[75%] rounded-[18px] py-2 px-3 ${isCurrentUser
          ? 'bg-social-blue text-white'
          : ('bg-white border border-gray-200')
          } ${continiousTop ? (isCurrentUser ? '!rounded-tr-[4px]' : '!rounded-tl-[4px]') : ''}
            ${continiousBottom ? (isCurrentUser ? '!rounded-br-[4px]' : '!rounded-bl-[4px]') : ''}
          `}
      >
        <p className={`${isCurrentUser ? 'text-white' : 'text-gray-800'} text-left`} style={{ fontSize: '15px' }}>
          {message.content}
        </p>
      </div>
    </div>
  </>;
};

export default MessageBubble;
