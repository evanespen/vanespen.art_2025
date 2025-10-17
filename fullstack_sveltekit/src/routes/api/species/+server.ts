import {db} from '$lib/services/db.ts';
import {json} from "@sveltejs/kit";

/** @type {import('./$types').RequestHandler} */
export async function GET() {
    return json({species: await db.species.withPictures()});
}

export async function PUT({request}) {
    const payload = await request.json();
    const pictureId = payload.pictureId;
    const specieId = payload.specieId;
    const action = payload.action;

    if (action === 'add') {
        await db.species.addPicture(specieId, pictureId);
    } else if (action === 'remove') {
        await db.species.removePicture(specieId, pictureId);
    }

    return json({status: 'ok'});
}
