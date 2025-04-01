
import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Image, Smile, MapPin } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

const PostForm = ({ onPostCreated }: { onPostCreated?: () => void }) => {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!content.trim()) {
      toast({
        title: "Cannot create empty post",
        description: "Please write something first.",
        variant: "destructive",
      });
      return;
    }
    
    setIsSubmitting(true);
    
    // Simulate API call
    setTimeout(() => {
      setContent('');
      setIsSubmitting(false);
      
      toast({
        title: "Post created",
        description: "Your post has been shared successfully.",
      });
      
      if (onPostCreated) {
        onPostCreated();
      }
    }, 1000);
  };
  
  return (
    <div className="bg-white rounded-xl p-4 card-shadow mb-6 animate-fade-in">
      <form onSubmit={handleSubmit}>
        <div className="flex items-start space-x-3">
          <Avatar className="h-10 w-10 mt-1">
            <AvatarImage src="/placeholder.svg" alt="User" />
            <AvatarFallback>US</AvatarFallback>
          </Avatar>
          
          <div className="flex-1">
            <Textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="What's on your mind?"
              className="min-h-20 resize-none border-none focus-visible:ring-0 focus-visible:ring-offset-0 p-0"
            />
            
            <div className="flex items-center justify-between mt-3 pt-3 border-t">
              <div className="flex space-x-2">
                <Button type="button" variant="ghost" size="icon" className="text-social-blue rounded-full h-9 w-9">
                  <Image className="h-5 w-5" />
                </Button>
                <Button type="button" variant="ghost" size="icon" className="text-social-blue rounded-full h-9 w-9">
                  <Smile className="h-5 w-5" />
                </Button>
                <Button type="button" variant="ghost" size="icon" className="text-social-blue rounded-full h-9 w-9">
                  <MapPin className="h-5 w-5" />
                </Button>
              </div>
              
              <Button 
                type="submit" 
                disabled={isSubmitting || !content.trim()}
                className="px-5 gradient-blue"
              >
                {isSubmitting ? 'Posting...' : 'Post'}
              </Button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default PostForm;
