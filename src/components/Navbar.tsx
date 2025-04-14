
import { Link, useLocation, useNavigate } from 'react-router';
import { Home, Search, TrendingUp, MessageCircle, UserCircle, LogOut, Menu, X, Lock } from 'lucide-react';
import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { useIsMobile } from '@/hooks/use-mobile';
import NotificationCenter from '@/components/NotificationCenter';

const Navbar = () => {
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const isMobile = useIsMobile();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const isActive = (path: string) => location.pathname === path;

  const navItems = [
    { icon: <Home size={24} />, text: 'Home', path: '/' },
    { icon: <Search size={24} />, text: 'Search', path: '/search' },
    { icon: <TrendingUp size={24} />, text: 'Trending', path: '/trending' },
    { icon: <MessageCircle size={24} />, text: 'Messages', path: '/messages' },
  ];

  if (user) {
    navItems.push({
      icon: <UserCircle size={24} />,
      text: 'Profile',
      path: `/profile/${user.username}`
    })
  }

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const renderNavItems = () => (
    <>
      {navItems.map((item) => (
        <Link
          key={item.path}
          to={item.path}
          className={`flex items-center space-x-2 py-2 px-3 rounded-md transition-colors hover:bg-gray-100 ${isActive(item.path) ? 'text-social-blue font-medium' : 'text-gray-700'
            }`}
          onClick={() => setIsMenuOpen(false)}
        >
          {item.icon}
          <span>{item.text}</span>
        </Link>
      ))}
      <Link
        to="/change-password"
        className={`flex items-center space-x-2 py-2 px-3 rounded-md transition-colors hover:bg-gray-100 ${isActive('/change-password') ? 'text-social-blue font-medium' : 'text-gray-700'
          }`}
        onClick={() => setIsMenuOpen(false)}
      >
        <Lock size={24} />
        <span>Change Password</span>
      </Link>
      <button
        onClick={handleLogout}
        className="flex items-center space-x-2 py-2 px-3 rounded-md transition-colors hover:bg-gray-100 text-gray-700 w-full text-left"
      >
        <LogOut size={24} />
        <span>Logout</span>
      </button>
    </>
  );

  return (
    <nav className="bg-white shadow-sm border-b sticky top-0 z-10">
      <div className="max-w-6xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link to="/" className="flex items-center">
              <h1 className="text-xl font-bold text-social-blue">SocialNet</h1>
            </Link>
          </div>

          {isMobile ? (
            <div className="flex items-center space-x-4">
              <NotificationCenter />

              <Sheet open={isMenuOpen} onOpenChange={setIsMenuOpen}>
                <SheetTrigger asChild>
                  <button className="p-1">
                    {isMenuOpen ? <X size={24} /> : <Menu size={24} />}
                  </button>
                </SheetTrigger>
                <SheetContent side="right">
                  <div className="py-4 flex flex-col space-y-3">
                    {user && (
                      <div className="flex items-center space-x-3 mb-4 pb-4 border-b">
                        <Avatar>
                          <AvatarImage src={user.avatar || undefined} alt={user.name} />
                          <AvatarFallback>{user.name?.charAt(0)}</AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium">{user.name}</p>
                          <p className="text-sm text-gray-500">@{user.username}</p>
                        </div>
                      </div>
                    )}
                    {renderNavItems()}
                  </div>
                </SheetContent>
              </Sheet>
            </div>
          ) : (
            <div className="flex items-center space-x-4">
              <div className="hidden md:flex space-x-4">
                {navItems.map((item) => (
                  <Link
                    key={item.path}
                    to={item.path}
                    className={`flex items-center py-2 px-3 rounded-md transition-colors hover:bg-gray-100 ${isActive(item.path) ? 'text-social-blue font-medium' : 'text-gray-700'
                      }`}
                  >
                    {item.icon}
                  </Link>
                ))}
              </div>

              <NotificationCenter />

              {user && (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button className="focus:outline-none">
                      <Avatar>
                        <AvatarImage src={user.avatar || undefined} alt={user.name} />
                        <AvatarFallback>{user.name?.charAt(0)}</AvatarFallback>
                      </Avatar>
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56">
                    <div className="flex items-center space-x-3 p-2 border-b">
                      <Avatar>
                        <AvatarImage src={user.avatar || undefined} alt={user.name} />
                        <AvatarFallback>{user.name?.charAt(0)}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium">{user.name}</p>
                        <p className="text-sm text-gray-500">@{user.username}</p>
                      </div>
                    </div>
                    <DropdownMenuItem asChild>
                      <Link to={`/profile/${user.username}`}>
                        <UserCircle className="mr-2 h-4 w-4" /> Profile
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link to="/change-password">
                        <Lock className="mr-2 h-4 w-4" /> Change Password
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={handleLogout}>
                      <LogOut className="mr-2 h-4 w-4" /> Logout
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
