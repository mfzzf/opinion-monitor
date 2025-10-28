'use client';

import Link from 'next/link';
import { useAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { usePathname } from 'next/navigation';

export function NavBar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();

  if (!user) return null;

  const isActive = (path: string) => pathname === path;

  return (
    <nav className="border-b bg-white">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center space-x-8">
            <Link href="/videos" className="text-xl font-bold">
              舆情监测
            </Link>
            <div className="flex space-x-4">
              <Link
                href="/videos"
                className={`px-3 py-2 rounded-md text-sm font-medium ${
                  isActive('/videos')
                    ? 'bg-gray-100 text-gray-900'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                视频
              </Link>
              <Link
                href="/upload"
                className={`px-3 py-2 rounded-md text-sm font-medium ${
                  isActive('/upload')
                    ? 'bg-gray-100 text-gray-900'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                上传
              </Link>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <span className="text-sm text-gray-600">
              欢迎，{user.username}
            </span>
            <Button variant="outline" size="sm" onClick={logout}>
              退出登录
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}

