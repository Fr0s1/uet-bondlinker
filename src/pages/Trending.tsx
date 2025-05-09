
import React from 'react';
import Post from '@/components/Post';
import { useTrendingPosts } from '@/hooks/use-posts';
import { Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';

const Trending = () => {
  const { posts, isLoading, page, setPage } = useTrendingPosts();

  if (isLoading && page === 1) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
        <span className="ml-2 text-gray-500">Loading trending posts...</span>
      </div>
    );
  }

  if (posts.length === 0 && !isLoading) {
    return (
      <div className="bg-white rounded-xl p-8 text-center card-shadow my-4">
        <h3 className="text-lg font-medium text-gray-700">No trending posts yet</h3>
        <p className="text-gray-500 mt-2">
          Check back later for the most popular content!
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      <h1 className="text-2xl font-bold mb-6">Trending Posts</h1>

      <div className="space-y-4">
        {posts.map((post) => (
          <Post
            key={post.id}
            post={post}
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
    </div>
  );
};

export default Trending;
