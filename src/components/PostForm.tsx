
import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Image, Smile, MapPin, Loader2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/AuthContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';

const PostForm = ({ onPostCreated }: { onPostCreated?: () => void }) => {
  const [content, setContent] = useState('');
  const { toast } = useToast();
  const { user, isAuthenticated } = useAuth();
  const queryClient = useQueryClient();
  
  const createPostMutation = useMutation({
    mutationFn: (postData: { content: string }) => 
      api.post('/posts', postData),
    onSuccess: () => {
      setContent('');
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      
      toast({
        title: "Post created",
        description: "Your post has been shared successfully.",
      });
      
      if (onPostCreated) {
        onPostCreated();
      }
    },
    onError: (error) => {
      console.error("Failed to create post:", error);
      toast({
        title: "Error creating post",
        description: "Please try again later.",
        variant: "destructive",
      });
    }
  });
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!isAuthenticated) {
      toast({
        title: "Authentication required",
        description: "Please log in to create a post",
        variant: "destructive",
      });
      return;
    }
    
    if (!content.trim()) {
      toast({
        title: "Cannot create empty post",
        description: "Please write something first.",
        variant: "destructive",
      });
      return;
    }
    
    createPostMutation.mutate({ content: content.trim() });
  };
  
  if (!isAuthenticated) {
    return (
      <div className="bg-white rounded-xl p-4 card-shadow mb-6 animate-fade-in">
        <div className="text-center py-4">
          <h3 className="font-medium text-lg mb-2">Join the conversation</h3>
          <p className="text-gray-500 mb-4">Sign in to share your thoughts and connect with others.</p>
          <div className="flex justify-center gap-3">
            <Button asChild variant="outline">
              <Link to="/login">Log in</Link>
            </Button>
            <Button asChild className="gradient-blue">
              <Link to="/register">Sign up</Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }
  
  return (
    <div className="bg-white rounded-xl p-4 card-shadow mb-6 animate-fade-in">
      <form onSubmit={handleSubmit}>
        <div className="flex items-start space-x-3">
          <Avatar className="h-10 w-10 mt-1">
            <AvatarImage src={user?.avatar || "/placeholder.svg"} alt={user?.name} />
            <AvatarFallback>{user?.name?.slice(0, 2).toUpperCase() || "U"}</AvatarFallback>
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
                disabled={createPostMutation.isPending || !content.trim()}
                className="px-5 gradient-blue"
              >
                {createPostMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Posting...
                  </>
                ) : (
                  'Post'
                )}
              </Button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default PostForm;
