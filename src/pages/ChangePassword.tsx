
import { useEffect } from 'react';
import { useNavigate } from 'react-router';
import { useAuth } from '@/contexts/AuthContext';
import ChangePasswordForm from '@/components/ChangePasswordForm';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

const ChangePasswordPage = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, isLoading, navigate]);

  return (
    <div className=" bg-gray-50">
      <div className="container max-w-6xl mx-auto p-4">
        <Button
          variant="ghost"
          className="mb-6 flex items-center"
          onClick={() => navigate(-1)}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>

        <div className="max-w-md mx-auto">
          <h1 className="text-2xl font-bold mb-8 text-center">Account Security</h1>
          <ChangePasswordForm />
        </div>
      </div>
    </div>
  );
};

export default ChangePasswordPage;
