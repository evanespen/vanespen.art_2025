export const getHeaders = () => {
    for (const cookie of document.cookie.split(';')) {
        if (cookie.includes('sessionid')) {
            const token = cookie.split('sessionid=')[1];
            return {
                'Authorization': 'jwt ' + token,
            }
        }
    }
    return {
        'Authorization': '',
    }
}