import fs from "fs";
import {db} from "$lib/services/db";


/** @type {import('./$types').PageServerLoad} */
export async function load({cookies, params}) {
    const review = await db.reviews.get(params.name);
    review.password = undefined;

    const reviewPath = import.meta.env.VITE_STORAGE_REVIEWS + '/' + review.name;
    const archives = fs.readdirSync(reviewPath).filter(f => f.endsWith('.zip'));

    return {
        review: review,
        archives: archives,
        authorized: cookies.get(`auth--${review.id}`) || false,
    };
}

/** @type {import('./$types').Actions} */
export const actions = {
    checkAuth: async ({cookies, params, request}) => {
        const data = await request.formData();
        const password = data.get('password');
        const review = await db.reviews.get(params.name);

        if (review.password === password) {
            console.log('password is correct');
            cookies.set(`auth--${review.id}`, true);
            return {auth: true};
        } else {
            console.log('password is wrong');
            return {auth: false};
        }
    }
};
