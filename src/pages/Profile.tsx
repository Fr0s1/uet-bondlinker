
import { useParams } from 'react-router';
import UserProfile from '@/components/UserProfile';
import Feed from '@/components/Feed';
import { useUserByUsername } from '@/hooks/use-users';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

const Profile = () => {
  const { username } = useParams<{ username: string }>();
  const { user: currentUser } = useAuth();
  const { data: profileUser, isLoading, error } = useUserByUsername(username || '');

  const isCurrentUser = currentUser?.id === profileUser?.id;

  // Show loading state
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-[calc(100vh-64px)]">
        <Loader2 className="h-8 w-8 animate-spin text-social-blue" />
        <span className="ml-2 text-gray-500">Loading profile...</span>
      </div>
    );
  }

  // Show error state
  if (error || !profileUser) {
    return (
      <div className="bg-white rounded-xl p-8 text-center card-shadow mt-4">
        <h2 className="text-2xl font-bold text-gray-800 mb-2">User not found</h2>
        <p className="text-gray-600">The user you're looking for doesn't exist or has been removed.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
      {/* Left Sidebar - Empty on profile page */}
      <aside className="hidden lg:block lg:col-span-3">
        {/* Intentionally empty for layout consistency */}
      </aside>

      {/* Main Content */}
      <main className="lg:col-span-6 pt-4">
        <UserProfile
          user={{
            id: profileUser.id,
            name: profileUser.name,
            username: profileUser.username,
            avatar: profileUser.avatar || "/user-avatar.png",
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
          <Feed type="public" userId={profileUser.id} />
        </div>
      </main>

      {/* Right Sidebar - Empty on profile page */}
      <aside className="hidden lg:block lg:col-span-3">
        {/* Intentionally empty for layout consistency */}
      </aside>
    </div>
  );
};

export default Profile;
