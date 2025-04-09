
import React, { useState } from 'react';
import { Bell } from 'lucide-react';
import { 
  Sheet, 
  SheetContent, 
  SheetHeader, 
  SheetTitle, 
  SheetTrigger 
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { useNotifications } from '@/hooks/use-notifications';
import NotificationItem from '@/components/NotificationItem';

const NotificationCenter: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { 
    notifications, 
    unreadCount, 
    isLoading, 
    markAsRead, 
    markAllAsRead 
  } = useNotifications();

  const handleMarkAsRead = (id: string) => {
    markAsRead(id);
  };

  const handleMarkAllAsRead = () => {
    markAllAsRead();
  };

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <span className="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] text-white">
              {unreadCount > 99 ? '99+' : unreadCount}
            </span>
          )}
        </Button>
      </SheetTrigger>
      <SheetContent side="right" className="w-full sm:max-w-md p-0">
        <SheetHeader className="px-4 py-3 border-b">
          <div className="flex justify-between items-center">
            <SheetTitle className="text-left">Notifications</SheetTitle>
            {(notifications && notifications.length > 0) && (
              <Button 
                variant="ghost" 
                size="sm"
                onClick={handleMarkAllAsRead}
              >
                Mark all as read
              </Button>
            )}
          </div>
        </SheetHeader>
        <div className="overflow-y-auto h-[calc(100vh-70px)]">
          {isLoading ? (
            <div className="flex justify-center items-center h-32">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-social-blue"></div>
            </div>
          ) : notifications && notifications.length > 0 ? (
            notifications.map((notification) => (
              <NotificationItem 
                key={notification.id} 
                notification={notification} 
                onRead={handleMarkAsRead}
              />
            ))
          ) : (
            <div className="flex flex-col justify-center items-center h-32 p-4 text-center text-gray-500">
              <Bell className="h-8 w-8 mb-2 text-gray-400" />
              <p>No notifications yet</p>
            </div>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default NotificationCenter;
