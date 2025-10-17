import fs from "fs";

/** @type {import('./$types').RequestHandler} */
export async function GET({params, url}) {
    const typesDirs = {
        thumb: import.meta.env.VITE_STORAGE_THUMB,
        half: import.meta.env.VITE_STORAGE_HALF,
        full: import.meta.env.VITE_STORAGE_FULL,
    };
    const wantedType = url.searchParams.get('type');

    const filepath = `${typesDirs[wantedType]}/${params.filename}`;
    let stat = fs.statSync(filepath);

    return new Response(
        fs.readFileSync(filepath),
        {
            headers: {
                'Content-Type': 'image/jpeg',
                'Content-Length': stat.size
            }
        }
    );
}
