
import React, { useState, useRef } from 'react';
import { Link } from 'react-router';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Heart, MessageSquare, Share2, MoreHorizontal, Trash2 } from 'lucide-react';
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
  DialogClose,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Textarea } from "@/components/ui/textarea";
import CommentSection from './CommentSection';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from '@/components/ui/use-toast';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import EmojiPicker from 'emoji-picker-react';
import { type Post as PostType } from '@/hooks/use-posts';

export interface PostProps {
  post: PostType
}

const Post = ({
  post
}: PostProps) => {
  const [liked, setLiked] = useState(post.isLiked);
  const [likeCount, setLikeCount] = useState(post.likes);
  const [shareCount, setShareCount] = useState(post.shares);
  const [showComments, setShowComments] = useState(false);
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const [shareContent, setShareContent] = useState('');
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const { isAuthenticated, user } = useAuth();
  const queryClient = useQueryClient();
  const deleteDialogRef = useRef<HTMLDivElement>(null);

  const likePostMutation = useMutation({
    mutationFn: () => api.post(`/posts/${post.id}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      queryClient.invalidateQueries({ queryKey: ['trending'] });
    },
    onError: () => {
      setLiked(!liked);
      setLikeCount(prev => liked ? prev + 1 : prev - 1);
      toast({
        title: "Error",
        description: "Failed to like post. Please try again.",
        variant: "destructive"
      });
    }
  });

  const unlikePostMutation = useMutation({
    mutationFn: () => api.delete(`/posts/${post.id}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });
      queryClient.invalidateQueries({ queryKey: ['trending'] });
    },
    onError: () => {
      setLiked(!liked);
      setLikeCount(prev => liked ? prev - 1 : prev + 1);
      toast({
        title: "Error",
        description: "Failed to unlike post. Please try again.",
        variant: "destructive"
      });
    }
  });

  const sharePostMutation = useMutation({
    mutationFn: (content: string) => api.post(`/posts/${post.id}/share`, { content }),
    onSuccess: () => {
      setShareDialogOpen(false);
      setShareContent('');
      setShareCount(prev => prev + 1);
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

  const deletePostMutation = useMutation({
    mutationFn: () => api.delete(`/posts/${post.id}`),
    onSuccess: () => {
      // Close the dialog first, then invalidate queries
      setDeleteDialogOpen(false);

      // Use setTimeout to ensure the dialog is fully closed before invalidating queries
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: ['posts'] });
        queryClient.invalidateQueries({ queryKey: ['feed'] });
        queryClient.invalidateQueries({ queryKey: ['trending'] });
        queryClient.invalidateQueries({ queryKey: ['search'] });

        toast({
          title: "Post deleted",
          description: "Your post has been deleted successfully!",
        });
      }, 100);
    },
    onError: () => {
      setDeleteDialogOpen(false);
      toast({
        title: "Error",
        description: "Failed to delete the post. Please try again.",
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

  const handleDeletePost = () => {
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    deletePostMutation.mutate();
  };

  const toggleComments = () => {
    setShowComments(!showComments);
  };

  const onEmojiClick = (emojiObject: any) => {
    setShareContent(prev => prev + emojiObject.emoji);
    setShowEmojiPicker(false);
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

  const isAuthor = user?.id === post.author.id;

  return (
    <div className="bg-white rounded-xl overflow-hidden mb-4 card-shadow animate-fade-in">
      <div className="p-4">
        <div className="flex items-start justify-between">
          <Link to={`/profile/${post.author.username}`} className="flex items-center space-x-3 group">
            <Avatar className="h-10 w-10">
              <AvatarImage src={post.author.avatar} alt={post.author.name} />
              <AvatarFallback>{post.author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
            </Avatar>
            <div>
              <h3 className="font-medium group-hover:text-social-blue transition-colors text-left">
                {post.author.name}
              </h3>
              <p className="text-sm text-gray-500">
                @{post.author.username} · {formatDate(post.createdAt)}
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
                <DropdownMenuItem className="text-red-500" onClick={handleDeletePost}>
                  <Trash2 className="h-4 w-4 mr-2" /> Delete Post
                </DropdownMenuItem>
              )}
              <DropdownMenuItem>Save Post</DropdownMenuItem>
              <DropdownMenuItem>Hide Post</DropdownMenuItem>
              <DropdownMenuItem>Report</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="mt-3 text-left">
          <div className="text-gray-800 whitespace-pre-line prose prose-sm max-w-none">
            <div className="prose prose-sm max-w-none">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {post.content}
              </ReactMarkdown>
            </div>
          </div>
        </div>

        {post.sharedPost && (
          <div className="mt-3 border border-gray-200 rounded-lg p-3">
            <div className="flex items-center space-x-2">
              <Avatar className="h-6 w-6">
                <AvatarImage src={post.sharedPost.author.avatar} alt={post.sharedPost.author.name} />
                <AvatarFallback>{post.sharedPost.author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>
              <div>
                <span className="font-medium text-sm">{post.sharedPost.author.name}</span>
                <span className="text-xs text-gray-500 ml-1">@{post.sharedPost.author.username} · {formatDate(post.sharedPost.createdAt)}</span>
              </div>
            </div>
            <div className="mt-2 text-left">
              <div className="text-gray-800 text-sm">
                <div className="prose prose-sm max-w-none">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {post.sharedPost.content}
                  </ReactMarkdown>
                </div>
              </div>
            </div>
            {post.sharedPost.image && (
              <div className="mt-2 rounded-lg overflow-hidden">
                <img
                  src={post.sharedPost.image}
                  alt="Shared post content"
                  className="w-full h-auto object-cover max-h-64"
                />
              </div>
            )}
          </div>
        )}

        {post.image && (
          <div className="mt-3 rounded-lg overflow-hidden">
            <img
              src={post.image}
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
            <span>{post.comments}</span>
          </Button>

          <Button
            variant="ghost"
            size="sm"
            className="flex items-center space-x-1 text-gray-500"
            onClick={handleShare}
          >
            <Share2 className="h-4 w-4" />
            <span>{shareCount}</span>
          </Button>
        </div>
      </div>

      {showComments && <CommentSection postId={post.id} />}

      <Dialog open={shareDialogOpen} onOpenChange={setShareDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Share this post</DialogTitle>
            <DialogDescription>
              Add a comment to share this post with your followers
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="relative">
              <Textarea
                placeholder="Write something about this post..."
                value={shareContent}
                onChange={(e) => setShareContent(e.target.value)}
                className="min-h-24"
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute bottom-2 right-2"
                onClick={() => setShowEmojiPicker(!showEmojiPicker)}
              >
                😊
              </Button>
              {showEmojiPicker && (
                <div className="absolute bottom-12 right-0 z-10">
                  <EmojiPicker onEmojiClick={onEmojiClick} />
                </div>
              )}
            </div>

            <div className="bg-gray-50 p-3 rounded-lg">
              <div className="flex items-center space-x-2">
                <Avatar className="h-6 w-6">
                  <AvatarImage src={post.author.avatar} alt={post.author.name} />
                  <AvatarFallback>{post.author.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium">{post.author.name}</span>
              </div>
              <div className="text-sm mt-2 line-clamp-2 text-left">
                <div className="prose prose-sm max-w-none">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {post.content}
                  </ReactMarkdown>
                </div>
              </div>
            </div>
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
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

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent ref={deleteDialogRef}>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you want to delete this post?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete your post.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-red-500 hover:bg-red-600"
              onClick={(e) => {
                e.preventDefault();
                confirmDelete();
                return false;
              }}
              disabled={deletePostMutation.isPending}
            >
              {deletePostMutation.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

export default Post;
