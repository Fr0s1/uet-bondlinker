
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Heart, MessageSquare, Share, MoreHorizontal } from 'lucide-react';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import CommentSection from './CommentSection';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from '@/components/ui/use-toast';

export interface PostProps {
  id: string;
  author: {
    id: string;
    name: string;
    username: string;
    avatar: string;
  };
  content: string;
  image?: string;
  createdAt: string;
  likes: number;
  comments: number;
  shares: number;
  isLiked?: boolean;
  sharedPost?: {
    id: string;
    author: {
      id: string;
      name: string;
      username: string;
      avatar: string;
    };
    content: string;
    image?: string;
    createdAt: string;
  };
}

const Post = ({ 
  id, 
  author, 
  content, 
  image, 
  createdAt, 
  likes, 
  comments, 
  shares, 
  isLiked = false,
  sharedPost
}: PostProps) => {
  const [liked, setLiked] = useState(isLiked);
  const [likeCount, setLikeCount] = useState(likes);
  const [showComments, setShowComments] = useState(false);
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const [shareContent, setShareContent] = useState('');
  const { isAuthenticated, user } = useAuth();
  const queryClient = useQueryClient();
  
  const likePostMutation = useMutation({
    mutationFn: () => api.post(`/posts/${id}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      queryClient.invalidateQueries({ queryKey: ['trending'] });
    },
    onError: () => {
      // Revert optimistic update on error
      setLiked(!liked);
      setLikeCount(prev => liked ? prev + 1 : prev - 1);
    }
  });
  
  const unlikePostMutation = useMutation({
    mutationFn: () => api.delete(`/posts/${id}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      queryClient.invalidateQueries({ queryKey: ['trending'] });
    },
    onError: () => {
      // Revert optimistic update on error
      setLiked(!liked);
      setLikeCount(prev => liked ? prev - 1 : prev + 1);
    }
  });
  
  const sharePostMutation = useMutation({
    mutationFn: (content: string) => api.post(`/posts/${id}/share`, { content }),
    onSuccess: () => {
      setShareDialogOpen(false);
      setShareContent('');
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      queryClient.invalidateQueries({ queryKey: ['trending'] });
      toast({
        title: "Post shared",
        description: "The post has been shared successfully!",
      });
    },
    onError: () => {
      toast({
        title: "Error",
        description: "Failed to share the post. Please try again.",
        variant: "destructive",
      });
    }
  });
  
  const handleLike = () => {
    if (!isAuthenticated) {
      toast({
        title: "Authentication required",
        description: "Please log in to like posts",
        variant: "destructive",
      });
      return;
    }
    
    // Optimistic update
    setLiked(!liked);
    setLikeCount(prev => liked ? prev - 1 : prev + 1);
    
    if (liked) {
      unlikePostMutation.mutate();
    } else {
      likePostMutation.mutate();
    }
  };
  
  const handleShare = () => {
    if (!isAuthenticated) {
      toast({
        title: "Authentication required",
        description: "Please log in to share posts",
        variant: "destructive",
      });
      return;
    }
    
    setShareDialogOpen(true);
  };
  
  const submitShare = () => {
    sharePostMutation.mutate(shareContent);
  };
  
  const toggleComments = () => {
    setShowComments(!showComments);
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
  
  const isAuthor = user?.id === author.id;
  
  return (
    <div className="bg-white rounded-xl overflow-hidden mb-4 card-shadow animate-fade-in">
      <div className="p-4">
        <div className="flex items-start justify-between">
          <Link to={`/profile/${author.username}`} className="flex items-center space-x-3 group">
            <Avatar className="h-10 w-10">
              <AvatarImage src={author.avatar || "/placeholder.svg"} alt={author.name} />
              <AvatarFallback>{author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
            </Avatar>
            <div>
              <h3 className="font-medium group-hover:text-social-blue transition-colors">
                {author.name}
              </h3>
              <p className="text-sm text-gray-500">
                @{author.username} · {formatDate(createdAt)}
              </p>
            </div>
          </Link>
          
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm" className="h-8 w-8">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {isAuthor && (
                <>
                  <DropdownMenuItem>Edit Post</DropdownMenuItem>
                  <DropdownMenuItem className="text-red-500">Delete Post</DropdownMenuItem>
                </>
              )}
              <DropdownMenuItem>Save Post</DropdownMenuItem>
              <DropdownMenuItem>Hide Post</DropdownMenuItem>
              <DropdownMenuItem>Report</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        
        <div className="mt-3">
          <p className="text-gray-800 whitespace-pre-line">{content}</p>
        </div>
        
        {/* Shared post */}
        {sharedPost && (
          <div className="mt-3 border border-gray-200 rounded-lg p-3">
            <div className="flex items-center space-x-2">
              <Avatar className="h-6 w-6">
                <AvatarImage src={sharedPost.author.avatar || "/placeholder.svg"} alt={sharedPost.author.name} />
                <AvatarFallback>{sharedPost.author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>
              <div>
                <span className="font-medium text-sm">{sharedPost.author.name}</span>
                <span className="text-xs text-gray-500 ml-1">@{sharedPost.author.username} · {formatDate(sharedPost.createdAt)}</span>
              </div>
            </div>
            <p className="text-gray-800 text-sm mt-2">{sharedPost.content}</p>
            {sharedPost.image && (
              <div className="mt-2 rounded-lg overflow-hidden">
                <img 
                  src={sharedPost.image} 
                  alt="Shared post content" 
                  className="w-full h-auto object-cover max-h-64"
                />
              </div>
            )}
          </div>
        )}
        
        {image && (
          <div className="mt-3 rounded-lg overflow-hidden">
            <img 
              src={image} 
              alt="Post content" 
              className="w-full h-auto object-cover max-h-96"
            />
          </div>
        )}
        
        <div className="flex items-center justify-between mt-4 pt-3 border-t border-gray-100">
          <Button 
            variant="ghost" 
            size="sm" 
            className={`flex items-center space-x-1 ${liked ? 'text-red-500' : 'text-gray-500'}`}
            onClick={handleLike}
          >
            <Heart className="h-4 w-4" fill={liked ? "currentColor" : "none"} />
            <span>{likeCount}</span>
          </Button>
          
          <Button 
            variant="ghost" 
            size="sm" 
            className="flex items-center space-x-1 text-gray-500"
            onClick={toggleComments}
          >
            <MessageSquare className="h-4 w-4" />
            <span>{comments}</span>
          </Button>
          
          <Button 
            variant="ghost" 
            size="sm" 
            className="flex items-center space-x-1 text-gray-500"
            onClick={handleShare}
          >
            <Share className="h-4 w-4" />
            <span>{shares}</span>
          </Button>
        </div>
      </div>
      
      {showComments && <CommentSection postId={id} />}
      
      <Dialog open={shareDialogOpen} onOpenChange={setShareDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Share this post</DialogTitle>
            <DialogDescription>
              Add a comment to share this post with your followers
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <Textarea
              placeholder="Write something about this post..."
              value={shareContent}
              onChange={(e) => setShareContent(e.target.value)}
              className="min-h-24"
            />
            
            <div className="bg-gray-50 p-3 rounded-lg">
              <div className="flex items-center space-x-2">
                <Avatar className="h-6 w-6">
                  <AvatarImage src={author.avatar || "/placeholder.svg"} alt={author.name} />
                  <AvatarFallback>{author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium">{author.name}</span>
              </div>
              <p className="text-sm mt-2 line-clamp-2">{content}</p>
            </div>
          </div>
          <DialogFooter>
            <Button
              onClick={submitShare}
              disabled={sharePostMutation.isPending}
              className="w-full sm:w-auto"
            >
              {sharePostMutation.isPending ? "Sharing..." : "Share Post"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default Post;
