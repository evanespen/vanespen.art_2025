import {username} from "$lib/services/userStore";
import * as cookie from 'cookie';
import jwt from "jsonwebtoken";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({event, resolve}) {

    if (event.url.pathname === '/logout') {
        event.cookies.set('sessionid', null);
        console.log(event.cookies.delete('sessionid'))
        event.request.headers.set('cookie', event.cookies);
        event.locals.user = null;
        username.set(null);
        console.log('logged out');

        return await resolve(event);
    }

    if (event.url.pathname.includes('admin')) {
        if (event.cookies.get('sessionid') === undefined) {
            return new Response('Redirect', {status: 303, headers: {Location: '/login'}});
        }
    }


    if (event.cookies.get('sessionid') !== undefined) {
        const cookies = cookie.parse(event.request.headers.get('cookie') || '');
        const token = cookies['sessionid'];
        try {
            const decoded = jwt.verify(token, import.meta.env.VITE_SECRET_KEY);
            event.locals.user = decoded.username;
            username.set(event.locals.user);
        } catch (e) {
            return new Response('Redirect', {status: 303, headers: {Location: '/login'}});
        }
    }

    return await resolve(event);
}
