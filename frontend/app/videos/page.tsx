'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { videoAPI, Video } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { NavBar } from '@/components/nav-bar';
import { FileVideo, Clock, Trash2, Eye, RefreshCw } from 'lucide-react';

export default function VideosPage() {
  const router = useRouter();
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [filter, setFilter] = useState<string>('');

  useEffect(() => {
    loadVideos();
    const interval = setInterval(() => {
      loadVideos(true); // Auto-refresh for status updates
    }, 5000);
    return () => clearInterval(interval);
  }, [page, filter]);

  const loadVideos = async (silent = false) => {
    if (!silent) setLoading(true);
    try {
      const params: any = { page, page_size: 12 };
      if (filter) params.status = filter;

      const response = await videoAPI.list(params);
      setVideos(response.data.videos || []);
      setTotal(response.data.total || 0);
    } catch (err) {
      console.error('Failed to load videos:', err);
    } finally {
      if (!silent) setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this video?')) return;

    try {
      await videoAPI.delete(id);
      loadVideos();
    } catch (err) {
      alert('Failed to delete video');
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: Record<string, any> = {
      pending: { variant: 'secondary', label: 'Pending' },
      processing: { variant: 'default', label: 'Processing' },
      completed: { variant: 'outline', label: 'Completed', className: 'bg-green-50 text-green-700 border-green-200' },
      failed: { variant: 'destructive', label: 'Failed' },
    };
    const config = variants[status] || variants.pending;
    return <Badge {...config}>{config.label}</Badge>;
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleString();
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavBar />
        <div className="container mx-auto px-4 py-8">
          <div className="text-center">Loading...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar />
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold">My Videos</h1>
            <p className="text-gray-600 mt-1">
              {total} video{total !== 1 ? 's' : ''} in total
            </p>
          </div>
          <div className="flex space-x-3">
            <Button variant="outline" onClick={() => loadVideos()}>
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
            <Button onClick={() => router.push('/upload')}>
              Upload Videos
            </Button>
          </div>
        </div>

        {/* Filters */}
        <div className="flex space-x-2 mb-6">
          <Button
            variant={filter === '' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilter('')}
          >
            All
          </Button>
          <Button
            variant={filter === 'pending' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilter('pending')}
          >
            Pending
          </Button>
          <Button
            variant={filter === 'processing' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilter('processing')}
          >
            Processing
          </Button>
          <Button
            variant={filter === 'completed' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilter('completed')}
          >
            Completed
          </Button>
          <Button
            variant={filter === 'failed' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setFilter('failed')}
          >
            Failed
          </Button>
        </div>

        {/* Video Grid */}
        {videos.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <FileVideo className="mx-auto h-12 w-12 text-gray-400 mb-4" />
              <p className="text-gray-600">No videos found</p>
              <Button
                className="mt-4"
                onClick={() => router.push('/upload')}
              >
                Upload Your First Video
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {videos.map((video) => (
              <Card key={video.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex justify-between items-start mb-2">
                    <FileVideo className="h-8 w-8 text-gray-400" />
                    {getStatusBadge(video.status)}
                  </div>
                  <CardTitle className="text-lg truncate">
                    {video.original_filename}
                  </CardTitle>
                  <CardDescription>
                    <div className="flex items-center space-x-4 text-xs mt-2">
                      <span>{formatFileSize(video.file_size)}</span>
                      {video.duration > 0 && (
                        <span className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          {Math.round(video.duration)}s
                        </span>
                      )}
                    </div>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-xs text-gray-500 mb-4">
                    Uploaded: {formatDate(video.created_at)}
                  </p>
                  <div className="flex space-x-2">
                    {video.status === 'completed' && (
                      <Button
                        size="sm"
                        className="flex-1"
                        onClick={() => router.push(`/reports/${video.id}`)}
                      >
                        <Eye className="h-4 w-4 mr-1" />
                        View Report
                      </Button>
                    )}
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleDelete(video.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Pagination */}
        {total > 12 && (
          <div className="flex justify-center space-x-2 mt-8">
            <Button
              variant="outline"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              Previous
            </Button>
            <span className="flex items-center px-4">
              Page {page} of {Math.ceil(total / 12)}
            </span>
            <Button
              variant="outline"
              onClick={() => setPage((p) => p + 1)}
              disabled={page >= Math.ceil(total / 12)}
            >
              Next
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

