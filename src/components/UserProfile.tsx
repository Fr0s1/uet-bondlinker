
import React from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { MapPin, Calendar, Link as LinkIcon } from 'lucide-react';

interface UserProfileProps {
  user: {
    id: string;
    name: string;
    username: string;
    avatar: string;
    bio: string;
    followers: number;
    following: number;
    location?: string;
    website?: string;
    joinedDate: string;
    isFollowing?: boolean;
  };
  isCurrentUser?: boolean;
}

const UserProfile = ({ user, isCurrentUser = false }: UserProfileProps) => {
  const [isFollowing, setIsFollowing] = React.useState(user.isFollowing || false);
  const [followerCount, setFollowerCount] = React.useState(user.followers);
  
  const handleFollowToggle = () => {
    if (isFollowing) {
      setFollowerCount(prev => prev - 1);
    } else {
      setFollowerCount(prev => prev + 1);
    }
    setIsFollowing(!isFollowing);
  };
  
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  };
  
  return (
    <div className="bg-white rounded-xl overflow-hidden card-shadow animate-fade-in">
      <div className="h-32 bg-gradient-to-r from-social-blue to-social-darkblue"></div>
      
      <div className="px-4 pb-4">
        <div className="flex justify-between items-end -mt-12 mb-4">
          <Avatar className="h-24 w-24 border-4 border-white avatar-shadow">
            <AvatarImage src={user.avatar || "/placeholder.svg"} alt={user.name} />
            <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>
          
          {isCurrentUser ? (
            <Button variant="outline" className="mb-2">Edit Profile</Button>
          ) : (
            <Button 
              variant={isFollowing ? "outline" : "default"} 
              className={isFollowing ? "mb-2" : "mb-2 gradient-blue"}
              onClick={handleFollowToggle}
            >
              {isFollowing ? 'Following' : 'Follow'}
            </Button>
          )}
        </div>
        
        <div>
          <h2 className="text-xl font-bold">{user.name}</h2>
          <p className="text-gray-500">@{user.username}</p>
          
          <p className="my-3">{user.bio}</p>
          
          <div className="flex flex-wrap text-sm text-gray-500 space-x-4 mb-3">
            {user.location && (
              <div className="flex items-center">
                <MapPin className="h-4 w-4 mr-1" />
                <span>{user.location}</span>
              </div>
            )}
            
            {user.website && (
              <div className="flex items-center">
                <LinkIcon className="h-4 w-4 mr-1" />
                <a href={user.website} target="_blank" rel="noopener noreferrer" className="text-social-blue hover:underline">
                  {user.website.replace(/(^\w+:|^)\/\//, '')}
                </a>
              </div>
            )}
            
            <div className="flex items-center">
              <Calendar className="h-4 w-4 mr-1" />
              <span>Joined {formatDate(user.joinedDate)}</span>
            </div>
          </div>
          
          <div className="flex space-x-5 text-sm">
            <Link to="#" className="hover:underline">
              <span className="font-semibold">{user.following}</span> Following
            </Link>
            <Link to="#" className="hover:underline">
              <span className="font-semibold">{followerCount}</span> Followers
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserProfile;
