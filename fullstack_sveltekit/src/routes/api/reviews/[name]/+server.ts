import fs from 'fs';
import {json} from "@sveltejs/kit";
import {db} from "$lib/services/db.ts";
import {withAuth} from "$lib/services/apiGuard";

/** @type {import('./$types').RequestHandler} */
export async function GET({params}) {
    return json({review: await db.reviews.get(params.name)});
}

export const DELETE = withAuth(async ({params, url}) => {
    console.log(params.name);
    await db.reviews.delete(params.name);
    return json({status: 'ok'});
});


export const PUT = withAuth(async ({url, request, params}) => {
    const payload = await request.json();

    if (payload.action === 'refresh') {
        const MAIN_STORAGE = import.meta.env.VITE_STORAGE_REVIEWS;
        const review = await db.reviews.get(params.name);
        const REVIEW_STORAGE = `${MAIN_STORAGE}/${review.name}`;
        const logFileName = `${REVIEW_STORAGE}/events.log`;
        const picturesList = fs.readdirSync(REVIEW_STORAGE).filter(f => f.endsWith('.jpg') || f.endsWith('.jpeg') || f.endsWith('.png'));

        fetch(`${url.protocol}//${url.host}/api/internal/reviews/${review.name}`, {
            method: 'POST',
            body: JSON.stringify({
                action: 'refresh',
                logFileName: logFileName
            }),
        })

        return json({picturesList});
    }
});