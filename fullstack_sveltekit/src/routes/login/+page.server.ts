import {invalid, redirect} from '@sveltejs/kit';
import jwt from 'jsonwebtoken';

/** @type {import('./$types').Actions} */
export const actions = {
    default: async ({cookies: cookies, request: request}) => {
        const data = await request.formData();
        const username = data.get('username');
        const password = data.get('password');

        if (username === import.meta.env.VITE_ADMIN_USER && password === import.meta.env.VITE_ADMIN_PASS) {
            const token = jwt.sign({username: username}, import.meta.env.VITE_SECRET_KEY);
            cookies.set('sessionid', token, {
                httpOnly: false,
            });
            throw redirect(303, '/');
        } else {
            return invalid(400, {username, incorrect: true});
        }
    },
};
