'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { reportAPI, videoAPI, Report, Video } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { NavBar } from '@/components/nav-bar';
import { ArrowLeft, FileVideo, Clock, TrendingUp, AlertTriangle } from 'lucide-react';
import Image from 'next/image';

export default function ReportDetailPage() {
  const router = useRouter();
  const params = useParams();
  const videoId = params.id as string;

  const [report, setReport] = useState<Report | null>(null);
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

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
      low: { variant: 'outline', label: 'Low Risk', className: 'bg-green-50 text-green-700 border-green-200' },
      medium: { variant: 'outline', label: 'Medium Risk', className: 'bg-yellow-50 text-yellow-700 border-yellow-200' },
      high: { variant: 'destructive', label: 'High Risk' },
    };
    const config = variants[riskLevel] || variants.low;
    return <Badge {...config}>{config.label}</Badge>;
  };

  const getSentimentBadge = (label: string) => {
    const variants: Record<string, any> = {
      positive: { variant: 'outline', label: 'Positive', className: 'bg-green-50 text-green-700 border-green-200' },
      neutral: { variant: 'secondary', label: 'Neutral' },
      negative: { variant: 'destructive', label: 'Negative' },
    };
    const config = variants[label] || variants.neutral;
    return <Badge {...config}>{config.label}</Badge>;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavBar />
        <div className="container mx-auto px-4 py-8">
          <div className="text-center">Loading report...</div>
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
              <p className="text-red-500 mb-4">{error || 'Report not found'}</p>
              <Button onClick={() => router.push('/videos')}>
                Back to Videos
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
          Back to Videos
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
                          {Math.round(video.duration)}s
                        </span>
                        <span>
                          Processed in {report.processing_time.toFixed(2)}s
                        </span>
                      </CardDescription>
                    </div>
                  </div>
                </div>
              </CardHeader>
              {video.cover_path && (
                <CardContent>
                  <div className="relative aspect-video bg-gray-100 rounded-lg overflow-hidden">
                    <Image
                      src={`${API_URL}/${video.cover_path}`}
                      alt="Video cover"
                      fill
                      className="object-contain"
                    />
                  </div>
                </CardContent>
              )}
            </Card>

            {/* Cover Text */}
            <Card>
              <CardHeader>
                <CardTitle>Extracted Text</CardTitle>
                <CardDescription>
                  Text extracted from video cover using AI
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm whitespace-pre-wrap bg-gray-50 p-4 rounded-md">
                  {report.cover_text}
                </p>
              </CardContent>
            </Card>

            {/* Detailed Analysis */}
            <Card>
              <CardHeader>
                <CardTitle>Detailed Analysis</CardTitle>
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
                  <CardTitle>Recommendations</CardTitle>
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
                  Sentiment
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  {getSentimentBadge(report.sentiment_label)}
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-2">
                    <span>Score</span>
                    <span className="font-semibold">
                      {(report.sentiment_score * 100).toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={report.sentiment_score * 100} />
                </div>
                <p className="text-xs text-gray-500">
                  0% = Very Negative, 50% = Neutral, 100% = Very Positive
                </p>
              </CardContent>
            </Card>

            {/* Risk Level */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <AlertTriangle className="h-5 w-5 mr-2" />
                  Risk Assessment
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
                  <CardTitle>Key Topics</CardTitle>
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
                <CardTitle>Report Metadata</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Report ID</span>
                  <span className="font-mono">#{report.id}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Video ID</span>
                  <span className="font-mono">#{video.id}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Generated</span>
                  <span>{new Date(report.created_at).toLocaleString()}</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

