
import Post from './Post';
import { usePosts, useFeed } from '@/hooks/use-posts';
import { Loader2 } from 'lucide-react';
import { Button } from './ui/button';

interface FeedProps {
  type?: 'public' | 'personal';
  userId?: string;
}

const Feed = ({ type = 'public', userId }: FeedProps) => {

  let posts: any[];
  let isLoading: boolean;
  let page: number;
  let setPage: any;

  if (type == 'personal') {
    const { posts: personalPosts, isLoading: isPersonalLoading, page: feedPage, setPage: setFeedPage } = useFeed();
    posts = personalPosts
    isLoading = isPersonalLoading
    page = feedPage
    setPage = setFeedPage
  } else {
    const { posts: publicPosts, isLoading: isPublicLoading, page: publicPage, setPage: setPublicPage } = usePosts(userId);
    posts = publicPosts
    isLoading = isPublicLoading
    page = publicPage
    setPage = setPublicPage
  }

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
          shares={post.shares}
          isLiked={post.is_liked}
          sharedPost={post.shared_post ? {
            id: post.shared_post.id,
            author: {
              id: post.shared_post.user_id,
              name: post.shared_post.author?.name || "Unknown User",
              username: post.shared_post.author?.username || "unknown",
              avatar: post.shared_post.author?.avatar || "/placeholder.svg",
            },
            content: post.shared_post.content,
            image: post.shared_post.image,
            createdAt: post.shared_post.created_at,
          } : undefined}
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
