
import React from 'react';
import Navbar from '@/components/Navbar';
import PostForm from '@/components/PostForm';
import Feed from '@/components/Feed';
import UserProfile from '@/components/UserProfile';
import { Button } from '@/components/ui/button';
import { Bell, Users, Link as LinkIcon } from 'lucide-react';

// Mock user data
const currentUser = {
  id: "u1",
  name: "John Doe",
  username: "johndoe",
  avatar: "/placeholder.svg",
  bio: "Web developer and designer. Passionate about creating beautiful user experiences.",
  followers: 245,
  following: 186,
  location: "San Francisco, CA",
  website: "https://johndoe.com",
  joinedDate: "2022-01-15T00:00:00Z",
  isFollowing: false,
};

// Mock suggestions data
const suggestedUsers = [
  {
    id: "u5",
    name: "Sarah Williams",
    username: "sarahw",
    avatar: "/placeholder.svg",
  },
  {
    id: "u6",
    name: "Michael Brown",
    username: "michaelb",
    avatar: "/placeholder.svg",
  },
  {
    id: "u7",
    name: "Jennifer Lee",
    username: "jenniferl",
    avatar: "/placeholder.svg",
  },
];

// Mock trends data
const trendingTopics = [
  { id: "t1", name: "#WebDevelopment", posts: "5.2K" },
  { id: "t2", name: "#ArtificialIntelligence", posts: "12K" },
  { id: "t3", name: "#Photography", posts: "3.8K" },
  { id: "t4", name: "#TravelTips", posts: "7.1K" },
];

const Index = () => {
  const [refreshFeed, setRefreshFeed] = React.useState(false);
  
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
            <div className="sticky top-20">
              <UserProfile user={currentUser} isCurrentUser={true} />
            </div>
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
              <div className="bg-white rounded-xl p-4 card-shadow animate-fade-in">
                <h3 className="font-semibold text-lg mb-4">Who to follow</h3>
                <div className="space-y-4">
                  {suggestedUsers.map((user) => (
                    <div key={user.id} className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <img 
                          src={user.avatar} 
                          alt={user.name} 
                          className="w-10 h-10 rounded-full avatar-shadow"
                        />
                        <div>
                          <p className="font-medium text-sm">{user.name}</p>
                          <p className="text-xs text-gray-500">@{user.username}</p>
                        </div>
                      </div>
                      <Button size="sm" className="h-8 gradient-blue">Follow</Button>
                    </div>
                  ))}
                </div>
                <Button variant="ghost" className="text-social-blue w-full mt-3">
                  Show more
                </Button>
              </div>
              
              {/* Trending topics */}
              <div className="bg-white rounded-xl p-4 card-shadow animate-fade-in">
                <h3 className="font-semibold text-lg mb-4">Trends for you</h3>
                <div className="space-y-4">
                  {trendingTopics.map((topic) => (
                    <div key={topic.id} className="group cursor-pointer">
                      <h4 className="font-medium group-hover:text-social-blue transition-colors">
                        {topic.name}
                      </h4>
                      <p className="text-xs text-gray-500">{topic.posts} posts</p>
                    </div>
                  ))}
                </div>
                <Button variant="ghost" className="text-social-blue w-full mt-3">
                  Show more
                </Button>
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
