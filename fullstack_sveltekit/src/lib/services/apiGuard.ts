import jwt from "jsonwebtoken";

export function withHandlers(...handlers) {
    return async (request) => {
        for (const handle of handlers) {
            const result = await handle(request)
            if (result !== undefined) {
                return result
            }
        }
    }
}

export function authHook(event, request) {
    try {
        const token = event.request.headers.get('Authorization').replace('jwt ', '');
        const decoded = jwt.verify(token, import.meta.env.VITE_SECRET_KEY);
        if (import.meta.env.VITE_ADMIN_USER !== decoded.username) {
            return new Response('Bad Credentials', {status: 401});
        }
    } catch (e) {
        return new Response('Bad Credentials', {status: 401});
    }
}

// create a new handler with auth check
export function withAuth(handle) {
    return withHandlers(authHook, handle);
}