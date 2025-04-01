
import React from 'react';
import Post, { PostProps } from './Post';

// Mock data for the feed
const mockPosts: PostProps[] = [
  {
    id: "p1",
    author: {
      id: "u1",
      name: "John Doe",
      username: "johndoe",
      avatar: "/placeholder.svg",
    },
    content: "Just launched my new website! Check it out and let me know what you think. I've been working on this project for months and I'm really excited to share it with everyone. #webdevelopment #design",
    createdAt: "2023-05-18T14:22:00Z",
    likes: 24,
    comments: 5,
    shares: 3,
    isLiked: true,
  },
  {
    id: "p2",
    author: {
      id: "u2",
      name: "Alice Smith",
      username: "alicesmith",
      avatar: "/placeholder.svg",
    },
    content: "Beautiful sunset today at the beach!",
    image: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2073&q=80",
    createdAt: "2023-05-17T18:45:00Z",
    likes: 56,
    comments: 8,
    shares: 12,
  },
  {
    id: "p3",
    author: {
      id: "u3",
      name: "Robert Johnson",
      username: "robertj",
      avatar: "/placeholder.svg",
    },
    content: "Just finished reading an amazing book about artificial intelligence and its future implications. Would highly recommend to anyone interested in tech!",
    createdAt: "2023-05-16T09:15:00Z",
    likes: 18,
    comments: 4,
    shares: 2,
  },
  {
    id: "p4",
    author: {
      id: "u4",
      name: "Emily Davis",
      username: "emilyd",
      avatar: "/placeholder.svg",
    },
    content: "My new photography project is coming along nicely. Here's a sneak peek!",
    image: "https://images.unsplash.com/photo-1554080353-a576cf803bda?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1887&q=80",
    createdAt: "2023-05-15T21:30:00Z",
    likes: 42,
    comments: 7,
    shares: 5,
  },
];

interface FeedProps {
  posts?: PostProps[];
}

const Feed = ({ posts = mockPosts }: FeedProps) => {
  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <Post key={post.id} {...post} />
      ))}
    </div>
  );
};

export default Feed;
