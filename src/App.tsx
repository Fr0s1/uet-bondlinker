
import { BrowserRouter as Router, Routes, Route } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '@/contexts/AuthContext';
import { Toaster } from '@/components/ui/toaster';

import Index from '@/pages/Index';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import ForgotPassword from '@/pages/ForgotPassword';
import ResetPassword from '@/pages/ResetPassword';
import VerifyEmail from '@/pages/VerifyEmail';
import Profile from '@/pages/Profile';
import Search from '@/pages/Search';
import Trending from '@/pages/Trending';
import Messages from '@/pages/Messages';
import ChangePassword from '@/pages/ChangePassword';
import NotFound from '@/pages/NotFound';

import './App.css';
import { MainLayout } from './components/MainLayout';
import { FCMProvider } from './components/FCMProvider';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <FCMProvider>
          <Router>
            <div className="min-h-screen flex flex-col">
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/forgot-password" element={<ForgotPassword />} />
                <Route path="/reset-password" element={<ResetPassword />} />
                <Route path="/verify-email" element={<VerifyEmail />} />
                <Route element={<MainLayout />}>
                  <Route path="/" element={<Index />} />
                  <Route path="/profile/:username" element={<Profile />} />
                  <Route path="/search" element={<Search />} />
                  <Route path="/trending" element={<Trending />} />
                  <Route path="/messages" element={<Messages />} />
                  <Route path="/messages/c/:conversationId" element={<Messages />} />
                  <Route path="/change-password" element={<ChangePassword />} />
                </Route>

                <Route path="*" element={<NotFound />} />
              </Routes>
              <Toaster />
            </div>
          </Router>
        </FCMProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
