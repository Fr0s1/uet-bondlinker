
import React from 'react';
import { Link } from 'react-router-dom';
import { Bell, Home, MessageSquare, Search, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/contexts/AuthContext';

const Navbar = () => {
  const { user } = useAuth();
  
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
          />
        </div>

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
