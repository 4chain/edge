import { useEffect, useState } from 'react';
import {useEventSource} from './hooks/eventSource';
import { HttpEntity, Stats, Tunnel } from './types';
import { Header } from './components/Header';
import { Tunnels } from './components/Tunnels.tsx';
import { Stats as StatsComponent } from './components/Stats';
import { RequestList } from './components/RequestList';
import { RequestDetail } from './components/RequestDetail';
import { mockRequests, mockStats, mockTunnel } from './mock/data';
import {Footer} from "./components/Footer.tsx";

const useMock = import.meta.env.VITE_USE_MOCK === 'true';

function App() {
    const [isDark, setIsDark] = useState(() => {
        const isDark = localStorage.getItem('isDark');
        return isDark ? JSON.parse(isDark) : false;
    });

    const [requests, setRequests] = useState<HttpEntity[]>(useMock ? mockRequests : []);
    const [stats, setStats] = useState<Stats>(useMock ? mockStats : {
        requests: 0,
        responses: 0,
        requestBytes: 0,
        responseBytes: 0,
        activeConnections: 0,
        totalConnections: 0,
    });
    const [tunnel, setTunnel] = useState<Tunnel | undefined>(useMock ? mockTunnel : undefined);
    const [selectedRequest, setSelectedRequest] = useState<HttpEntity | null>(null);
    const [tab, setTab] = useState<'request' | 'response'>('request');

    useEffect(() => {
        localStorage.setItem('isDark', JSON.stringify(isDark));
        if (isDark) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }, [isDark]);

    // 只在非 mock 模式下使用 EventSource
    if (!useMock) {
        useEventSource("/events", {
            onRequest: request => {
                setRequests((old)=> [request, ...old])
            },
            onRequests: reqs => {
                setRequests(reqs)
            },
            onStats: s => {
                setStats({...s})
            },
            onTunnel: t => {
              setTunnel({...t})
            }
        });
    }

    const handleToggleTheme = () => {
        setIsDark(!isDark);
    };

    return (
        <div className="h-screen flex flex-col bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark text-echogy-text-primary dark:text-echogy-text-primary-dark">
            <Header isDark={isDark} onToggleTheme={handleToggleTheme} />

            <main className="flex-1 flex flex-col min-h-0 px-4">
                {/* Tunnel URLs and Stats */}
                <div className="flex gap-4 py-1 flex-none">
                    <div className="flex-1 bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark">
                        <Tunnels tunnel={tunnel} />
                    </div>
                    <div className="flex-1 bg-echogy-bg-primary dark:bg-echogy-bg-primary-dark">
                        <StatsComponent stats={stats} />
                    </div>
                </div>

                {/* Request List and Detail */}
                <div className="flex-1 flex gap-4 min-h-0 pb-4">
                    <div className="flex-1 min-h-0">
                        <RequestList
                            requests={requests}
                            selectedRequest={selectedRequest}
                            onSelectRequest={setSelectedRequest}
                        />
                    </div>
                    <div className="flex-1 min-h-0">
                        <RequestDetail
                            request={selectedRequest}
                            tab={tab}
                            onTabChange={setTab}
                        />
                    </div>
                </div>
            </main>
            <Footer/>
        </div>
    );
}

export default App;
