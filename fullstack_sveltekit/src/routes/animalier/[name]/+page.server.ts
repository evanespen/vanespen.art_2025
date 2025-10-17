import {db} from '$lib/services/db.ts';

/**
 * @typedef {{
 *   pictures: [];
 * }} Specie
 */


/** @type {import('./$types').PageServerLoad} */
export async function load({params}) {
    const specie = await db.species.get(params.name);
    return {
        /** @type {Specie[]} */
        specie
    }
}
