import {db} from "$lib/services/db.ts";
import Moment from "moment";

/**
 * @typedef {{
 *   url: string;
 * }} Picture
 */

/** @type {import('./$types').PageServerLoad} */
export const load = async () => {
    let _pictures = await db.pictures.all();

    let pictures = {};
    _pictures.forEach(p => {
        const month = Moment(p.timestamp).format('MMMM YYYY');
        if (!Object.keys(pictures).includes(month)) pictures[month] = [];
        pictures[month].push(p);
    })

    return {
        pictures
    }
};
