from datetime import datetime
from typing import Optional, List, Dict

from pydantic import BaseModel, validator

from .species import Species, SpeciesBase, SpeciesWithoutPictures


# Shared properties
class PictureBase(BaseModel):
    aperture: str
    cam_model: str
    exposure: str
    flash: str
    focal: str
    focal_equiv: str
    iso: str
    lens: str
    mode: str
    timestamp: datetime
    path: str
    stared: bool = False
    blured: bool = False


class PictureBaseInDB(PictureBase):
    id: int

    class Config:
        orm_mode = True


# Properties to receive via API on creation
class PictureCreate():
    pass


# Properties to receive via API on update
class PictureUpdate(PictureBaseInDB):
    # species_id: Optional[int] = None
    pass


# Additional properties to return via API
class Picture(PictureBase):
    # species: Optional[SpeciesWithoutPictures] = None
    # species_id: Optional[int] = None
    id: int
    species: Optional[SpeciesWithoutPictures] = None

    @validator('species')
    def validate_species(cls, v):
        if v is None:
            return {}


# Additional properties stored in DB
class PictureInDB(PictureBaseInDB):
    species_id: Optional[int]
