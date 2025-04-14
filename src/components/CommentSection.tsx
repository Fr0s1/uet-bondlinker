
import React, { useState } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Heart, Send, Loader2 } from 'lucide-react';
import { useComments, Comment, CreateCommentData } from '@/hooks/use-comments';
import { useAuth } from '@/contexts/AuthContext';
import { formatDistanceToNow } from 'date-fns';

const CommentSection = ({ postId }: { postId: string }) => {
  const { user } = useAuth();
  const [newComment, setNewComment] = useState('');
  const {
    comments,
    isLoading,
    createComment,
    deleteComment,
    updateComment
  } = useComments(postId);

  const handleSubmitComment = (e: React.FormEvent) => {
    e.preventDefault();

    if (!newComment.trim()) return;

    const commentData: CreateCommentData = {
      content: newComment
    };

    createComment.mutate(commentData, {
      onSuccess: () => {
        setNewComment('');
      }
    });
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return formatDistanceToNow(date, { addSuffix: true });
  };

  return (
    <div className="bg-gray-50 border-t p-4">
      <form onSubmit={handleSubmitComment} className="flex items-start space-x-3 mb-4">
        <Avatar className="h-8 w-8">
          <AvatarImage src={user?.avatar || "/placeholder.svg"} alt={user?.name || "User"} />
          <AvatarFallback>{user?.name?.slice(0, 2).toUpperCase() || "U"}</AvatarFallback>
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
            disabled={createComment.isPending || !newComment.trim()}
          >
            {createComment.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
          </Button>
        </div>
      </form>

      {isLoading ? (
        <div className="flex justify-center py-4">
          <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
          <span className="ml-2 text-gray-500">Loading comments...</span>
        </div>
      ) : comments.length > 0 ? (
        <div className="space-y-4">
          {comments.map((comment) => (
            <div key={comment.id} className="flex space-x-3">
              <Avatar className="h-8 w-8">
                <AvatarImage
                  src={comment.author?.avatar || "/placeholder.svg"}
                  alt={comment.author?.name || "User"}
                />
                <AvatarFallback>
                  {comment.author?.name?.slice(0, 2).toUpperCase() || "U"}
                </AvatarFallback>
              </Avatar>

              <div className="flex-1">
                <div className="bg-white p-3 rounded-lg">
                  <div className="flex justify-between items-start">
                    <h4 className="font-medium text-sm">{comment.author?.name || "Unknown User"}</h4>
                    <span className="text-xs text-gray-500">{formatDate(comment.created_at)}</span>
                  </div>
                  <p className="text-sm mt-1 text-left">{comment.content}</p>
                </div>

                <div className="flex items-center space-x-4 mt-1 ml-1">
                  {user?.id === comment.user_id && (
                    <button
                      className="text-xs text-gray-500 hover:text-red-500"
                      onClick={() => deleteComment.mutate(comment.id)}
                    >
                      {deleteComment.isPending ? "Deleting..." : "Delete"}
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-4 text-gray-500">
          No comments yet. Be the first to comment!
        </div>
      )}
    </div>
  );
};

export default CommentSection;
