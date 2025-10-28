'use client';

import { useEffect, useRef } from 'react';
import Plyr from 'plyr';
import 'plyr/dist/plyr.css';

interface VideoPlayerProps {
  videoUrl: string;
  posterUrl?: string;
  autoplay?: boolean;
  controls?: string[];
  className?: string;
}

export function VideoPlayer({
  videoUrl,
  posterUrl,
  autoplay = false,
  controls = ['play-large', 'play', 'progress', 'current-time', 'mute', 'volume', 'fullscreen'],
  className = '',
}: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<Plyr | null>(null);

  useEffect(() => {
    if (!videoRef.current) return;

    // Initialize Plyr
    playerRef.current = new Plyr(videoRef.current, {
      controls,
      autoplay,
      hideControls: true,
      resetOnEnd: true,
      keyboard: { focused: true, global: true },
      tooltips: { controls: true, seek: true },
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

    return () => {
      if (playerRef.current) {
        playerRef.current.destroy();
      }
    };
  }, [videoUrl, controls, autoplay]);

  return (
    <div className={`plyr-wrapper ${className}`}>
      <video
        ref={videoRef}
        className="plyr-video"
        poster={posterUrl}
        controls
        playsInline
      >
        <source src={videoUrl} type="video/mp4" />
        您的浏览器不支持视频播放。
      </video>
    </div>
  );
}

