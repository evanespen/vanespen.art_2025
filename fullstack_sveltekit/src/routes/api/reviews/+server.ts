import {json} from '@sveltejs/kit';
import fs from 'fs';
import {db} from "$lib/services/db.ts";

import {withAuth} from "$lib/services/apiGuard";

/** @type {import('./$types').RequestHandler} */
export const GET = withAuth(async ({request}) => {
    return json({reviews: await db.reviews.all()});
});

export const POST = withAuth(async ({request}) => {
    const payload = await request.json();
    await db.reviews.post(payload.name, payload.password);

    const MAIN_STORAGE = import.meta.env.VITE_STORAGE_REVIEWS;
    const REVIEW_STORAGE = `${MAIN_STORAGE}/${payload.name}`;
    if (!fs.existsSync(REVIEW_STORAGE)) {
        fs.mkdirSync(REVIEW_STORAGE, {recursive: true});
    }

    return json({status: 'ok'});
});

