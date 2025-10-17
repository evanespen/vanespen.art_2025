from typing import Optional, List

from pydantic import BaseModel


# Shared properties
class SpeciesBase(BaseModel):
    name: str
    scientific_name: str
    threat: str
    info_page: str
    description: Optional[str] = None


class SpeciesBaseInDB(SpeciesBase):
    id: int = None

    class Config:
        orm_mode = True


# Properties to receive via API on creation
class SpeciesCreate(SpeciesBaseInDB):
    pass


# Properties to receive via API on update
class SpeciesUpdate(SpeciesBaseInDB):
    pass


# Additional properties to return via API
class Species(SpeciesBaseInDB):
    pictures = []

    class Config:
        orm_mode = True

# Additional properties stored in DB


class SpeciesInDB(SpeciesBaseInDB):
    id: int


# Species schema without pictures field
class SpeciesWithoutPictures(BaseModel):
    id: int
    name: str
    scientific_name: str
    threat: str
    info_page: str
    description: str
