import {json} from "@sveltejs/kit";
import {db} from "$lib/services/db.ts";
import fs from "fs";
import crypto from "crypto";
import {resize} from "easyimage";
import ExifReader from 'exifreader';
import JSZip from "jszip";
import {Review} from "$lib/types";

/** @type {import('./$types').RequestHandler} */

function logStatus(fileName: string, step: string, status: boolean, logFileName: string) {
    const payload = {
        fileName,
        step,
        status,
    }
    console.log(payload);
    fs.appendFileSync(logFileName, JSON.stringify(payload) + "\n");
    console.log(logFileName);
}

async function process(fileName: string, review: Review, logFileName: string) {
    console.log('Processing file ' + fileName)
    const MAIN_STORAGE = import.meta.env.VITE_STORAGE_REVIEWS;
    const REVIEW_STORAGE = `${MAIN_STORAGE}/${review.name}`;

    const filePath = `${REVIEW_STORAGE}/${fileName}`;
    const buffer = fs.readFileSync(filePath);
    const HALF_DIR = REVIEW_STORAGE + '/half';
    if (!fs.existsSync(HALF_DIR)) {
        fs.mkdirSync(HALF_DIR);
    }

    const tags = ExifReader.load(buffer, {
        expanded: false,
    });

    const hash = crypto.createHash('md5');
    hash.update(buffer);
    const hex = hash.digest('hex');

    const reviewPicture = {
        path: filePath,
        name: fileName,
        hash: hex,
        review_id: review.id,
        review_name: review.name,
        status: 0,
        comment: '',
        landscape: false,
    }

    if (review.pictures.map(rp => rp.name).includes(reviewPicture.name)) {
        const existingPicture = await db.reviews.picture.get(review.name, reviewPicture.name);
        if (reviewPicture.hash !== existingPicture.hash) {
            console.log(reviewPicture.name, 'already exists but hash differs, updating hash and half');

            try {
                await db.reviews.picture.updateHash(review.name, reviewPicture.name, reviewPicture.hash)
                logStatus(fileName, 'db', true, logFileName);
            } catch (err) {
                logStatus(fileName, 'db', false, logFileName);
                console.error(err);
            }

            try {
                await resize({
                    src: filePath,
                    dst: HALF_DIR + `/${fileName}`,
                    width: tags['Image Width'].value / 2,
                    height: tags['Image Height'].value / 2,
                });
                logStatus(fileName, 'half', true, logFileName);
            } catch (err) {
                console.error(err);
                logStatus(fileName, 'half', false, logFileName);
            }

        } else {
            console.log(reviewPicture.name, 'already exists, same hash, skipping db and half');
            logStatus(fileName, 'db', true, logFileName);
            logStatus(fileName, 'half', true, logFileName);
        }

    } else {
        console.log(reviewPicture.name, 'does not exist, creating in db and half')
        try {
            await db.reviews.picture.insert(reviewPicture);
            logStatus(fileName, 'db', true, logFileName);

            try {
                await resize({
                    src: filePath,
                    dst: HALF_DIR + `/${fileName}`,
                    width: tags['Image Width'].value / 2,
                    height: tags['Image Height'].value / 2,
                });
                logStatus(fileName, 'half', true, logFileName);
            } catch (err) {
                console.error(err);
                logStatus(fileName, 'half', false, logFileName);
            }

        } catch (err) {
            console.error(err);
            logStatus(fileName, 'db', false, logFileName);
        }
    }

}

export async function POST({request, params}) {
    const payload = await request.json();

    if (payload.action === 'refresh') {
        console.log('STARTING INTERNAL PROCESSING')
        const MAIN_STORAGE = import.meta.env.VITE_STORAGE_REVIEWS;
        const REVIEW_STORAGE = `${MAIN_STORAGE}/${params.name}`;

        const review = await db.reviews.get(params.name);
        const logFileName = payload.logFileName;
        const picturesList = fs.readdirSync(REVIEW_STORAGE).filter(f => f.endsWith('.jpg') || f.endsWith('.jpeg') || f.endsWith('.png'));
        console.log(picturesList);

        fs.stat(logFileName, function (err, stats) {
            if (err) return;
            fs.unlinkSync(logFileName);
            console.log('DELETING OLD EVENTS LOG')
        });

        for (const fileName of picturesList) {
            if (fileName.endsWith('.jpg') || fileName.endsWith('.jpeg') || fileName.endsWith('.png')) {
                await process(fileName, review, logFileName);
            }
        }

        // ARCHIVES CREATION
        let maxCount = 50, currentCount = 0;

        let neededArchives = Math.ceil(picturesList.length / maxCount);
        neededArchives = neededArchives == 0 ? 1 : neededArchives;
        let archives = [];
        for (let i = 0; i < neededArchives; i++) {
            let fname = `${review.name}.zip`;
            if (neededArchives > 1) {
                fname = `${review.name}-partie-${i + 1}.zip`;
            }
            console.log('creating archive', fname);
            archives[i] = {
                fname: fname,
                zip: new JSZip(),
                files: picturesList.slice(currentCount, currentCount + maxCount)
            };
            currentCount += maxCount;

            for (const fileName of archives[i].files) {
                archives[i].zip.file(fileName, fs.readFileSync(`${REVIEW_STORAGE}/${fileName}`));
            }

            const zipFname = `${import.meta.env.VITE_STORAGE_REVIEWS}/${review.name}/${archives[i].fname}`;
            archives[i].zip
                .generateNodeStream({type: 'nodebuffer', streamFiles: true})
                .pipe(fs.createWriteStream(zipFname))
                .on('finish', () => {
                    console.log('zipfile created -> ', zipFname);
                    logStatus('archive', 'archive', true, logFileName);
                    logStatus('ALL', 'ALL', true, logFileName);
                })
        }
    }

    return json({status: 'ok'});
}