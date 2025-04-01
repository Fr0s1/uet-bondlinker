
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
import CommentSection from './CommentSection';

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
  isLiked = false 
}: PostProps) => {
  const [liked, setLiked] = useState(isLiked);
  const [likeCount, setLikeCount] = useState(likes);
  const [showComments, setShowComments] = useState(false);
  
  const handleLike = () => {
    if (liked) {
      setLikeCount(prev => prev - 1);
    } else {
      setLikeCount(prev => prev + 1);
    }
    setLiked(!liked);
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
  
  return (
    <div className="bg-white rounded-xl overflow-hidden mb-4 card-shadow animate-fade-in">
      <div className="p-4">
        <div className="flex items-start justify-between">
          <Link to={`/profile/${author.id}`} className="flex items-center space-x-3 group">
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
              <DropdownMenuItem>Save Post</DropdownMenuItem>
              <DropdownMenuItem>Hide Post</DropdownMenuItem>
              <DropdownMenuItem>Report</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        
        <div className="mt-3">
          <p className="text-gray-800 whitespace-pre-line">{content}</p>
        </div>
        
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
          >
            <Share className="h-4 w-4" />
            <span>{shares}</span>
          </Button>
        </div>
      </div>
      
      {showComments && <CommentSection postId={id} />}
    </div>
  );
};

export default Post;
