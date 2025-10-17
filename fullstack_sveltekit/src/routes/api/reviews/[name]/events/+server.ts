import {json} from "@sveltejs/kit";
import fs from 'fs';
import {db} from "$lib/services/db";
import {withAuth} from "$lib/services/apiGuard";


/** @type {import('./$types').RequestHandler} */
export const GET = withAuth(async ({params}) => {
    const MAIN_STORAGE = import.meta.env.VITE_STORAGE_REVIEWS;
    const review = await db.reviews.get(params.name);
    const REVIEW_STORAGE = `${MAIN_STORAGE}/${review.name}`;
    const logFileName = `${REVIEW_STORAGE}/events.log`;

    let events = [];
    try {
        const rawEvents = fs.readFileSync(logFileName, 'utf-8');
        rawEvents.split('\n').forEach((event) => {
            if (event !== '') events.push(JSON.parse(event));
        })
    } catch (err) {

    }
    return json({events});
});