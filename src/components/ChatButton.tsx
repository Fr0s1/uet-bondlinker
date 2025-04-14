
import React from 'react';
import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router';
import { MessageSquare } from 'lucide-react';
import { useCreateDirectMessage } from '@/hooks/use-messages';
import { toast } from '@/components/ui/use-toast';

interface ChatButtonProps {
  username: string;
  userId: string;
}

const ChatButton = ({ username, userId }: ChatButtonProps) => {
  const navigate = useNavigate();
  const createDirectMessage = useCreateDirectMessage();

  const handleStartChat = async () => {
    try {
      const conversation = await createDirectMessage.mutateAsync({
        username,
        content: 'Hello', // No initial message
      });

      if (conversation && conversation.id) {
        navigate(`/messages`);
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to start conversation. Please try again.",
        variant: "destructive",
      });
    }
  };

  return (
    <Button
      variant="outline"
      className="flex items-center space-x-1"
      onClick={handleStartChat}
      disabled={createDirectMessage.isPending}
    >
      <MessageSquare className="h-4 w-4" />
      <span>Message</span>
    </Button>
  );
};

export default ChatButton;
