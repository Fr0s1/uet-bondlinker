
import React from 'react';
import Post from './Post';
import { usePosts, useFeed, Post as PostType } from '@/hooks/use-posts';
import { useAuth } from '@/contexts/AuthContext';
import { Loader2 } from 'lucide-react';

interface FeedProps {
  type?: 'public' | 'personal';
}

const Feed = ({ type = 'public' }: FeedProps) => {
  const { isAuthenticated } = useAuth();
  const { posts: publicPosts, isLoading: isPublicLoading } = usePosts();
  const { posts: personalPosts, isLoading: isPersonalLoading } = useFeed();
  
  const posts = type === 'personal' && isAuthenticated ? personalPosts : publicPosts;
  const isLoading = type === 'personal' && isAuthenticated ? isPersonalLoading : isPublicLoading;
  
  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
        <span className="ml-2 text-gray-500">Loading posts...</span>
      </div>
    );
  }
  
  if (posts.length === 0) {
    return (
      <div className="bg-white rounded-xl p-8 text-center card-shadow my-4">
        <h3 className="text-lg font-medium text-gray-700">No posts yet</h3>
        <p className="text-gray-500 mt-2">
          {type === 'personal' 
            ? "Follow users to see their posts in your feed!"
            : "Be the first to create a post!"}
        </p>
      </div>
    );
  }
  
  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <Post 
          key={post.id} 
          id={post.id}
          author={{
            id: post.user_id,
            name: post.author?.name || "Unknown User",
            username: post.author?.username || "unknown",
            avatar: post.author?.avatar || "/placeholder.svg",
          }}
          content={post.content}
          image={post.image}
          createdAt={post.created_at}
          likes={post.likes}
          comments={post.comments}
          shares={0} // API doesn't support shares yet
          isLiked={post.is_liked}
        />
      ))}
    </div>
  );
};

export default Feed;
