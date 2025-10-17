import type {LayoutServerLoad} from './$types';

export const load: LayoutServerLoad = async ({url, locals}) => {
    const {user} = locals;
    return {
        user
    };
};
