'use client';

import { Video } from '@/lib/api';
import { VideoPlayer } from './video-player';
import { X, FileVideo, Clock } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';

interface VideoModalProps {
  video: Video;
  isOpen: boolean;
  onClose: () => void;
  apiUrl: string;
}

export function VideoModal({ video, isOpen, onClose, apiUrl }: VideoModalProps) {
  if (!isOpen) return null;

  const videoUrl = `${apiUrl}/${video.file_path}`;
  const posterUrl = video.cover_path ? `${apiUrl}/${video.cover_path}` : undefined;

  const getStatusBadge = (status: string) => {
    const variants: Record<string, any> = {
      pending: { variant: 'secondary', label: '待处理' },
      processing: { variant: 'default', label: '处理中' },
      completed: { variant: 'outline', label: '已完成', className: 'bg-green-50 text-green-700 border-green-200' },
      failed: { variant: 'destructive', label: '失败' },
    };
    const config = variants[status] || variants.pending;
    return <Badge {...config}>{config.label}</Badge>;
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 字节';
    const k = 1024;
    const sizes = ['字节', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4"
      onClick={onClose}
    >
      <div
        className="relative w-full max-w-7xl bg-white rounded-lg shadow-2xl overflow-hidden flex flex-col"
        style={{ height: '90vh', maxHeight: '90vh' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b bg-gray-50 flex-shrink-0">
          <div className="flex items-center space-x-3 flex-1 min-w-0">
            <FileVideo className="h-5 w-5 text-gray-400 flex-shrink-0" />
            <div className="min-w-0 flex-1">
              <h3 className="text-lg font-semibold truncate">{video.original_filename}</h3>
              <div className="flex items-center space-x-3 text-sm text-gray-500 mt-1">
                {getStatusBadge(video.status)}
                <span className="flex items-center">
                  <Clock className="h-3 w-3 mr-1" />
                  {Math.round(video.duration)}秒
                </span>
                <span>{formatFileSize(video.file_size)}</span>
              </div>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={onClose}
            className="flex-shrink-0"
          >
            <X className="h-5 w-5" />
          </Button>
        </div>

        {/* Video Player - 自适应高度 */}
        <div className="bg-black flex-1 overflow-hidden relative">
          <VideoPlayer
            videoUrl={videoUrl}
            posterUrl={posterUrl}
            maxHeight={1200}
          />
        </div>

        {/* Footer Info */}
        <div className="p-4 bg-gray-50 text-sm text-gray-600 flex-shrink-0">
          <p>上传时间：{new Date(video.created_at).toLocaleString('zh-CN')}</p>
        </div>
      </div>
    </div>
  );
}

