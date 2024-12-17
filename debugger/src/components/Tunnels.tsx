import React from 'react';
import {Tunnel} from '../types';

interface TunnelsProps {
    tunnel?: Tunnel;
}

export const Tunnels: React.FC<TunnelsProps> = ({tunnel}) => {
    return (
        <div className="h-full flex flex-col">
            <div className="px-4 py-2">
                <span className="text-sm font-medium text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                    Tunnels
                </span>
            </div>
            <div className="flex-1 px-4 flex items-center">
                <div className="w-full space-y-3">
                    {['http', 'https'].map((protocol) => (
                        <div key={protocol} className="flex items-center">
                            <div className="flex-grow bg-echogy-bg-secondary dark:bg-echogy-bg-secondary-dark rounded overflow-hidden flex border border-echogy-border dark:border-echogy-border-dark">
                                <input
                                    type="text"
                                    value={`${protocol}://${tunnel?.href}`}
                                    readOnly
                                    className="bg-transparent text-echogy-text-primary dark:text-echogy-text-primary-dark px-3 py-2 flex-grow font-mono text-sm outline-none"
                                />
                                <div className="flex">
                                    <button
                                        className="px-3 py-2 hover:bg-echogy-bg-secondary/50 transition-colors border-l border-echogy-border dark:border-echogy-border-dark">
                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8m-6 4h8m-6 4h8"/>
                                        </svg>
                                    </button>
                                    <button
                                        className="px-3 py-2 hover:bg-echogy-bg-secondary/50 transition-colors border-l border-echogy-border dark:border-echogy-border-dark">
                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                                  d="M12 4v1m6 11h2m-6 0h-2v-4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"/>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
