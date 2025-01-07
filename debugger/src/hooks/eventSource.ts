import { useEffect, useState } from 'react';
import {HttpEntity, Stats, Tunnel} from '../types';

export interface EventHandlers {
    onRequest?: (request: HttpEntity) => void;
    onRequests?: (requests: HttpEntity[]) => void;
    onStats?: (stats: Stats) => void;
    onTunnel?: (t : Tunnel) => void;
}

export function useEventSource(url: string, handlers: EventHandlers) {
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        const eventSource = new EventSource(url);

        eventSource.onopen = () => {
            setConnected(true);
        };

        eventSource.onerror = (err) => {
            console.error(err)
            setConnected(false);
        };

        eventSource.addEventListener('sync', (e: MessageEvent) => {
            const {name, data} = JSON.parse(e.data);
            if (name == 'update') {
                handlers.onStats?.(data.stats);
                handlers.onRequest?.(data.httpEntity);
            }else if (name == 'all') {
                handlers.onStats?.(data.stats)
                handlers.onRequests?.(data.httpEntities);
                handlers.onTunnel?.({
                    href: data.tunnel,
                })
            }
        });

        return () => {
            eventSource.close();
        };
    }, [url]);

    return connected;
}