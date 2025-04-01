
import React, { useState } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Heart, Send } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

interface Comment {
  id: string;
  author: {
    name: string;
    avatar: string;
  };
  content: string;
  createdAt: string;
  likes: number;
  isLiked: boolean;
}

// Mock comments data
const mockComments: Comment[] = [
  {
    id: "c1",
    author: {
      name: "Jane Doe",
      avatar: "/placeholder.svg",
    },
    content: "This is a great post! Thanks for sharing this insight.",
    createdAt: "2023-05-15T10:30:00Z",
    likes: 5,
    isLiked: false,
  },
  {
    id: "c2",
    author: {
      name: "Bob Smith",
      avatar: "/placeholder.svg",
    },
    content: "I completely agree with your points. Very well said!",
    createdAt: "2023-05-15T12:45:00Z",
    likes: 3,
    isLiked: true,
  },
];

const CommentSection = ({ postId }: { postId: string }) => {
  const [comments, setComments] = useState<Comment[]>(mockComments);
  const [newComment, setNewComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();
  
  const handleSubmitComment = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newComment.trim()) return;
    
    setIsSubmitting(true);
    
    // Simulate API call
    setTimeout(() => {
      const comment: Comment = {
        id: `c${Date.now()}`,
        author: {
          name: "Current User",
          avatar: "/placeholder.svg",
        },
        content: newComment,
        createdAt: new Date().toISOString(),
        likes: 0,
        isLiked: false,
      };
      
      setComments([comment, ...comments]);
      setNewComment('');
      setIsSubmitting(false);
      
      toast({
        title: "Comment added",
        description: "Your comment has been posted successfully.",
      });
    }, 500);
  };
  
  const handleLikeComment = (commentId: string) => {
    setComments(comments.map(comment => {
      if (comment.id === commentId) {
        const newIsLiked = !comment.isLiked;
        return {
          ...comment,
          isLiked: newIsLiked,
          likes: newIsLiked ? comment.likes + 1 : comment.likes - 1,
        };
      }
      return comment;
    }));
  };
  
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };
  
  return (
    <div className="bg-gray-50 border-t p-4">
      <form onSubmit={handleSubmitComment} className="flex items-start space-x-3 mb-4">
        <Avatar className="h-8 w-8">
          <AvatarImage src="/placeholder.svg" alt="Current User" />
          <AvatarFallback>CU</AvatarFallback>
        </Avatar>
        
        <div className="flex-1 relative">
          <Textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Write a comment..."
            className="pr-10 min-h-12 resize-none bg-white"
          />
          <Button 
            type="submit" 
            size="icon" 
            className="absolute right-2 bottom-2 h-7 w-7 text-social-blue" 
            variant="ghost"
            disabled={isSubmitting || !newComment.trim()}
          >
            <Send className="h-4 w-4" />
          </Button>
        </div>
      </form>
      
      <div className="space-y-4">
        {comments.map((comment) => (
          <div key={comment.id} className="flex space-x-3">
            <Avatar className="h-8 w-8">
              <AvatarImage src={comment.author.avatar} alt={comment.author.name} />
              <AvatarFallback>{comment.author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
            </Avatar>
            
            <div className="flex-1">
              <div className="bg-white p-3 rounded-lg">
                <div className="flex justify-between items-start">
                  <h4 className="font-medium text-sm">{comment.author.name}</h4>
                  <span className="text-xs text-gray-500">{formatDate(comment.createdAt)}</span>
                </div>
                <p className="text-sm mt-1">{comment.content}</p>
              </div>
              
              <div className="flex items-center space-x-4 mt-1 ml-1">
                <button 
                  className={`text-xs flex items-center space-x-1 ${comment.isLiked ? 'text-social-blue' : 'text-gray-500'}`}
                  onClick={() => handleLikeComment(comment.id)}
                >
                  <Heart className="h-3 w-3" fill={comment.isLiked ? "currentColor" : "none"} />
                  <span>{comment.likes > 0 && comment.likes}</span>
                </button>
                <button className="text-xs text-gray-500">Reply</button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CommentSection;
