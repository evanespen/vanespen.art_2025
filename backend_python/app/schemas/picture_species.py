from datetime import datetime
from typing import Optional, List

from pydantic import BaseModel, validator, Field


class Species(BaseModel):
    name: str
    scientific_name: str
    threat: str
    info_page: str
    description: str
    id: int


class PictureSpecies(BaseModel):
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
    species: Optional[Species] = None

    @validator('species')
    def check_species(cls, v):
        print('!!!!!!!!!! VALIDATOR SPECIES', cls, v)
        if v is None:
            return dict

    # @validator('species_id')
    # def check_species_id(cls, v):
    #     print('?????????? VALIDATOR SPECIES_ID', cls, v)
    #     return v

    class Config:
        orm_mode = True
        arbitrary_types_allowed = True
