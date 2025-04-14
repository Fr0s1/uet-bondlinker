
import { Navigate, useLocation } from 'react-router';
import { useAuth } from '@/contexts/AuthContext';

const ProtectedRoute = ({ children }: React.PropsWithChildren) => {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  // If authentication is still loading, show a loading state
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-social-blue"></div>
      </div>
    );
  }

  // If not authenticated, redirect to login with the intended destination in state
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  }

  // If authenticated, render the children
  return <>{children}</>;
};

export default ProtectedRoute;
