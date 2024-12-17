import React from 'react';
import { HttpEntity } from '../types';

interface RequestDetailProps {
    request: HttpEntity | null;
    tab: 'request' | 'response';
    onTabChange: (tab: 'request' | 'response') => void;
}

export const RequestDetail: React.FC<RequestDetailProps> = ({ request, tab, onTabChange }) => {
    return (
        <div className="h-full flex flex-col bg-echogy-bg-secondary dark:bg-echogy-bg-secondary-dark rounded-lg overflow-hidden border border-echogy-border dark:border-echogy-border-dark">
            <div className="h-11 px-3 bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark flex items-center border-b border-echogy-border dark:border-echogy-border-dark">
                <button
                    className={`px-4 py-1.5 text-sm font-medium rounded-md transition-colors ${
                        tab === 'request'
                            ? 'bg-echogy-bg-hover dark:bg-echogy-bg-hover-dark text-echogy-text-primary dark:text-echogy-text-primary-dark'
                            : 'text-echogy-text-secondary dark:text-echogy-text-secondary-dark hover:text-echogy-text-primary dark:hover:text-echogy-text-primary-dark hover:bg-echogy-bg-hover/50 dark:hover:bg-echogy-bg-hover-dark/50'
                    }`}
                    onClick={() => onTabChange('request')}
                >
                    Request
                </button>
                <button
                    className={`px-4 py-1.5 text-sm font-medium rounded-md ml-2 transition-colors ${
                        tab === 'response'
                            ? 'bg-echogy-bg-hover dark:bg-echogy-bg-hover-dark text-echogy-text-primary dark:text-echogy-text-primary-dark'
                            : 'text-echogy-text-secondary dark:text-echogy-text-secondary-dark hover:text-echogy-text-primary dark:hover:text-echogy-text-primary-dark hover:bg-echogy-bg-hover/50 dark:hover:bg-echogy-bg-hover-dark/50'
                    }`}
                    onClick={() => onTabChange('response')}
                >
                    Response
                </button>
            </div>

            <div className="flex-1 overflow-auto p-4">
                {request ? (
                    <div className="space-y-4">
                        <div className="flex justify-end items-center mb-4">
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                {request.request.time}
                            </div>
                        </div>
                        <div className="font-mono text-sm space-y-2">
                            {tab === 'request' ? (
                                /* Request Headers */
                                Object.entries(request.request.headers).map(([key, value]) => (
                                    <div key={key} className="flex">
                                        <span className="text-echogy-text-secondary dark:text-echogy-text-secondary-dark w-32">
                                            {key}:
                                        </span>
                                        <span className="text-echogy-text-primary dark:text-echogy-text-primary-dark flex-1">
                                            {value}
                                        </span>
                                    </div>
                                ))
                            ) : (
                                /* Response Headers */
                                <>
                                    <div className="flex mb-4">
                                        <span className="text-echogy-text-secondary dark:text-echogy-text-secondary-dark w-32">
                                            Status:
                                        </span>
                                        <span className={`text-sm rounded-md px-2 py-0.5 ${
                                            request.response.status >= 200 && request.response.status < 300
                                                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                                : request.response.status >= 400
                                                    ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                                                    : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200'
                                        }`}>
                                            {request.response.status}
                                        </span>
                                    </div>
                                    {Object.entries(request.response.headers).map(([key, value]) => (
                                        <div key={key} className="flex">
                                            <span className="text-echogy-text-secondary dark:text-echogy-text-secondary-dark w-32">
                                                {key}:
                                            </span>
                                            <span className="text-echogy-text-primary dark:text-echogy-text-primary-dark flex-1">
                                                {value}
                                            </span>
                                        </div>
                                    ))}
                                </>
                            )}
                        </div>
                    </div>
                ) : (
                    <div className="h-full flex items-center justify-center text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                        Select a request to view details
                    </div>
                )}
            </div>
        </div>
    );
};
