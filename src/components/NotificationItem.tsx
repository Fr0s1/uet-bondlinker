
import React from 'react';
import { Link } from 'react-router';
import { formatDistanceToNow } from 'date-fns';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Notification } from '@/hooks/use-notifications';
import { UserPlus, Heart, MessageCircle, Share2, Bell } from 'lucide-react';

interface NotificationItemProps {
  notification: Notification;
  onRead: (id: string) => void;
}

const NotificationItem: React.FC<NotificationItemProps> = ({ notification, onRead }) => {
  const getNotificationIcon = () => {
    switch (notification.type) {
      case 'follow':
        return <UserPlus className="h-4 w-4 text-blue-500" />;
      case 'like':
        return <Heart className="h-4 w-4 text-red-500" />;
      case 'comment':
        return <MessageCircle className="h-4 w-4 text-green-500" />;
      case 'share':
        return <Share2 className="h-4 w-4 text-purple-500" />;
      case 'message':
        return <MessageCircle className="h-4 w-4 text-blue-500" />;
      default:
        return <Bell className="h-4 w-4 text-gray-500" />;
    }
  };

  const getNotificationLink = () => {
    switch (notification.type) {
      case 'follow':
        return notification.senderId ? `/profile/${notification.senderId}` : '#';
      case 'like':
      case 'comment':
      case 'share':
        return notification.relatedEntityId ? `/posts/${notification.relatedEntityId}` : '#';
      case 'message':
        return notification.relatedEntityId ? `/messages?conversation=${notification.relatedEntityId}` : '/messages';
      default:
        return '#';
    }
  };

  const handleClick = () => {
    if (!notification.isRead) {
      onRead(notification.id);
    }
  };

  return (
    <Link
      to={getNotificationLink()}
      className={`block p-4 border-b hover:bg-gray-50 transition-colors ${!notification.isRead ? 'bg-blue-50' : ''}`}
      onClick={handleClick}
    >
      <div className="flex items-start space-x-3">
        <div className="flex-shrink-0">
          {notification.sender ? (
            <Avatar className="h-10 w-10">
              <AvatarImage src={notification.sender.avatar} alt={notification.sender.name} />
              <AvatarFallback>{notification.sender.name.charAt(0)}</AvatarFallback>
            </Avatar>
          ) : (
            <div className="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
              {getNotificationIcon()}
            </div>
          )}
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm text-gray-900">{notification.message}</p>
          <p className="text-xs text-gray-500 mt-1">
            {formatDistanceToNow(new Date(notification.createdAt), { addSuffix: true })}
          </p>
        </div>
        {!notification.isRead && (
          <div className="flex-shrink-0">
            <div className="h-2 w-2 rounded-full bg-blue-500"></div>
          </div>
        )}
      </div>
    </Link>
  );
};

export default NotificationItem;
