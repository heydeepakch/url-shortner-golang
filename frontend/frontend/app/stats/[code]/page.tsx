'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '@/app/lib/api'; // Adjust path if necessary based on your folder structure

export default function StatsPage() {
  const { code } = useParams();
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (code) {
      fetchStats();
    }
  }, [code]);

  const fetchStats = async () => {
    try {
      const response = await api.get(`/url/${code}/stats`);
      setStats(response.data);
    } catch (err) {
      setError('Failed to load stats');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="p-8 text-center">Loading stats...</div>;
  if (error) return <div className="p-8 text-center text-red-500">{error}</div>;
  if (!stats) return null;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-8 flex items-center justify-center">
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-8 max-w-2xl w-full">
        <h1 className="text-2xl font-bold mb-6 text-gray-900 dark:text-white">Stats for /{code}</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-indigo-50 dark:bg-indigo-900/20 p-4 rounded-lg">
                <p className="text-sm text-gray-500 dark:text-gray-400">Total Clicks</p>
                <p className="text-4xl font-bold text-indigo-600 dark:text-indigo-400">{stats.clicks}</p>
            </div>
            
             <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-500 dark:text-gray-400">Created At</p>
                <p className="text-lg font-semibold text-gray-900 dark:text-white">
                    {new Date(stats.created_at).toLocaleString()}
                </p>
            </div>

            <div className="col-span-1 md:col-span-2 bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Original URL</p>
                <a href={stats.original_url} target="_blank" className="text-indigo-600 hover:underline break-all">
                    {stats.original_url}
                </a>
            </div>
             <div className="col-span-1 md:col-span-2 bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Short URL</p>
                 <a href={stats.short_url} target="_blank" className="text-indigo-600 hover:underline">
                    {stats.short_url}
                </a>
            </div>
        </div>
      </div>
    </div>
  );
}
