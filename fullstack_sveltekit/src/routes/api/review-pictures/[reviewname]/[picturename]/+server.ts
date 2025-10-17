import fs from 'fs';
import {json} from '@sveltejs/kit';
import {db} from "$lib/services/db.ts";

/** @type {import('./$types').RequestHandler} */
export async function GET({params, url}) {
    const wantedType = url.searchParams.get('type');
    const reviewsStorage = import.meta.env.VITE_STORAGE_REVIEWS;

    let filePath;
    if (wantedType === null) {
        filePath = `${reviewsStorage}/${params.reviewname}/${params.picturename}`;
    } else if (wantedType === 'half') {
        filePath = `${reviewsStorage}/${params.reviewname}/half/${params.picturename}`;
    }


    let headers;
    if (filePath.endsWith('.zip')) {
        headers = {
            'Content-Type': 'application/zip',
            'Content-Length': fs.statSync(filePath).size,
            'Content-Disposition': `attachment; filename=${params.picturename}`
        }
    } else {
        headers = {
            'Content-Type': 'image/jpeg',
            'Content-Length': fs.statSync(filePath).size
        }
    }

    return new Response(
        fs.readFileSync(filePath),
        {headers}
    );
}

export async function PUT({params, request}) {
    const payload = await request.json();

    if (payload.action === 'setStatus') {
        await db.reviews.picture.setStatus(params.reviewname, params.picturename, payload.value);
    } else if (payload.action === 'setComment') {
        await db.reviews.picture.setComment(params.reviewname, params.picturename, payload.value);
    }

    return json({status: 'ok'});
}