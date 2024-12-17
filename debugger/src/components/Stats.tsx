import React from 'react';
import { Stats as StatsType } from '../types';

interface StatsProps {
    stats: StatsType;
}

export const Stats: React.FC<StatsProps> = ({ stats }) => {
    return (
        <div className="h-full flex flex-col">
            <div className="px-4 py-2">
                <span className="text-sm font-medium text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                    Stats
                </span>
            </div>

            <div className="px-4 py-4">
                <div className="grid grid-cols-2 gap-6">
                    <div className="flex items-start space-x-3">
                        <div className="mt-1">
                            <svg className="w-5 h-5 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                      d="M13 10V3L4 14h7v7l9-11h-7z"/>
                            </svg>
                        </div>
                        <div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Active Connections
                            </div>
                            <div
                                className="text-xl font-medium text-echogy-text-primary dark:text-echogy-text-primary-dark">
                                {stats.activeConnections}
                            </div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Total: {stats.totalConnections}
                            </div>
                        </div>
                    </div>

                    <div className="flex items-start space-x-3">
                        <div className="mt-1">
                            <svg className="w-5 h-5 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                      d="M20 3h-3a2 2 0 00-2 2v14a2 2 0 002 2h3v-3.5M4 3h3a2 2 0 012 2v14a2 2 0 01-2 2H4v-3.5M8 7h12m-8 4h8m-8 4h8"/>
                            </svg>
                        </div>
                        <div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Requests
                            </div>
                            <div
                                className="text-xl font-medium text-echogy-text-primary dark:text-echogy-text-primary-dark">
                                {stats.requests}
                            </div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                {stats.requestBytes} K
                            </div>
                        </div>
                    </div>

                    <div className="flex items-start space-x-3">
                        <div className="mt-1">
                            <svg className="w-5 h-5 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                      d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
                            </svg>
                        </div>
                        <div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Traffic
                            </div>
                            <div
                                className="text-xl font-medium text-echogy-text-primary dark:text-echogy-text-primary-dark">
                                {Math.round((stats.requestBytes + stats.responseBytes) / 1024)} MB
                            </div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Total
                            </div>
                        </div>
                    </div>

                    <div className="flex items-start space-x-3">
                        <div className="mt-1">
                            <svg className="w-5 h-5 text-echogy-text-secondary dark:text-echogy-text-secondary-dark"
                                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                            </svg>
                        </div>
                        <div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                Responses
                            </div>
                            <div
                                className="text-xl font-medium text-echogy-text-primary dark:text-echogy-text-primary-dark">
                                {stats.responses}
                            </div>
                            <div className="text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                                {stats.responseBytes} K
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    );
};
