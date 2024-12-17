import React from 'react';
import { HttpEntity } from '../types';
import { formatTime } from '../utils/format';

interface RequestListProps {
    requests: HttpEntity[];
    selectedRequest: HttpEntity | null;
    onSelectRequest: (request: HttpEntity) => void;
}

export const RequestList: React.FC<RequestListProps> = ({ requests, selectedRequest, onSelectRequest }) => {
    return (
        <div className="h-full flex flex-col bg-echogy-bg-secondary dark:bg-echogy-bg-secondary-dark rounded-lg border border-echogy-border dark:border-echogy-border-dark">
            {/* Header */}
            <div className="h-11 flex-none px-3 bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark flex justify-between items-center border-b border-echogy-border dark:border-echogy-border-dark">
                <div className="flex space-x-2">
                    <button className="p-1.5 hover:bg-echogy-bg-hover dark:hover:bg-echogy-bg-hover-dark rounded transition-colors">
                        <svg className="w-4 h-4 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                             fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                        </svg>
                    </button>
                    <button className="p-1.5 hover:bg-echogy-bg-hover dark:hover:bg-echogy-bg-hover-dark rounded transition-colors">
                        <svg className="w-4 h-4 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                             fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/>
                        </svg>
                    </button>
                </div>
            </div>
            
            {/* Table Container */}
            <div className="flex-1 min-h-0 overflow-auto relative">
                <table className="w-full">
                    <thead className="sticky top-0 bg-echogy-bg-secondary dark:bg-echogy-bg-secondary-dark">
                        <tr className="text-xs text-echogy-text-secondary dark:text-echogy-text-secondary-dark border-b border-echogy-border dark:border-echogy-border-dark">
                            <th className="px-4 py-2 font-medium text-left w-16">#</th>
                            <th className="px-4 py-2 font-medium text-left w-24">Method</th>
                            <th className="px-4 py-2 font-medium text-left min-w-[300px]">URI</th>
                            <th className="px-4 py-2 font-medium text-left w-24">Status</th>
                            <th className="px-4 py-2 font-medium text-left w-24">Time</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-echogy-border dark:divide-echogy-border-dark">
                    {requests.map((entity, index) => (
                        <tr
                            key={index}
                            className={`border-b border-echogy-border dark:border-echogy-border-dark hover:bg-echogy-bg-hover dark:hover:bg-echogy-bg-hover-dark cursor-pointer ${
                                selectedRequest?.id === entity.id ? 'bg-echogy-bg-hover dark:bg-echogy-bg-hover-dark' : ''
                            }`}
                            onClick={() => onSelectRequest(entity)}
                        >
                            <td className="px-4 py-2 text-sm">{index + 1}</td>
                            <td className="px-4 py-2 text-sm">{entity.request.method}</td>
                            <td className="px-4 py-2 text-sm font-mono">{entity.request.uri}</td>
                            <td className="px-4 py-2">
                                <span className={`px-2 py-0.5 text-xs rounded ${
                                    entity.response.status >= 200 && entity.response.status < 300
                                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                        : entity.response.status >= 400
                                            ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                                            : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200'
                                }`}>
                                    {entity.response.status}
                                </span>
                            </td>
                            <td className="px-4 py-2 text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                {formatTime(entity.useTime)}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};
