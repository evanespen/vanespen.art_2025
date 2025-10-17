from typing import List

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from starlette import status
from loguru import logger


from app import crud
from app.api.utils.db import get_db
from app.schemas.species import Species, SpeciesCreate, SpeciesUpdate
from app.schemas.species_picture import SpeciesPicture
from app.models.user import User as DBUser
from app.api.utils.security import get_current_active_superuser


router = APIRouter()


@router.get("/", response_model=List[Species])
def get_species(
        db: Session = Depends(get_db),
):
    db_species = crud.species.get_multi(db)
    db_species = sorted(db_species, key=lambda s: s.name.lower())
    return db_species


# @router.get("/{id}", response_model=Species)
# def get_one_species(
#         *,
#         db: Session = Depends(get_db),
#         id: int
# ):
#     species = crud.species.get(db_session=db, id=id)
#     return species

@router.get("/{name}")
def get_one_specie_by_name(
        *,
        db: Session = Depends(get_db),
        name: str
):
    db_species = crud.species.get_multi(db)
    for specie in db_species:
        if specie.name == name:
            specie.pictures
            return specie


@router.post("/", response_model=Species, status_code=status.HTTP_201_CREATED)
def create_species(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        species_in: SpeciesCreate
):
    """
    Create a new species (god mode)
    """
    species = crud.species.get_by_name(db, name=species_in.name)
    print(species)
    if species:
        raise HTTPException(
            status_code=409,
            detail="A species with this name already exists",
        )
    user = crud.species.create(db, obj_in=species_in)
    return user


@router.put("/{id}", response_model=Species)
def update_product(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
        species_in: SpeciesUpdate,
):
    """
    Update a species
    """
    species = crud.species.get(db_session=db, id=id)
    if not species:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="species not found")

    species = crud.species.update(
        db_session=db, db_obj=species, obj_in=species_in)
    return species


@router.delete("/{id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_species(
        *,
        db: Session = Depends(get_db),
        current_user: DBUser = Depends(get_current_active_superuser),
        id: int,
):
    """
    Delete a species (thanos mode)
    """
    species = crud.species.get(db_session=db, id=id)
    if not species:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="species not found")

    species = crud.species.remove(db_session=db, id=id)
    return species
