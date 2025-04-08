
import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Search as SearchIcon, Loader2 } from 'lucide-react';
import Post from '@/components/Post';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Link } from 'react-router-dom';
import { useSearch, useSearchUsers, useSearchPosts } from '@/hooks/use-search';

const Search = () => {
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const initialQuery = searchParams.get('q') || '';
  
  const [query, setQuery] = useState(initialQuery);
  const [searchTerm, setSearchTerm] = useState(initialQuery);
  const [activeTab, setActiveTab] = useState('all');
  
  const { results, isLoading } = useSearch(searchTerm);
  const { users, isLoading: isUsersLoading } = useSearchUsers(searchTerm);
  const { posts, isLoading: isPostsLoading } = useSearchPosts(searchTerm);
  
  useEffect(() => {
    // Update the URL when search term changes
    if (searchTerm) {
      const newUrl = `/search?q=${encodeURIComponent(searchTerm)}`;
      window.history.pushState({}, '', newUrl);
    }
  }, [searchTerm]);
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchTerm(query);
  };
  
  const renderResults = () => {
    if (activeTab === 'all') {
      if (isLoading) {
        return (
          <div className="flex justify-center items-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
            <span className="ml-2 text-gray-500">Searching...</span>
          </div>
        );
      }
      
      if (results.users.length === 0 && results.posts.length === 0) {
        return (
          <div className="text-center py-12">
            <p className="text-gray-500">No results found for "{searchTerm}"</p>
          </div>
        );
      }
      
      return (
        <>
          {results.users.length > 0 && (
            <div className="mb-8">
              <h3 className="text-lg font-bold mb-4">People</h3>
              <div className="space-y-4">
                {results.users.slice(0, 3).map((user) => (
                  <Link 
                    key={user.id} 
                    to={`/profile/${user.username}`}
                    className="flex items-center p-4 bg-white rounded-xl card-shadow hover:shadow-md transition-shadow"
                  >
                    <Avatar className="h-12 w-12">
                      <AvatarImage src={user.avatar || "/placeholder.svg"} alt={user.name} />
                      <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                    </Avatar>
                    <div className="ml-4">
                      <h4 className="font-medium">{user.name}</h4>
                      <p className="text-sm text-gray-500">@{user.username}</p>
                    </div>
                  </Link>
                ))}
                {results.users.length > 3 && (
                  <Button 
                    variant="link" 
                    className="w-full text-social-blue"
                    onClick={() => setActiveTab('people')}
                  >
                    View all {results.users.length} people
                  </Button>
                )}
              </div>
            </div>
          )}
          
          {results.posts.length > 0 && (
            <div>
              <h3 className="text-lg font-bold mb-4">Posts</h3>
              <div className="space-y-4">
                {results.posts.slice(0, 3).map((post) => (
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
                    shares={post.shares || 0}
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
                {results.posts.length > 3 && (
                  <Button 
                    variant="link" 
                    className="w-full text-social-blue"
                    onClick={() => setActiveTab('posts')}
                  >
                    View all {results.posts.length} posts
                  </Button>
                )}
              </div>
            </div>
          )}
        </>
      );
    }
    
    if (activeTab === 'people') {
      if (isUsersLoading) {
        return (
          <div className="flex justify-center items-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
            <span className="ml-2 text-gray-500">Loading users...</span>
          </div>
        );
      }
      
      if (users.length === 0) {
        return (
          <div className="text-center py-12">
            <p className="text-gray-500">No users found for "{searchTerm}"</p>
          </div>
        );
      }
      
      return (
        <div className="space-y-4">
          {users.map((user) => (
            <Link 
              key={user.id} 
              to={`/profile/${user.username}`}
              className="flex items-center p-4 bg-white rounded-xl card-shadow hover:shadow-md transition-shadow"
            >
              <Avatar className="h-12 w-12">
                <AvatarImage src={user.avatar || "/placeholder.svg"} alt={user.name} />
                <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>
              <div className="ml-4">
                <h4 className="font-medium">{user.name}</h4>
                <p className="text-sm text-gray-500">@{user.username}</p>
                {user.bio && <p className="text-sm text-gray-700 mt-1">{user.bio}</p>}
              </div>
            </Link>
          ))}
        </div>
      );
    }
    
    if (activeTab === 'posts') {
      if (isPostsLoading) {
        return (
          <div className="flex justify-center items-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
            <span className="ml-2 text-gray-500">Loading posts...</span>
          </div>
        );
      }
      
      if (posts.length === 0) {
        return (
          <div className="text-center py-12">
            <p className="text-gray-500">No posts found for "{searchTerm}"</p>
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
              shares={post.shares || 0}
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
        </div>
      );
    }
    
    return null;
  };
  
  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      <form onSubmit={handleSearch} className="mb-6">
        <div className="flex gap-2">
          <Input
            type="text"
            placeholder="Search for people or posts..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="flex-1"
          />
          <Button type="submit">
            <SearchIcon className="h-4 w-4 mr-2" />
            Search
          </Button>
        </div>
      </form>
      
      {searchTerm && (
        <Tabs value={activeTab} onValueChange={setActiveTab} className="mb-6">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="people">People</TabsTrigger>
            <TabsTrigger value="posts">Posts</TabsTrigger>
          </TabsList>
          <TabsContent value="all" className="mt-6">{renderResults()}</TabsContent>
          <TabsContent value="people" className="mt-6">{renderResults()}</TabsContent>
          <TabsContent value="posts" className="mt-6">{renderResults()}</TabsContent>
        </Tabs>
      )}
    </div>
  );
};

export default Search;
