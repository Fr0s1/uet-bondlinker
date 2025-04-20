
import React from 'react';
import PostForm from '@/components/PostForm';
import Feed from '@/components/Feed';
import UserProfile from '@/components/UserProfile';
import { useAuth } from '@/contexts/AuthContext';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useSuggestedPosts } from '@/hooks/use-posts';
import { useFeed } from '@/hooks/use-posts';
import Post from '@/components/Post';
import { Loader2 } from 'lucide-react';

const Index = () => {
  const { user } = useAuth();
  const [refreshFeed, setRefreshFeed] = React.useState(false);
  const [activeTab, setActiveTab] = React.useState('following');

  const { posts: followingPosts, isLoading: isFollowingLoading } = useFeed();
  const { posts: suggestedPosts, isLoading: isSuggestedLoading } = useSuggestedPosts();

  const handlePostCreated = () => {
    setRefreshFeed(!refreshFeed);
  };

  const renderFollowingFeed = () => {
    if (isFollowingLoading) {
      return (
        <div className="flex justify-center items-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
          <span className="ml-2 text-gray-500">Loading feed...</span>
        </div>
      );
    }

    if (followingPosts.length === 0) {
      return (
        <div className="bg-white rounded-xl p-6 card-shadow text-center">
          <h3 className="text-lg font-medium mb-2">Your feed is empty</h3>
          <p className="text-gray-500 mb-4">Follow other users to see their posts in your feed</p>
        </div>
      );
    }

    return (
      <div className="space-y-4">
        {followingPosts.map((post) => (
          <Post
            key={post.id}
            post={post}
          />
        ))}
      </div>
    );
  };

  const renderSuggestedFeed = () => {
    if (isSuggestedLoading) {
      return (
        <div className="flex justify-center items-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
          <span className="ml-2 text-gray-500">Loading suggestions...</span>
        </div>
      );
    }

    if (suggestedPosts.length === 0) {
      return (
        <div className="bg-white rounded-xl p-6 card-shadow text-center">
          <h3 className="text-lg font-medium mb-2">No suggestions available</h3>
          <p className="text-gray-500">We'll have more personalized content for you soon</p>
        </div>
      );
    }

    return (
      <div className="space-y-4">
        {suggestedPosts.map((post) => (
          <Post
            key={post.id}
            post={post}
          />
        ))}
      </div>
    );
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
      <aside className="hidden lg:block lg:col-span-3">
        {user && (
          <div className="sticky top-20">
            <UserProfile
              user={{
                id: user.id,
                name: user.name,
                username: user.username,
                avatar: user.avatar || "/user-avatar.png",
                bio: user.bio || "No bio provided",
                location: user.location,
                website: user.website,
                joinedDate: user.createdAt,
                followers: user.followers || 0,
                following: user.following || 0,
                isFollowing: false,
              }}
              isCurrentUser={true}
            />
          </div>
        )}
      </aside>

      <main className="lg:col-span-6 pt-4">
        <PostForm onPostCreated={handlePostCreated} />

        <Tabs value={activeTab} onValueChange={setActiveTab} className="mb-6">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="following">Following</TabsTrigger>
            <TabsTrigger value="suggested">Suggested</TabsTrigger>
          </TabsList>
          <TabsContent value="following" className="mt-6">
            {renderFollowingFeed()}
          </TabsContent>
          <TabsContent value="suggested" className="mt-6">
            {renderSuggestedFeed()}
          </TabsContent>
        </Tabs>
      </main>

      <aside className="hidden lg:block lg:col-span-3">
        <div className="sticky top-20 space-y-4">
          <div className="bg-white rounded-xl p-4 card-shadow animate-fade-in">
            <h3 className="font-semibold text-lg mb-4">Trending Topics</h3>
            <div className="space-y-4">
              <div className="group cursor-pointer">
                <h4 className="font-medium group-hover:text-social-blue transition-colors">
                  #WebDevelopment
                </h4>
                <p className="text-xs text-gray-500">1,234 posts</p>
              </div>
              <div className="group cursor-pointer">
                <h4 className="font-medium group-hover:text-social-blue transition-colors">
                  #ArtificialIntelligence
                </h4>
                <p className="text-xs text-gray-500">985 posts</p>
              </div>
              <div className="group cursor-pointer">
                <h4 className="font-medium group-hover:text-social-blue transition-colors">
                  #ProgrammingHumor
                </h4>
                <p className="text-xs text-gray-500">743 posts</p>
              </div>
            </div>
          </div>

          <div className="p-4 text-xs text-gray-500">
            <div className="flex flex-wrap gap-2">
              <a href="#" className="hover:underline">Terms of Service</a>
              <a href="#" className="hover:underline">Privacy Policy</a>
              <a href="#" className="hover:underline">Cookie Policy</a>
              <a href="#" className="hover:underline">Accessibility</a>
              <a href="#" className="hover:underline">Ads Info</a>
            </div>
            <p className="mt-2">© 2023 SocialNet</p>
          </div>
        </div>
      </aside>
    </div>
  );
};

export default Index;
