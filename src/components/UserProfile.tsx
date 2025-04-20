
import React, { useState, useRef } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router';
import { MapPin, Calendar, Link as LinkIcon, Upload, UserPlus, UserCheck } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { format } from 'date-fns';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { toast } from '@/components/ui/use-toast';
import { useFollowUser, useUnfollowUser } from '@/hooks/use-users';
import ChatButton from './ChatButton';
import { User } from '@/contexts/AuthContext';

interface UserProfileProps {
  user: User;
  isCurrentUser?: boolean;
}

const UserProfile = ({ user, isCurrentUser = false }: UserProfileProps) => {
  const [isFollowing, setIsFollowing] = React.useState(user.isFollowed || false);
  const [followerCount, setFollowerCount] = React.useState(user.followers);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const avatarInputRef = useRef<HTMLInputElement>(null);
  const coverInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  const followMutation = useFollowUser();
  const unfollowMutation = useUnfollowUser();

  const handleFollowToggle = () => {
    if (isFollowing) {
      unfollowMutation.mutate(user.id, {
        onSuccess: () => {
          setIsFollowing(false);
          setFollowerCount(prev => prev - 1);
        }
      });
    } else {
      followMutation.mutate(user.id, {
        onSuccess: () => {
          setIsFollowing(true);
          setFollowerCount(prev => prev + 1);
        }
      });
    }
  };

  // File upload mutation for S3
  const fileUploadMutation = useMutation({
    mutationFn: (file: File) => {
      const formData = new FormData();
      formData.append("file", file);
      return api.post<{ url: string }>("/uploads", formData, true);
    },
    onError: () => {
      toast({
        title: "Error",
        description: "Failed to upload image. Please try again.",
        variant: "destructive",
      });
      setIsUploading(false);
    }
  });

  const updateProfileMutation = useMutation({
    mutationFn: (data: { avatar?: string, cover?: string }) =>
      api.put<any>(`/users/${user.id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', 'username', user.username] });
      queryClient.invalidateQueries({ queryKey: ['auth'] });
      setIsEditDialogOpen(false);
      setIsUploading(false);
      toast({
        title: "Profile updated",
        description: "Your profile has been updated successfully!",
      });
    },
    onError: () => {
      toast({
        title: "Error",
        description: "Failed to update profile. Please try again.",
        variant: "destructive",
      });
      setIsUploading(false);
    }
  });

  const handleAvatarUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    toast({
      title: "Uploading...",
      description: "Your avatar is being uploaded, please wait.",
    });

    // Upload file to S3 via our backend
    fileUploadMutation.mutate(file, {
      onSuccess: ({ url }) => {
        // After successful S3 upload, update the user profile with the new avatar URL
        updateProfileMutation.mutate({
          avatar: url
        });
      }
    });
  };

  const handleCoverUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    toast({
      title: "Uploading...",
      description: "Your cover image is being uploaded, please wait.",
    });

    // Upload file to S3 via our backend
    fileUploadMutation.mutate(file, {
      onSuccess: ({ url }) => {
        // After successful S3 upload, update the user profile with the new avatar URL
        updateProfileMutation.mutate({
          cover: url
        });
      }
    });
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return format(date, 'MMMM yyyy');
    } catch (error) {
      return "Unknown date";
    }
  };

  return (
    <div className="bg-white rounded-xl overflow-hidden card-shadow animate-fade-in">
      <div className="relative">
        <div className="h-64 bg-gradient-to-r from-social-blue to-social-darkblue" style={{
          backgroundImage: user.cover ? `url(${user.cover})` : null,
          backgroundRepeat: 'no-repeat',
          backgroundPosition: 'center',
          backgroundSize: 'cover'
        }}></div>
        {isCurrentUser && (
          <Button
            variant="outline"
            size="icon"
            className="absolute top-2 right-2 bg-white/80 hover:bg-white"
            onClick={() => coverInputRef.current?.click()}
            disabled={isUploading}
          >
            <Upload className="h-4 w-4" />
            <input
              type="file"
              ref={coverInputRef}
              className="hidden"
              accept="image/*"
              onChange={handleCoverUpload}
            />
          </Button>
        )}
      </div>

      <div className="px-4 pb-4">
        <div className="flex justify-between items-end -mt-14 mb-4">
          <div className="relative">
            <Avatar className="h-28 w-28 border-4 border-white avatar-shadow bg-white">
              <AvatarImage src={user.avatar} alt={user.name} />
              <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
            </Avatar>
            {isCurrentUser && (
              <Button
                variant="outline"
                size="icon"
                className="absolute bottom-0 right-0 rounded-full w-7 h-7 bg-white"
                onClick={() => avatarInputRef.current?.click()}
                disabled={isUploading}
              >
                <Upload className="h-3 w-3" />
                <input
                  type="file"
                  ref={avatarInputRef}
                  className="hidden"
                  accept="image/*"
                  onChange={handleAvatarUpload}
                />
              </Button>
            )}
          </div>

          <div className="flex space-x-2 mb-2">
            {isCurrentUser ? (
              <Button variant="outline" onClick={() => setIsEditDialogOpen(true)}>
                Edit Profile
              </Button>
            ) : (
              <>
                <Button
                  variant={isFollowing ? "outline" : "default"}
                  className={isFollowing ? "" : "gradient-blue"}
                  onClick={handleFollowToggle}
                  disabled={followMutation.isPending || unfollowMutation.isPending}
                >
                  {isFollowing ? (
                    <>
                      <UserCheck className="h-4 w-4 mr-1" />
                      Following
                    </>
                  ) : (
                    <>
                      <UserPlus className="h-4 w-4 mr-1" />
                      Follow
                    </>
                  )}
                </Button>
                <ChatButton username={user.username} userId={user.id} />
              </>
            )}
          </div>
        </div>

        <div>
          <h2 className="text-xl font-bold text-left">{user.name}</h2>
          <p className="text-gray-500 text-left">@{user.username}</p>

          <p className="my-3 text-left">{user.bio}</p>

          <div className="flex flex-wrap text-sm text-gray-500 space-x-4 mb-3 text-left">
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
              <span>Joined {formatDate(user.createdAt)}</span>
            </div>
          </div>

          <div className="flex space-x-5 text-sm text-left">
            <Link to="#" className="hover:underline">
              <span className="font-semibold">{user.following}</span> Following
            </Link>
            <Link to="#" className="hover:underline">
              <span className="font-semibold">{followerCount}</span> Followers
            </Link>
          </div>
        </div>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Profile</DialogTitle>
            <DialogDescription>
              Update your profile information
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="flex justify-center">
              <div className="relative">
                <Avatar className="h-20 w-20 border-2 border-white avatar-shadow">
                  <AvatarImage src={user.avatar} alt={user.name} />
                  <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <Button
                  variant="outline"
                  size="icon"
                  className="absolute bottom-0 right-0 rounded-full w-7 h-7 bg-white"
                  onClick={() => avatarInputRef.current?.click()}
                  disabled={isUploading}
                >
                  <Upload className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <p className="text-center text-sm text-gray-500">Click on the icon to change your profile picture</p>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default UserProfile;
