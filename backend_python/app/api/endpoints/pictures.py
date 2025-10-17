import json
from random import randint
from hashlib import md5
from os import path, remove
from shutil import move
from typing import List

import pendulum
from fastapi import APIRouter, Depends, File, UploadFile, HTTPException
from fastapi.responses import FileResponse
from sqlalchemy.orm import Session
from sqlalchemy.exc import IntegrityError
from starlette import status
from loguru import logger
import PIL

from app import crud
from app.api.utils.db import get_db
from app.schemas.picture import Picture, PictureCreate, PictureBase
from app.schemas.picture_species import PictureSpecies
from app.schemas.filters_list import FiltersList
from app.core import storage_driver, exif_extractor, config
from app.models.user import User as DBUser
from app.api.utils.security import get_current_active_superuser


router = APIRouter()


def __get_file_hash(file_descriptor):
    hash_md5 = md5()
    for chunk in iter(lambda: file_descriptor.read(4096), b""):
        hash_md5.update(chunk)
    return hash_md5.hexdigest()

@router.get("/refresh")
def refresh(
        db: Session = Depends(get_db),
):
    db_pictures = crud.pictures.get_multi(db)
    for p in db_pictures:
        p.day = p.timestamp.strftime('%Y-%m-%d')
        try:
            print(p.species.name)
        except:
            pass

        storage_driver.make_thumbs(p)
    return 'done'



@router.get("/")  # , response_model=List[Picture])
def get_pictures(
        db: Session = Depends(get_db),
):
    db_pictures = crud.pictures.get_multi(db)
    for p in db_pictures:
        print(p.__dict__)
        p.day = p.timestamp.strftime('%Y-%m-%d')


        if path.isfile(path.join(config.HALFRES_DIR, p.path)):
            image = PIL.Image.open(path.join(config.HALFRES_DIR, p.path))
            width, height = image.size
            print(width, height)
            if height > width:
                p.landscape = False
            elif width > height:
                p.landscape = True

        else:
            print('NO FILE')
        
        try:
            print(p.species.name)
        except:
            pass

    # return []

    db.commit()

    tss = []
    images_by_month = {}
    images = sorted(db_pictures, key=lambda p: p.timestamp.timestamp(), reverse=True)
    for p in images:
        tss.append(p.timestamp.timestamp())
        month = pendulum.parse(p.day).format('MMMM YYYY', locale='fr').capitalize()
        if month not in images_by_month.keys():
            images_by_month[month] = []
            
        images_by_month[month].append(p)


    return images_by_month

@router.post("/")  # , response_model=List[Picture])
def get_pictures_w_filters(
        db: Session = Depends(get_db),
        filters_list: FiltersList = None
):
    db_pictures = crud.pictures.get_multi(db)
    sorted(db_pictures, key=lambda p: p.timestamp.timestamp())

    pictures = db_pictures

    for filter_name, filter_values in filters_list.dict().items():
        print(len(pictures))
        if filter_values != []:
            for picture in pictures:
                if picture.__dict__[filter_name] not in filter_values:
                    pictures.remove(picture)

    for p in pictures:
        try:
            print(p.species.name)
        except:
            pass

    return pictures[::-1]


@router.get("/filters")  # , response_model=List[Picture])
def get_filters(
        db: Session = Depends(get_db),
):
    db_pictures = crud.pictures.get_multi(db)
    filters = {
        "aperture": set([]),
        "exposure": set([]),
        "flash": set([]),
        "focal": set([]),
        "iso": set([]),
        "lens": set([])
    }

    for picture in db_pictures:
        for key, value in picture.__dict__.items():
            if key in filters.keys() and value != '':
                filters[key].add(value)

    return filters


@router.get("/favories")  # , response_model=List[Picture])
def get_pictures_stared(
        db: Session = Depends(get_db),
):
    db_pictures = crud.pictures.get_stared(db)
    sorted(db_pictures, key=lambda p: p.timestamp.timestamp())
    return db_pictures[::-1]


