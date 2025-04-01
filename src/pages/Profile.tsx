
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from '@/components/Navbar';
import UserProfile from '@/components/UserProfile';
import Feed from '@/components/Feed';
import { useUserByUsername } from '@/hooks/use-users';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api-client';

const Profile = () => {
  const { username } = useParams<{ username: string }>();
  const { user: currentUser } = useAuth();
  const { data: profileUser, isLoading, error } = useUserByUsername(username || '');
  
  // Fetch the user's posts
  const { data: posts, isLoading: isPostsLoading } = useQuery({
    queryKey: ['user-posts', profileUser?.id],
    queryFn: () => api.get(`/posts?user_id=${profileUser?.id}`),
    enabled: !!profileUser?.id,
  });
  
  const isCurrentUser = currentUser?.id === profileUser?.id;
  
  // Show loading state
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="flex justify-center items-center h-[calc(100vh-64px)]">
          <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
          <span className="ml-2 text-gray-500">Loading profile...</span>
        </div>
      </div>
    );
  }
  
  // Show error state
  if (error || !profileUser) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="container py-12">
          <div className="bg-white rounded-xl p-8 text-center card-shadow">
            <h2 className="text-2xl font-bold text-gray-800 mb-2">User not found</h2>
            <p className="text-gray-600">The user you're looking for doesn't exist or has been removed.</p>
          </div>
        </div>
      </div>
    );
  }
  
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="container py-4">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
          {/* Left Sidebar - Empty on profile page */}
          <aside className="hidden lg:block lg:col-span-3">
            {/* Intentionally empty for layout consistency */}
          </aside>
          
          {/* Main Content */}
          <main className="lg:col-span-6">
            <UserProfile 
              user={{
                id: profileUser.id,
                name: profileUser.name,
                username: profileUser.username,
                avatar: profileUser.avatar || "/placeholder.svg",
                bio: profileUser.bio || "No bio provided",
                location: profileUser.location,
                website: profileUser.website,
                joinedDate: profileUser.createdAt,
                followers: profileUser.followers || 0,
                following: profileUser.following || 0,
                isFollowing: profileUser.isFollowed,
              }}
              isCurrentUser={isCurrentUser}
            />
            
            <div className="mt-6">
              <h2 className="font-semibold text-xl mb-4 px-4">Posts</h2>
              {isPostsLoading ? (
                <div className="flex justify-center items-center py-12">
                  <Loader2 className="h-6 w-6 animate-spin text-social-blue" />
                  <span className="ml-2 text-gray-500">Loading posts...</span>
                </div>
              ) : posts && posts.length > 0 ? (
                <Feed />
              ) : (
                <div className="bg-white rounded-xl p-8 text-center card-shadow">
                  <h3 className="text-lg font-medium text-gray-700">No posts yet</h3>
                  <p className="text-gray-500 mt-2">
                    {isCurrentUser 
                      ? "You haven't created any posts yet. Start sharing!"
                      : `${profileUser.name} hasn't shared any posts yet.`
                    }
                  </p>
                </div>
              )}
            </div>
          </main>
          
          {/* Right Sidebar - Empty on profile page */}
          <aside className="hidden lg:block lg:col-span-3">
            {/* Intentionally empty for layout consistency */}
          </aside>
        </div>
      </div>
    </div>
  );
};

export default Profile;
