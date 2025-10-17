from typing import Optional, List

from pydantic import BaseModel

from app.schemas.picture import Picture


# Shared properties
class GalleryBase(BaseModel):
    name: str
    description: Optional[str] = None


class GalleryBaseInDB(GalleryBase):
    id: int = None

    class Config:
        orm_mode = True


# Properties to receive via API on creation
class GalleryCreate(GalleryBaseInDB):
    pass


# Properties to receive via API on update
class GalleryUpdate(GalleryBaseInDB):
    pass


# Additional properties to return via API
class Gallery(GalleryBaseInDB):
    id: int
    pictures = []

    class Config:
        orm_mode = True

# Additional properties stored in DB


class GalleryInDB(GalleryBaseInDB):
    id: int

class GalleryAddImages(BaseModel):
    pictures: List[int]
