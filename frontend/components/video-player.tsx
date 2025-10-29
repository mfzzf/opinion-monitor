'use client';

import { useEffect, useRef, useState } from 'react';
import 'plyr/dist/plyr.css';

// 动态导入 Plyr 以避免 TypeScript 导入错误
declare const Plyr: any;

interface VideoPlayerProps {
  videoUrl: string;
  posterUrl?: string;
  autoplay?: boolean;
  controls?: string[];
  className?: string;
  maxHeight?: number; // 最大高度限制（像素）
}

// 检测视频格式并返回对应的MIME类型
function getVideoMimeType(url: string): string {
  const extension = url.split('.').pop()?.toLowerCase() || '';
  const mimeTypes: Record<string, string> = {
    'mp4': 'video/mp4',
    'webm': 'video/webm',
    'ogg': 'video/ogg',
    'ogv': 'video/ogg',
    'avi': 'video/x-msvideo',
    'mov': 'video/quicktime',
    'mkv': 'video/x-matroska',
    'flv': 'video/x-flv',
    'wmv': 'video/x-ms-wmv',
    'm4v': 'video/x-m4v',
    '3gp': 'video/3gpp',
    'ts': 'video/mp2t',
  };
  return mimeTypes[extension] || 'video/mp4';
}

export function VideoPlayer({
  videoUrl,
  posterUrl,
  autoplay = false,
  controls = ['play-large', 'play', 'progress', 'current-time', 'mute', 'volume', 'fullscreen'],
  className = '',
  maxHeight = 1080, // 默认最大高度1080px
}: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<Plyr | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [aspectRatio, setAspectRatio] = useState<number>(16 / 9); // 默认16:9

  useEffect(() => {
    if (!videoRef.current) return;

    // 动态导入 Plyr
    import('plyr').then((PlyrModule) => {
      const PlyrConstructor = PlyrModule.default;
      if (!videoRef.current || playerRef.current) return;

      // 初始化 Plyr
      playerRef.current = new PlyrConstructor(videoRef.current, {
        controls,
        autoplay,
        hideControls: true,
        resetOnEnd: true,
        keyboard: { focused: true, global: true },
        tooltips: { controls: true, seek: true },
        ratio: `${Math.round(aspectRatio * 100)}:100`,
        i18n: {
          restart: '重新播放',
          rewind: '快退 {seektime}秒',
          play: '播放',
          pause: '暂停',
          fastForward: '快进 {seektime}秒',
          seek: '定位',
          seekLabel: '{currentTime} / {duration}',
          played: '已播放',
          buffered: '已缓冲',
          currentTime: '当前时间',
          duration: '总时长',
          volume: '音量',
          mute: '静音',
          unmute: '取消静音',
          enableCaptions: '启用字幕',
          disableCaptions: '禁用字幕',
          download: '下载',
          enterFullscreen: '进入全屏',
          exitFullscreen: '退出全屏',
          frameTitle: '视频播放器 - {title}',
          captions: '字幕',
          settings: '设置',
          menuBack: '返回上级菜单',
          speed: '播放速度',
          normal: '正常',
          quality: '画质',
          loop: '循环播放',
        },
      });
    });

    // 监听视频元数据加载，获取实际分辨率
    const handleLoadedMetadata = () => {
      if (videoRef.current) {
        const { videoWidth, videoHeight } = videoRef.current;
        if (videoWidth && videoHeight) {
          const ratio = videoWidth / videoHeight;
          setAspectRatio(ratio);
        }
      }
    };

    if (videoRef.current) {
      videoRef.current.addEventListener('loadedmetadata', handleLoadedMetadata);
    }

    return () => {
      if (videoRef.current) {
        videoRef.current.removeEventListener('loadedmetadata', handleLoadedMetadata);
      }
      if (playerRef.current) {
        playerRef.current.destroy();
      }
    };
  }, [videoUrl, controls, autoplay, maxHeight, aspectRatio]);

  const mimeType = getVideoMimeType(videoUrl);

  return (
    <div 
      ref={containerRef}
      className={`plyr-wrapper ${className}`}
      style={{ 
        width: '100%',
        height: '100%',
        maxHeight: `${maxHeight}px`,
        margin: '0 auto',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <video
        ref={videoRef}
        className="plyr-video"
        poster={posterUrl}
        controls
        playsInline
        style={{ 
          maxWidth: '100%',
          maxHeight: '100%',
          width: 'auto',
          height: 'auto',
          objectFit: 'contain',
        }}
      >
        <source src={videoUrl} type={mimeType} />
        {/* 提供多个备用格式 */}
        <source src={videoUrl} type="video/mp4" />
        <source src={videoUrl} type="video/webm" />
        您的浏览器不支持视频播放。请尝试使用现代浏览器（Chrome、Firefox、Safari、Edge）。
      </video>
    </div>
  );
}

