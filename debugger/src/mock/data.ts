import { HttpEntity, Stats, Tunnel } from '../types';

export const mockTunnel: Tunnel = {
    href: 'demo.echogy.dev'
};

export const mockStats: Stats = {
    requests: 128,
    responses: 125,
    requestBytes: 25600,
    responseBytes: 51200,
    activeConnections: 3,
    totalConnections: 15
};

export const mockRequests: HttpEntity[] = [
    {
        id: '1',
        request: {
            method: 'GET',
            uri: '/api/users',
            headers: {
                'Accept': 'application/json',
                'User-Agent': 'Mozilla/5.0'
            },
            time: new Date(Date.now() - 5000).toISOString()
        },
        response: {
            status: 200,
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': '1024'
            }
        },
        useTime: 150
    },
    {
        id: '2',
        request: {
            method: 'POST',
            uri: '/api/orders',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer token123'
            },
            time: new Date(Date.now() - 10000).toISOString()
        },
        response: {
            status: 201,
            headers: {
                'Content-Type': 'application/json',
                'Location': '/api/orders/123'
            }
        },
        useTime: 300
    },
    {
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },
    {
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },
    {
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },
    {
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },
    {
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    },{
        id: '3',
        request: {
            method: 'PUT',
            uri: '/api/products/456',
            headers: {
                'Content-Type': 'application/json',
                'If-Match': 'etag123'
            },
            time: new Date(Date.now() - 15000).toISOString()
        },
        response: {
            status: 204,
            headers: {
                'ETag': 'etag124'
            }
        },
        useTime: 200
    }

];
