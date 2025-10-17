// @ts-ignore
import {db} from '$lib/services/db.ts';

/**
 * @typedef {{
 *   pictures: []
 * }} Specie
 */


/** @type {import('./$types').PageServerLoad} */
export const load = async () => {
    let species = await db.species.all()
    // @ts-ignore
    species = species.sort((a, b) => {
        if (a.name > b.name) return 1
        else if (a.name < b.name) return -1
        else return 0
    })

    for (const specie of species) {
        specie.pictures = await db.pictures.bySpecie(specie.id)
    }

    return {
        species
    }
};
