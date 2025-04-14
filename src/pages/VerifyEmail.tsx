
import React, { useEffect, useState } from 'react';
import { useLocation, Link } from 'react-router';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { api } from '@/lib/api-client';
import { CheckCircle, XCircle, Loader2 } from 'lucide-react';

const VerifyEmail: React.FC = () => {
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('Verifying your email...');
  const location = useLocation();

  useEffect(() => {
    const verifyEmail = async () => {
      try {
        const searchParams = new URLSearchParams(location.search);
        const token = searchParams.get('token');

        if (!token) {
          setStatus('error');
          setMessage('Invalid verification link. Token is missing.');
          return;
        }

        await api.get(`/auth/verify-email?token=${token}`);
        setStatus('success');
        setMessage('Your email has been successfully verified!');
      } catch (error) {
        setStatus('error');
        setMessage('Email verification failed. The link may be invalid or expired.');
      }
    };

    verifyEmail();
  }, [location.search]);

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-50 p-4">
      <Card className="w-full max-w-md shadow-lg">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Email Verification</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center text-center">
          {status === 'loading' && (
            <>
              <Loader2 className="h-16 w-16 text-social-blue animate-spin mb-4" />
              <p>{message}</p>
            </>
          )}
          {status === 'success' && (
            <>
              <CheckCircle className="h-16 w-16 text-green-500 mb-4" />
              <p className="text-lg font-medium">{message}</p>
              <p className="mt-2 text-gray-600">
                You can now enjoy all features of the platform.
              </p>
            </>
          )}
          {status === 'error' && (
            <>
              <XCircle className="h-16 w-16 text-red-500 mb-4" />
              <p className="text-lg font-medium">{message}</p>
              <p className="mt-2 text-gray-600">
                Please try again or contact support if the problem persists.
              </p>
            </>
          )}
        </CardContent>
        <CardFooter className="flex justify-center">
          <Button asChild className="w-full max-w-xs">
            <Link to="/login">Go to Login</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default VerifyEmail;
