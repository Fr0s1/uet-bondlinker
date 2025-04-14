
import React, { useState, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Image, Smile, MapPin, Loader2, X } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/AuthContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { Link } from 'react-router';
import EmojiPicker from 'emoji-picker-react';

interface UploadResponse {
  imageUrl: string;
}

const PostForm = ({ onPostCreated }: { onPostCreated?: () => void }) => {
  const [content, setContent] = useState('');
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [image, setImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { toast } = useToast();
  const { user, isAuthenticated } = useAuth();
  const queryClient = useQueryClient();

  const uploadImageMutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('image', file);
      return api.post<UploadResponse>('/uploads/image', formData, true);
    }
  });

  const createPostMutation = useMutation({
    mutationFn: (postData: { content: string, image?: string }) =>
      api.post('/posts', postData),
    onSuccess: () => {
      setContent('');
      setImage(null);
      setImagePreview(null);
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['feed'] });

      toast({
        title: "Post created",
        description: "Your post has been shared successfully.",
      });

      if (onPostCreated) {
        onPostCreated();
      }
    },
    onError: (error) => {
      console.error("Failed to create post:", error);
      toast({
        title: "Error creating post",
        description: "Please try again later.",
        variant: "destructive",
      });
    }
  });

  const handleImageClick = () => {
    fileInputRef.current?.click();
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onloadend = () => {
      setImagePreview(reader.result as string);
    };
    reader.readAsDataURL(file);

    setImage(file);
  };

  const removeImage = () => {
    setImage(null);
    setImagePreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const onEmojiClick = (emojiObject: any) => {
    setContent(prev => prev + emojiObject.emoji);
    setShowEmojiPicker(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isAuthenticated) {
      toast({
        title: "Authentication required",
        description: "Please log in to create a post",
        variant: "destructive",
      });
      return;
    }

    if (!content.trim() && !image) {
      toast({
        title: "Cannot create empty post",
        description: "Please write something or add an image.",
        variant: "destructive",
      });
      return;
    }

    try {
      setIsUploading(true);
      let imageUrl;

      if (image) {
        const response = await uploadImageMutation.mutateAsync(image);
        imageUrl = response.imageUrl;
      }

      await createPostMutation.mutateAsync({
        content: content.trim(),
        image: imageUrl
      });
    } finally {
      setIsUploading(false);
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="bg-white rounded-xl p-4 card-shadow mb-6 animate-fade-in">
        <div className="text-center py-4">
          <h3 className="font-medium text-lg mb-2">Join the conversation</h3>
          <p className="text-gray-500 mb-4">Sign in to share your thoughts and connect with others.</p>
          <div className="flex justify-center gap-3">
            <Button asChild variant="outline">
              <Link to="/login">Log in</Link>
            </Button>
            <Button asChild className="gradient-blue">
              <Link to="/register">Sign up</Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl p-4 card-shadow mb-6 animate-fade-in">
      <form onSubmit={handleSubmit}>
        <div className="flex items-start space-x-3">
          <Avatar className="h-10 w-10 mt-1">
            <AvatarImage src={user?.avatar || "/placeholder.svg"} alt={user?.name} />
            <AvatarFallback>{user?.name?.slice(0, 2).toUpperCase() || "U"}</AvatarFallback>
          </Avatar>

          <div className="flex-1">
            <div className="relative">
              <Textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="What's on your mind?"
                className="min-h-20 resize-none border-none focus-visible:ring-0 focus-visible:ring-offset-0 p-0"
              />
              <div className="text-xs text-gray-500 mt-1">
                <span>Supports markdown and emojis</span>
              </div>

              {imagePreview && (
                <div className="relative mt-2 rounded-lg overflow-hidden">
                  <img
                    src={imagePreview}
                    alt="Post image preview"
                    className="w-full h-auto max-h-64 object-cover rounded-lg"
                  />
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    className="absolute top-2 right-2 h-8 w-8 rounded-full"
                    onClick={removeImage}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </div>

            <div className="flex items-center justify-between mt-3 pt-3 border-t">
              <div className="flex space-x-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="text-social-blue rounded-full h-9 w-9"
                  onClick={handleImageClick}
                >
                  <Image className="h-5 w-5" />
                  <input
                    type="file"
                    ref={fileInputRef}
                    className="hidden"
                    accept="image/*"
                    onChange={handleImageChange}
                  />
                </Button>

                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="text-social-blue rounded-full h-9 w-9"
                  onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                >
                  <Smile className="h-5 w-5" />
                </Button>

                {showEmojiPicker && (
                  <div className="absolute mt-10 z-10">
                    <EmojiPicker onEmojiClick={onEmojiClick} />
                  </div>
                )}

                <Button type="button" variant="ghost" size="icon" className="text-social-blue rounded-full h-9 w-9">
                  <MapPin className="h-5 w-5" />
                </Button>
              </div>

              <Button
                type="submit"
                disabled={isUploading || createPostMutation.isPending || (!content.trim() && !image)}
                className="px-5 gradient-blue"
              >
                {isUploading || createPostMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Posting...
                  </>
                ) : (
                  'Post'
                )}
              </Button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default PostForm;