@router.get("/random-landscape")  # , response_model=List[Picture])
def get_random_landscape_pictures(
        db: Session = Depends(get_db),
):
    db_pictures = crud.pictures.get_stared(db)
    picture = db_pictures[randint(0, len(db_pictures) - 1)]
    return FileResponse(path.join(config.HALFRES_DIR, picture.path))


@router.post("/upload")
@router.post("/upload/species-id/{species_id}")
def post_pictures(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        files: List[UploadFile] = File(...),
        species_id: int = None,
):
    species = None
    if species_id is not None:
        species = crud.species.get(db, id=species_id)

        if species is None:
            raise HTTPException(
                status_code=404,
                detail="This species doenst exists",
            )

    for _file in files:
        file_hash = __get_file_hash(_file.file)
        _file.file.seek(0)
        image_filename = file_hash + '.' + _file.filename.split('.')[-1]
        temp_filepath = path.join(config.TEMP_DIR, image_filename)
        with open(temp_filepath, 'wb') as temp_file:
            temp_file.write(_file.file.read())

        final_filepath = path.join(config.PICTURES_DIR, image_filename)
        move(temp_filepath, final_filepath)

        picture_in = exif_extractor.get_exif(final_filepath)
        picture_in['path'] = image_filename

        try:
            picture = crud.pictures.create(db, obj_in=picture_in)
            if species is not None:
                picture.species = species
                logger.info('SPCIES')
                picture = crud.pictures.update_wout_data(db, db_obj=picture)

            storage_driver.make_thumbs(picture)

        except IntegrityError:
            pass

    return 'ok'


@router.get("/stats")
def get_stats(
        *,
        db: Session = Depends(get_db)
):
    db_pictures = crud.pictures.get_multi(db)
    stats = {
        'lens': {},
        'flash': {},
        'exposure': {},
        'aperture': {},
        'mode': {},
        'iso': {},
        'focal': {}
    }
    stats_dates = {}

    for p in db_pictures:
        pdate = p.timestamp.strftime('%Y-%m-%d')
        if pdate in stats_dates.keys():
            stats_dates[pdate] += 1
        else:
            stats_dates[pdate] = 1

        for k in stats.keys():
            v = p.__dict__[k]
            if v != '':
                if v in stats[k].keys():
                    stats[k][v] += 1
                else:
                    stats[k][v] = 1
    return {
        'stats': stats,
        'dates': stats_dates
    }


@router.delete("/{id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_picture(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
):
    """
    Delete a species (thanos mode)
    """
    picture = crud.pictures.get(db_session=db, id=id)
    if not picture:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="picture not found")

    picture = crud.pictures.remove(db_session=db, picture=picture)

    for d in [
            config.PICTURES_DIR,
            config.THUMBS_DIR,
            config.HALFRES_DIR,
            config.TEMP_DIR
    ]:
        if path.isfile(path.join(d, picture.path)):
            remove(path.join(d, picture.path))

    return picture


@router.post("/{id}/togglestar")
def toggle_star_picture(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
):
    picture = crud.pictures.get(db_session=db, id=id)
    if not picture:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="picture not found")

    return crud.pictures.toggle_star(db_session=db, picture=picture)


@router.post("/{id}/toggleblur")
def toggle_blur_picture(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
):
    picture = crud.pictures.get(db_session=db, id=id)
    if not picture:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="picture not found")

    return crud.pictures.toggle_blur(db_session=db, picture=picture)


@router.get("/{id}")
@router.get("/{id}/{data_type}")
def get_picture(
        *,
        db: Session = Depends(get_db),
        id: int,
        data_type: str = None,
):
    picture = crud.pictures.get(db_session=db, id=id)
    if not picture:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="picture not found")

    try:
        print(picture.species.name)
    except:
        pass

    if data_type == 'meta' or data_type is None:
        return picture

    elif data_type == 'full_file':
        return FileResponse(path.join(config.PICTURES_DIR, picture.path))

    elif data_type == 'thumb_file':
        return FileResponse(path.join(config.THUMBS_DIR, picture.path))

    elif data_type == 'halfres_file':
        return FileResponse(path.join(config.HALFRES_DIR, picture.path))

    else:
        raise HTTPException(
            status_code=400, detail='We dont know what you want...')
