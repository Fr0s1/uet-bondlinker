
import React, { useState } from 'react';
import Post from './Post';
import { usePosts, useFeed, Post as PostType } from '@/hooks/use-posts';
import { useAuth } from '@/contexts/AuthContext';
import { Loader2 } from 'lucide-react';
import { Button } from './ui/button';

interface FeedProps {
  type?: 'public' | 'personal';
  userId?: string;
}

const Feed = ({ type = 'public', userId }: FeedProps) => {
  const { isAuthenticated } = useAuth();
  const { posts: publicPosts, isLoading: isPublicLoading, page: publicPage, setPage: setPublicPage } = usePosts(userId);
  const { posts: personalPosts, isLoading: isPersonalLoading, page: feedPage, setPage: setFeedPage } = useFeed();
  
  const posts = type === 'personal' && isAuthenticated ? personalPosts : publicPosts;
  const isLoading = type === 'personal' && isAuthenticated ? isPersonalLoading : isPublicLoading;
  const page = type === 'personal' && isAuthenticated ? feedPage : publicPage;
  const setPage = type === 'personal' && isAuthenticated ? setFeedPage : setPublicPage;
  
  if (isLoading && page === 1) {
    return (
      <div className="flex justify-center items-center py-8">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <span className="ml-2 text-muted-foreground">Loading posts...</span>
      </div>
    );
  }
  
  if (posts.length === 0 && !isLoading) {
    return (
      <div className="bg-card rounded-xl p-6 text-center shadow-sm my-4">
        <h3 className="text-lg font-medium">No posts yet</h3>
        <p className="text-muted-foreground mt-2">
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
      
      {posts.length > 0 && (
        <div className="flex justify-center my-4">
          <Button
            variant="outline"
            className="mx-auto"
            onClick={() => setPage(prevPage => prevPage + 1)}
            disabled={isLoading || posts.length < 10}
          >
            {isLoading ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Loading...
              </>
            ) : (
              'Load more'
            )}
          </Button>
        </div>
      )}
    </div>
  );
};

export default Feed;
