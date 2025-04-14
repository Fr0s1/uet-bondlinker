
import React, { useState } from 'react';
import { Link } from 'react-router';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Loader2, Mail, ArrowLeft } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { api } from '@/lib/api-client';

const schema = z.object({
  email: z.string().email({ message: 'Please enter a valid email address' }),
});

type FormValues = z.infer<typeof schema>;

const ForgotPassword: React.FC = () => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const { toast } = useToast();

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
    },
  });

  const onSubmit = async (data: FormValues) => {
    setIsSubmitting(true);
    try {
      await api.post('/auth/forgot-password', data);
      setIsSubmitted(true);
      toast({
        title: 'Reset link sent',
        description: 'If your email is registered, you will receive a password reset link',
      });
    } catch (error) {
      // Don't show error for security reasons
      setIsSubmitted(true);
      toast({
        title: 'Reset link sent',
        description: 'If your email is registered, you will receive a password reset link',
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-50 p-4">
      <Card className="w-full max-w-md shadow-lg">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Forgot Password</CardTitle>
          <CardDescription className="text-center">
            {!isSubmitted
              ? 'Enter your email address and we\'ll send you a link to reset your password.'
              : 'Check your email for the reset link'}
          </CardDescription>
        </CardHeader>

        <CardContent>
          {!isSubmitted ? (
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="your@email.com"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button
                  type="submit"
                  className="w-full gradient-blue"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Sending reset link...
                    </>
                  ) : (
                    'Send Reset Link'
                  )}
                </Button>
              </form>
            </Form>
          ) : (
            <div className="text-center p-6">
              <Mail className="h-12 w-12 text-social-blue mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Check your inbox</h3>
              <p className="text-gray-600 mb-4">
                We've sent you an email with a link to reset your password. The link will expire in 15 minutes.
              </p>
              <p className="text-sm text-gray-500">
                Didn't receive an email? Check your spam folder or try again.
              </p>
            </div>
          )}
        </CardContent>

        <CardFooter className="flex justify-center">
          <Button asChild variant="outline" className="w-full">
            <Link to="/login" className="flex items-center justify-center">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Login
            </Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default ForgotPassword;
