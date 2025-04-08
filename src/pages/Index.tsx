
import React from 'react';
import Navbar from '@/components/Navbar';
import PostForm from '@/components/PostForm';
import Feed from '@/components/Feed';
import UserProfile from '@/components/UserProfile';
import { useAuth } from '@/contexts/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';

const Index = () => {
  const { user, isAuthenticated } = useAuth();
  const [refreshFeed, setRefreshFeed] = React.useState(false);
  
  // Fetch suggested users to follow
  const { data: suggestedUsers, isLoading: isSuggestedUsersLoading } = useQuery({
    queryKey: ['suggested-users'],
    queryFn: () => api.get('/users/suggested'),
    enabled: isAuthenticated
  });
  
  // Fetch trending topics
  const { data: trendingTopics, isLoading: isTrendingTopicsLoading } = useQuery({
    queryKey: ['trending-topics'],
    queryFn: () => api.get('/posts/trending'),
  });
  
  const handlePostCreated = () => {
    setRefreshFeed(!refreshFeed);
  };
  
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="container py-4">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
          {/* Left Sidebar - User Profile */}
          <aside className="hidden lg:block lg:col-span-3">
            {isAuthenticated && user && (
              <div className="sticky top-20">
                <UserProfile 
                  user={{
                    id: user.id,
                    name: user.name,
                    username: user.username,
                    avatar: user.avatar || "/placeholder.svg",
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
          
          {/* Main Content */}
          <main className="lg:col-span-6">
            <PostForm onPostCreated={handlePostCreated} />
            <Feed />
          </main>
          
          {/* Right Sidebar - Suggestions and Trends */}
          <aside className="hidden lg:block lg:col-span-3">
            <div className="sticky top-20 space-y-4">
              {/* Who to follow */}
              {isAuthenticated && (
                <div className="bg-white rounded-xl p-4 card-shadow animate-fade-in">
                  <h3 className="font-semibold text-lg mb-4">Who to follow</h3>
                  {isSuggestedUsersLoading ? (
                    <div className="flex justify-center py-6">
                      <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
                    </div>
                  ) : suggestedUsers && suggestedUsers.length > 0 ? (
                    <div className="space-y-4">
                      {suggestedUsers.map((user: any) => (
                        <div key={user.id} className="flex items-center justify-between">
                          <div className="flex items-center space-x-2">
                            <img 
                              src={user.avatar || "/placeholder.svg"} 
                              alt={user.name} 
                              className="w-10 h-10 rounded-full avatar-shadow"
                            />
                            <div>
                              <p className="font-medium text-sm">{user.name}</p>
                              <p className="text-xs text-gray-500">@{user.username}</p>
                            </div>
                          </div>
                          <Button 
                            size="sm" 
                            className="h-8 gradient-blue"
                            onClick={() => {
                              api.post(`/users/follow/${user.id}`);
                            }}
                          >
                            Follow
                          </Button>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-center py-2 text-gray-500">No suggestions available</p>
                  )}
                </div>
              )}
              
              {/* Trending topics */}
              <div className="bg-white rounded-xl p-4 card-shadow animate-fade-in">
                <h3 className="font-semibold text-lg mb-4">Trends for you</h3>
                {isTrendingTopicsLoading ? (
                  <div className="flex justify-center py-6">
                    <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
                  </div>
                ) : trendingTopics && trendingTopics.length > 0 ? (
                  <div className="space-y-4">
                    {trendingTopics.map((topic: any) => (
                      <div key={topic.id} className="group cursor-pointer">
                        <h4 className="font-medium group-hover:text-social-blue transition-colors">
                          {topic.name}
                        </h4>
                        <p className="text-xs text-gray-500">{topic.posts} posts</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-center py-2 text-gray-500">No trending topics available</p>
                )}
              </div>
              
              {/* Footer Links */}
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
      </div>
    </div>
  );
};

export default Index;
