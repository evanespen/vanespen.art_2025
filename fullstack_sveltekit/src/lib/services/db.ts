import pkg from 'pg';

import type {AlbumWithPictures, Picture, ReviewPicture, Specie} from "$lib/types";

const {Pool} = pkg;


const client = new Pool({
    user: import.meta.env.VITE_PG_USER,
    password: import.meta.env.VITE_PG_PASS,
    host: import.meta.env.VITE_PG_HOST,
    port: import.meta.env.VITE_PG_PORT,
    database: import.meta.env.VITE_PG_DB,
});

export const db = {
    pictures: {
        all: async () => {
            try {
                const pictures = await client.query('SELECT * FROM pictures ORDER BY timestamp DESC');
                return pictures.rows;
            } catch (err) {
                console.error(err);
                return [];
            }
        },

        insert: async (picture: Picture) => {
            const queryStr = `
                INSERT INTO pictures(aperture, cam_model, exposure, flash, focal, focal_equiv, iso, lens, mode,
                                     timestamp,
                                     path, stared, blured, landscape, notes)
                VALUES ('${picture.aperture}',
                        '${picture.camera}',
                        '${picture.exposure}',
                        '${picture.flash}',
                        '${picture.focal}',
                        '${picture.focal}',
                        '${picture.iso}',
                        '${picture.lens}',
                        '${picture.mode}',
                        timestamp '${picture.dateString}',
                        '${picture.path}',
                        false,
                        false,
                        ${picture.landscape},
                        '${picture.notes}')
            `;
            return client.query(queryStr);
        },

        delete: async (id: number) => {
            try {
                await client.query(`DELETE
                                    FROM pictures
                                    WHERE id = ${id}`)
            } catch (err) {
                console.error(err);
                return;
            }
        },

        star: async (id: number) => {
            try {
                await client.query(`UPDATE pictures
                                    SET stared = true
                                    WHERE id = ${id}`)
            } catch (err) {
                console.error(err);
                return;
            }
        },

        unstar: async (id: number) => {
            try {
                await client.query(`UPDATE pictures
                                    SET stared = false
                                    WHERE id = ${id}`)
            } catch (err) {
                console.error(err);
                return;
            }
        },

        bySpecie: async (specieId: number) => {
            try {
                const pictures = await client.query(`SELECT *
                                                     FROM pictures
                                                     WHERE pictures.species_id = ${specieId}`);
                return pictures.rows;
            } catch (err) {
                console.error(err);
                return [];
            }
        }
    },
    species: {
        all: async () => {
            try {
                const species = await client.query('SELECT * FROM species');
                return species.rows;
            } catch (err) {
                console.error(err);
                return [];
            }
        },

        post: async (specie: Specie) => {
            try {
                const query = await client.query(`INSERT INTO species(name, scientific_name, threat, info_page, description)
                                                  VALUES ('${specie.name}', '${specie.scientific_name}',
                                                          '${specie.threat}', '${specie.info_page}',
                                                          '${specie.description}')`);
                console.log(query);
                return;
            } catch (err) {
                console.error(err);
            }
        },

        withPictures: async () => {
            try {
                const _species = await client.query('SELECT * FROM species');

                let species = [];
                for
                    (const s of _species.rows) {
                    const picturesQuery = await client.query(`
                        SELECT *
                        FROM pictures
                        WHERE pictures.species_id IS NOT NULL
                          AND pictures.species_id = ${s.id}`);
                    s.pictures = picturesQuery.rows;
                    species.push(s);
                }

                return species;
            } catch (err) {
                console.error(err);
                return [];
            }
        },

        get: async (name: string) => {
            try {
                const specieQuery = await client.query(`
                    SELECT *
                    FROM species
                    WHERE name = '${name}'`);
                const specie = specieQuery.rows[0];
                const picturesQuery = await client.query(`
                    SELECT *
                    FROM pictures
                    WHERE pictures.species_id IS NOT NULL
                      AND pictures.species_id = ${specie.id}`);
                const pictures = picturesQuery.rows;
                specie.pictures = pictures;
                return specie;
            } catch (err) {
                console.error(err);
            }
        },

        addPicture: async (specieId: number, pictureId: number) => {
            await client.query(`
                UPDATE
                    pictures
                SET species_id = ${specieId}
                WHERE id = ${pictureId}`);
        },

        removePicture: async (specieId: number, pictureId: number) => {
            await client.query(`
                UPDATE
                    pictures
                SET species_id = NULL
                WHERE id = ${pictureId}`);
        },
    },
    albums: {
        all: async () => {
            try {
                const albumsQuery = await client.query('SELECT * FROM albums');
                const albums = albumsQuery.rows;

                const albumsPicturesJoinQuery = await client.query('SELECT * FROM albums_pictures');
                const albumsPicturesJoin = albumsPicturesJoinQuery.rows;

                const picturesQuery = await client.query('SELECT * FROM pictures');
                const pictures = picturesQuery.rows;

                albums.forEach(album => {
                    const picturesIds = albumsPicturesJoin.filter(apj => apj.gallery_id === album.id).map(apj => apj.picture_id);
                    album.pictures = [];
                    album.pictures = pictures.filter((picture: Picture) => picturesIds.includes(picture.id));
                })

                return albums;
            } catch (err) {
                console.error(err);
                return [];
            }
        },

        get: async (name: string): Promise<AlbumWithPictures> => {
            return client.query(`SELECT *
                                 FROM albums
                                 WHERE name = '${name}'`).then(res => {
                const album = res.rows[0];

                return client.query(`SELECT *
                                     FROM pictures
                                     WHERE pictures.id IN
                                           (SELECT picture_id FROM albums_pictures WHERE gallery_id = ${album.id})`).then(res => {
                    album.pictures = res.rows;
                    return album;
                })
            });
        },

        post: async (name: string, description: string) => {
            const query = await client.query(`
                INSERT
                INTO albums(name, description)
                VALUES ('${name}', '${description}')`);
            console.log(query);
        },

        addPicture: async (albumId: number, pictureId: number) => {
            try {
                const query = await client.query(`
                    INSERT
                    INTO albums_pictures(picture_id, gallery_id)
                    VALUES (${pictureId}, ${albumId})`)
                console.log(query.rows);
            } catch (err) {
                console.error(err);
                return;
            }
        },

        removePicture: async (albumId: number, pictureId: number) => {
            try {
                const query = await client.query(`
                    DELETE
                    FROM albums_pictures
                    WHERE picture_id = ${pictureId}
                      AND gallery_id = ${albumId}`)
                console.log(query.rows);
            } catch (err) {
                console.error(err);
                return;
            }
        },


    },
    reviews: {
        all: async () => {
            try {
                const reviews = await client.query('SELECT * FROM reviews');
                for (const review of reviews.rows) {
                    const picturesQuery = await client.query('SELECT * FROM review_pictures WHERE review_id = ' + review.id);
                    review.pictures = picturesQuery.rows;
                }

                return reviews.rows;
            } catch (err) {
                console.error(err);
                return [];
            }
        },

        get: async (name: string) => {
            try {
                const reviewQuery = await client.query(`
                    SELECT *
                    FROM reviews
                    WHERE name = '${name}'`);
                let review = reviewQuery.rows[0];
                const picturesQuery = await client.query(`SELECT *
                                                          FROM review_pictures
                                                          WHERE review_id = ${review.id}
                                                          ORDER BY name`);
                review.pictures = picturesQuery.rows;
                return review;
            } catch (err) {
                console.error(err);
                return {};
            }
        },

        post: async (name: string, password: string) => {
            try {
                const query = await client.query(`
                    INSERT
                    INTO reviews(name, password)
                    VALUES ('${name}', '${password}')`);
                console.log(query.rows);
            } catch (err) {
                console.error(err);
                return;
            }
        },

        delete: async (name: string) => {
            try {
                const query = await client.query(`
                    DELETE
                    FROM reviews
                    WHERE name = '${name}'`);
                console.log(query.rows);
            } catch (err) {
                console.error(err);
                return;
            }
        },

        picture: {
            insert: async (picture: ReviewPicture) => {
                try {
                    const query = await client.query(`
                        INSERT
                        INTO review_pictures(path, name, hash, review_id, review_name, landscape, status, comment)
                        VALUES ('${picture.path}', '${picture.name}',
                                '${picture.hash}',
                                '${picture.review_id}', '${picture.review_name}',
                                ${picture.landscape}, '${picture.status}',
                                '${picture.comment}')`);
                    console.log(query.rows);
                } catch (err) {
                    console.error(err);
                    return;
                }
            },

            setStatus: async (reviewName: string, pictureName: string, value: number) => {
                try {
                    console.log('IN DB', reviewName, pictureName, value);
                    const query = await client.query(`
                        UPDATE
                            review_pictures
                        SET status = ${value}
                        WHERE review_name = '${reviewName}'
                          AND name = '${pictureName}'`);
                    console.log(query.rows);
                } catch (err) {
                    console.error(err);
                    return;
                }
            },

            setComment: async (reviewName: string, pictureName: string, value: string) => {
                try {
                    console.log('IN DB', reviewName, pictureName, value);
                    const query = await client.query(`
                        UPDATE
                            review_pictures
                        SET comment = '${value}'
                        WHERE review_name = '${reviewName}'
                          AND name = '${pictureName}'`);
                    console.log(query.rows);
                } catch (err) {
                    console.error(err);
                    return;
                }
            },

            get: async (reviewName: string, pictureName: string) => {
                try {
                    const query = await client.query(`
                        SELECT *
                        FROM review_pictures
                        WHERE review_name = '${reviewName}'
                          AND name = '${pictureName}'`);
                    return query.rows[0];
                } catch (err) {
                    console.error(err);
                    return;
                }
            },

            updateHash: async (reviewName: string, pictureName: string, newHash: string) => {
                try {
                    const query = await client.query(`
                        UPDATE
                            review_pictures
                        SET hash = '${newHash}'
                        WHERE review_name = '${reviewName}'
                          AND name = '${pictureName}'`);
                    console.log(query.rows);
                } catch (err) {
                    console.error(err);
                    return;
                }
            }
        }
    }
}
