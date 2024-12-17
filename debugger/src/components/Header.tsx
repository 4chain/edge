import React from 'react';

interface HeaderProps {
    isDark: boolean;
    onToggleTheme: () => void;
}

export const Header: React.FC<HeaderProps> = ({ isDark, onToggleTheme }) => {
    return (
        <header className="bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark border-b border-echogy-border dark:border-echogy-border-dark">
            <div className="px-6 py-4 flex justify-between items-center">
                <div className="flex items-center space-x-4">
                    <h1 className="text-2xl font-bold text-echogy-text-primary dark:text-echogy-text-primary-dark">
                        <a href="https://echogy.io" target="_blank">Echogy</a>
                    </h1>
                    <span className="text-echogy-text-secondary dark:text-echogy-text-secondary-dark">Web Debugger</span>
                </div>
                <button
                    onClick={onToggleTheme}
                    className="p-2 rounded-full hover:bg-echogy-bg-secondary/50 transition-colors"
                >
                    {isDark ? (
                        <svg className="w-6 h-6 text-echogy-text-primary dark:text-echogy-text-primary-dark"
                             fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/>
                        </svg>
                    ) : (
                        <svg className="w-6 h-6 text-echogy-text-primary dark:text-echogy-text-primary-dark"
                             fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/>
                        </svg>
                    )}
                </button>
            </div>
        </header>
    );
};
