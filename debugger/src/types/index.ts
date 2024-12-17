export interface Tunnel {
    href: string;
}

export interface Stats {
    requests: number;
    responses: number;
    requestBytes: number;
    responseBytes: number;
    activeConnections: number;
    totalConnections: number;
}

export interface HttpRequest {
    method: string;
    uri: string;
    headers: Record<string, string>;
    time: string;
}

export interface HttpResponse {
    status: number;
    headers: Record<string, string>;
}

export interface HttpEntity {
    id: string;
    request: HttpRequest;
    response: HttpResponse;
    useTime: number;
}
