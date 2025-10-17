import {username} from '$lib/services/userStore.ts';

export const load = async ({data}) => {
    if (data.user) {
        username.set(data.user);
    }
};
