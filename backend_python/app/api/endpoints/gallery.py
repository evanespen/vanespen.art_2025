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

from app import crud
from app.api.utils.db import get_db
from app.schemas.picture import Picture, PictureCreate, PictureBase
from app.schemas.picture_species import PictureSpecies
from app.schemas.gallery import Gallery, GalleryCreate, GalleryBase, GalleryAddImages
from app.schemas.filters_list import FiltersList
from app.core import storage_driver, exif_extractor, config
from app.models.user import User as DBUser
from app.api.utils.security import get_current_active_superuser


router = APIRouter()


@router.get("/", response_model=List[Gallery])
def get_galleries(
        db: Session = Depends(get_db),
):
    db_galleries = crud.galleries.get_all(db)
    for g in db_galleries:
        print(g.pictures)
    return db_galleries


@router.get("/{name}", response_model=Gallery)
def get_gallery_by_name(*,
                        db: Session = Depends(get_db),
                        name: str
                        ):
    db_galleries = crud.galleries.get_by_name(db, name=name)
    return db_galleries


@router.post("/", response_model=Gallery, status_code=status.HTTP_201_CREATED)
def create_gallery(
        *,
        db: Session = Depends(get_db),
        # current_user: DBUser = Depends(get_current_active_superuser),
        gallery_in: GalleryCreate
):
    """
    Create a new gallery
    """
    gallery = crud.galleries.get_by_name(db, name=gallery_in.name)
    print(gallery)
    if gallery:
        raise HTTPException(
            status_code=409,
            detail="A gallery with this name already exists",
        )
    gallery = crud.galleries.create(db, obj_in=gallery_in)

    return gallery


@router.post("/{gallery_name}/add", response_model=Gallery)
def add_picture(
        *,
        gallery_name: str,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        add_in: GalleryAddImages
):
    gallery = crud.galleries.get_by_name(db, name=gallery_name)

    print(gallery)
    print(add_in.pictures)

    for id in add_in.pictures:
        picture = crud.pictures.get(db, id)
        crud.galleries.add_picture(db, gallery, picture)

    return gallery


@router.post("/{gallery_name}/remove/{picture_id}", response_model=Gallery)
def remove_picture(
        *,
        gallery_name: str,
        picture_id: int,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
):
    gallery = crud.galleries.get_by_name(db, name=gallery_name)
    picture = crud.pictures.get(db, picture_id)

    gallery = crud.galleries.remove_picture(db, gallery, picture)

    return gallery


@router.delete("/{id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_gallery(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
):
    """
    Delete a gallery
    """
    gallery = crud.galleries.get(db_session=db, id=id)
    if not gallery:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="gallery not found")

    gallery = crud.galleries.remove(db_session=db, id=id)
    return gallery
