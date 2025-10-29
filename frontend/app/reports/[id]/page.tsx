'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { reportAPI, videoAPI, Report, Video } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { NavBar } from '@/components/nav-bar';
import { ArrowLeft, FileVideo, Clock, TrendingUp, AlertTriangle, ChevronDown, ChevronUp } from 'lucide-react';
import { VideoPlayer } from '@/components/video-player';

export default function ReportDetailPage() {
  const router = useRouter();
  const params = useParams();
  const videoId = params.id as string;

  const [report, setReport] = useState<Report | null>(null);
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isTextExpanded, setIsTextExpanded] = useState(false);
  const [isTranscriptExpanded, setIsTranscriptExpanded] = useState(false);

  useEffect(() => {
    loadReport();
  }, [videoId]);

  const loadReport = async () => {
    setLoading(true);
    try {
      const [reportRes, videoRes] = await Promise.all([
        reportAPI.getByVideoId(Number(videoId)),
        videoAPI.get(Number(videoId)),
      ]);
      setReport(reportRes.data);
      setVideo(videoRes.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load report');
    } finally {
      setLoading(false);
    }
  };

  const getRiskBadge = (riskLevel: string) => {
    const variants: Record<string, any> = {
      low: { variant: 'outline', label: '低风险', className: 'bg-green-50 text-green-700 border-green-200' },
      medium: { variant: 'outline', label: '中风险', className: 'bg-yellow-50 text-yellow-700 border-yellow-200' },
      high: { variant: 'destructive', label: '高风险' },
    };
    const config = variants[riskLevel] || variants.low;
    return <Badge {...config}>{config.label}</Badge>;
  };

  const getSentimentBadge = (label: string) => {
    const variants: Record<string, any> = {
      positive: { variant: 'outline', label: '积极', className: 'bg-green-50 text-green-700 border-green-200' },
      neutral: { variant: 'secondary', label: '中性' },
      negative: { variant: 'destructive', label: '消极' },
    };
    const config = variants[label] || variants.neutral;
    return <Badge {...config}>{config.label}</Badge>;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavBar />
        <div className="container mx-auto px-4 py-8">
          <div className="text-center">加载报告中...</div>
        </div>
      </div>
    );
  }

  if (error || !report || !video) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavBar />
        <div className="container mx-auto px-4 py-8">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-red-500 mb-4">{error || '未找到报告'}</p>
              <Button onClick={() => router.push('/videos')}>
                返回视频列表
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  let keyTopics: string[] = [];
  let recommendations: string[] = [];
  
  try {
    keyTopics = JSON.parse(report.key_topics);
  } catch (e) {
    keyTopics = [];
  }

  try {
    recommendations = JSON.parse(report.recommendations);
  } catch (e) {
    recommendations = [];
  }

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar />
      <div className="container mx-auto px-4 py-8">
        <Button
          variant="ghost"
          className="mb-6"
          onClick={() => router.push('/videos')}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          返回视频列表
        </Button>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Video Info */}
            <Card>
              <CardHeader>
                <div className="flex justify-between items-start">
                  <div className="flex items-center space-x-3">
                    <FileVideo className="h-6 w-6 text-gray-400" />
                    <div>
                      <CardTitle>{video.original_filename}</CardTitle>
                      <CardDescription className="flex items-center space-x-3 mt-1">
                        <span className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          {Math.round(video.duration)}秒
                        </span>
                        <span>
                          处理耗时 {report.processing_time.toFixed(2)}秒
                        </span>
                      </CardDescription>
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="relative aspect-video bg-gray-100 rounded-lg overflow-hidden">
                  <VideoPlayer
                    videoUrl={`${API_URL}/${video.file_path}`}
                    posterUrl={video.cover_path ? `${API_URL}/${video.cover_path}` : undefined}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Cover Text */}
            <Card>
              <CardHeader>
                <CardTitle>封面文字</CardTitle>
                <CardDescription>
                  使用AI从视频封面提取的文本
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="relative">
                  <p 
                    className={`text-sm whitespace-pre-wrap bg-gray-50 p-4 rounded-md transition-all ${
                      isTextExpanded ? '' : 'line-clamp-3'
                    }`}
                  >
                    {report.cover_text}
                  </p>
                  {!isTextExpanded && report.cover_text.length > 100 && (
                    <div className="absolute bottom-0 left-0 right-0 h-12 bg-gradient-to-t from-gray-50 to-transparent rounded-b-md pointer-events-none" />
                  )}
                </div>
                {report.cover_text.length > 100 && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setIsTextExpanded(!isTextExpanded)}
                    className="w-full"
                  >
                    {isTextExpanded ? (
                      <>
                        <ChevronUp className="h-4 w-4 mr-1" />
                        收起
                      </>
                    ) : (
                      <>
                        <ChevronDown className="h-4 w-4 mr-1" />
                        展开全部
                      </>
                    )}
                  </Button>
                )}
              </CardContent>
            </Card>

            {/* Audio Transcript */}
            {report.transcript_text && (
              <Card>
                <CardHeader>
                  <CardTitle>音频转录</CardTitle>
                  <CardDescription>
                    使用Whisper AI从视频音频提取的文本
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="relative">
                    <p 
                      className={`text-sm whitespace-pre-wrap bg-blue-50 p-4 rounded-md transition-all ${
                        isTranscriptExpanded ? '' : 'line-clamp-3'
                      }`}
                    >
                      {report.transcript_text}
                    </p>
                    {!isTranscriptExpanded && report.transcript_text.length > 100 && (
                      <div className="absolute bottom-0 left-0 right-0 h-12 bg-gradient-to-t from-blue-50 to-transparent rounded-b-md pointer-events-none" />
                    )}
                  </div>
                  {report.transcript_text.length > 100 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setIsTranscriptExpanded(!isTranscriptExpanded)}
                      className="w-full"
                    >
                      {isTranscriptExpanded ? (
                        <>
                          <ChevronUp className="h-4 w-4 mr-1" />
                          收起
                        </>
                      ) : (
                        <>
                          <ChevronDown className="h-4 w-4 mr-1" />
                          展开全部
                        </>
                      )}
                    </Button>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Detailed Analysis */}
            <Card>
              <CardHeader>
                <CardTitle>详细分析</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-700 leading-relaxed">
                  {report.detailed_analysis}
                </p>
              </CardContent>
            </Card>

            {/* Recommendations */}
            {recommendations.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>建议</CardTitle>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-2">
                    {recommendations.map((rec, index) => (
                      <li key={index} className="flex items-start space-x-2">
                        <span className="text-primary mt-1">•</span>
                        <span className="text-sm text-gray-700">{rec}</span>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Sentiment Score */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <TrendingUp className="h-5 w-5 mr-2" />
                  情感倾向
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  {getSentimentBadge(report.sentiment_label)}
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-2">
                    <span>得分</span>
                    <span className="font-semibold">
                      {(report.sentiment_score * 100).toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={report.sentiment_score * 100} />
                </div>
                <p className="text-xs text-gray-500">
                  0% = 非常消极，50% = 中性，100% = 非常积极
                </p>
              </CardContent>
            </Card>

            {/* Risk Level */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <AlertTriangle className="h-5 w-5 mr-2" />
                  风险评估
                </CardTitle>
              </CardHeader>
              <CardContent>
                {getRiskBadge(report.risk_level)}
              </CardContent>
            </Card>

            {/* Key Topics */}
            {keyTopics.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>关键主题</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {keyTopics.map((topic, index) => (
                      <Badge key={index} variant="secondary">
                        {topic}
                      </Badge>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Metadata */}
            <Card>
              <CardHeader>
                <CardTitle>报告元数据</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">报告ID</span>
                  <span className="font-mono">#{report.id}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">视频ID</span>
                  <span className="font-mono">#{video.id}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">生成时间</span>
                  <span>{new Date(report.created_at).toLocaleString('zh-CN')}</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

