'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import api from '../lib/api';

// Define the type for a URL object based on your Go model
interface URLData {
  id: number;
  short_code: string;
  original_url: string;
  clicks: number;
  created_at: string;
  short_url: string; // Assuming backend returns this or we construct it
}

export default function Dashboard() {
  const [urls, setUrls] = useState<URLData[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetchUrls();
  }, [router]);

  const fetchUrls = async () => {
    try {
      const response = await api.get('/my-urls');
      // Ensure response.data is an array; if it's null or undefined, use []
      // Also handles if backend returns { urls: [...] } wrapper
      const data = response.data;
      if (Array.isArray(data)) {
          setUrls(data);
      } else if (data && Array.isArray(data.urls)) {
          setUrls(data.urls);
      } else {
          setUrls([]);
      }
    } catch (error) {
      console.error('Failed to fetch URLs', error);
      setUrls([]); // Fallback to empty array on error
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
      localStorage.removeItem('token');
      router.push('/login');
  };

  if (loading) return <div className="p-8 text-center">Loading dashboard...</div>;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">My Dashboard</h1>
            <div className="flex gap-4">
                <Link href="/" className="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">
                    Create New Link
                </Link>
                <button onClick={handleLogout} className="px-4 py-2 border border-red-600 text-red-600 rounded hover:bg-red-50">
                    Logout
                </button>
            </div>
        </div>

        <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-300">Original URL</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-300">Short Link</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-300">Clicks</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-300">Created At</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider dark:text-gray-300">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200 dark:divide-gray-700 dark:bg-gray-800">
              {urls.map((url) => (
                <tr key={url.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100 truncate max-w-xs">
                    <a href={url.original_url} target="_blank" rel="noopener noreferrer" className="hover:text-indigo-600">
                        {url.original_url}
                    </a>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-indigo-600 hover:text-indigo-900">
                    <a href={url.short_url} target="_blank" rel="noopener noreferrer">{url.short_code}</a>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {url.clicks}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {new Date(url.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link href={`/stats/${url.short_code}`} className="text-indigo-600 hover:text-indigo-900 mr-4">
                      Stats
                    </Link>
                  </td>
                </tr>
              ))}
              {urls.length === 0 && (
                  <tr>
                      <td colSpan={5} className="px-6 py-4 text-center text-gray-500">No URLs found. Create one!</td>
                  </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
