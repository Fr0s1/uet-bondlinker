
import React, { useState, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Bell, Home, MessageSquare, Search, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/contexts/AuthContext';
import { 
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from '@/components/ui/command';
import { useSearch } from '@/hooks/use-search';

const Navbar = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);
  const { results, isLoading } = useSearch(searchQuery);
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery)}`);
      setOpen(false);
    }
  };

  const handleOpenSearch = () => {
    setOpen(true);
  };
  
  return (
    <nav className="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm py-3">
      <div className="container flex items-center justify-between">
        <div className="flex items-center">
          <Link to="/" className="flex items-center space-x-2">
            <div className="w-9 h-9 rounded-full gradient-blue flex items-center justify-center">
              <span className="text-white font-bold text-lg">S</span>
            </div>
            <span className="font-bold text-xl hidden md:block">SocialNet</span>
          </Link>
        </div>

        <div className="hidden md:flex items-center w-1/3 relative">
          <Search className="absolute left-3 h-4 w-4 text-gray-400" />
          <Input 
            placeholder="Search..." 
            className="pl-10 h-9 bg-gray-50 rounded-xl"
            onClick={handleOpenSearch}
            ref={inputRef}
          />
        </div>

        <CommandDialog open={open} onOpenChange={setOpen}>
          <form onSubmit={handleSearch}>
            <CommandInput
              placeholder="Search for people or posts..."
              value={searchQuery}
              onValueChange={setSearchQuery}
              autoFocus
            />
          </form>
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            {results.users && results.users.length > 0 && (
              <CommandGroup heading="People">
                {results.users.slice(0, 4).map((user) => (
                  <CommandItem
                    key={user.id}
                    onSelect={() => {
                      navigate(`/profile/${user.username}`);
                      setOpen(false);
                    }}
                  >
                    <Avatar className="h-6 w-6 mr-2">
                      <AvatarImage src={user.avatar || "/placeholder.svg"} alt={user.name} />
                      <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                    </Avatar>
                    <span>{user.name}</span>
                    <span className="text-sm text-gray-500 ml-1">@{user.username}</span>
                  </CommandItem>
                ))}
                {results.users.length > 4 && (
                  <CommandItem
                    onSelect={() => {
                      navigate(`/search?q=${encodeURIComponent(searchQuery)}&tab=people`);
                      setOpen(false);
                    }}
                  >
                    <span className="text-social-blue">See all people</span>
                  </CommandItem>
                )}
              </CommandGroup>
            )}
            {results.posts && results.posts.length > 0 && (
              <CommandGroup heading="Posts">
                {results.posts.slice(0, 4).map((post) => (
                  <CommandItem
                    key={post.id}
                    onSelect={() => {
                      navigate(`/post/${post.id}`);
                      setOpen(false);
                    }}
                  >
                    <Avatar className="h-6 w-6 mr-2">
                      <AvatarImage src={post.author?.avatar || "/placeholder.svg"} alt={post.author?.name} />
                      <AvatarFallback>{(post.author?.name || "U").slice(0, 2).toUpperCase()}</AvatarFallback>
                    </Avatar>
                    <span className="truncate">{post.content.slice(0, 40)}{post.content.length > 40 ? '...' : ''}</span>
                  </CommandItem>
                ))}
                {results.posts.length > 4 && (
                  <CommandItem
                    onSelect={() => {
                      navigate(`/search?q=${encodeURIComponent(searchQuery)}&tab=posts`);
                      setOpen(false);
                    }}
                  >
                    <span className="text-social-blue">See all posts</span>
                  </CommandItem>
                )}
              </CommandGroup>
            )}
          </CommandList>
        </CommandDialog>

        <div className="flex items-center space-x-1 sm:space-x-3">
          <Button variant="ghost" size="icon" className="text-gray-600" asChild>
            <Link to="/"><Home className="h-5 w-5" /></Link>
          </Button>
          <Button variant="ghost" size="icon" className="text-gray-600" asChild>
            <Link to="/messages"><MessageSquare className="h-5 w-5" /></Link>
          </Button>
          <Button variant="ghost" size="icon" className="text-gray-600" asChild>
            <Link to="/notifications"><Bell className="h-5 w-5" /></Link>
          </Button>
          <Button variant="ghost" size="icon" className="text-gray-600 hidden sm:flex" asChild>
            <Link to={user ? `/profile/${user.username}` : "/login"}><User className="h-5 w-5" /></Link>
          </Button>
          <Link to={user ? `/profile/${user.username}` : "/login"} className="ml-2">
            <Avatar className="h-8 w-8 avatar-shadow">
              <AvatarImage src={user?.avatar || "/placeholder.svg"} alt="Profile" />
              <AvatarFallback>{user ? user.name.substring(0, 2).toUpperCase() : "US"}</AvatarFallback>
            </Avatar>
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
