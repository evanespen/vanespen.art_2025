import {db} from "$lib/services/db";

/** @type {import('./$types').PageServerLoad} */
export async function load({params}) {
    return {
        album: await db.albums.get(params.name)
    };
}
