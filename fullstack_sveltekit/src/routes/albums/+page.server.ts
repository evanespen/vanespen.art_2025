import {db} from "$lib/services/db.ts";

/**
 * @typedef {{
 *   pictures: []
 * }} Album
 */


/** @type {import('./$types').PageServerLoad} */
export const load = async () => {
    return {
        /** @type {Album[]} */
        albums: await db.albums.all()
    }
};
