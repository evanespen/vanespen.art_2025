from datetime import datetime
from typing import Optional, List

from pydantic import BaseModel


class Picture(BaseModel):
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


class SpeciesPicture(BaseModel):
    name: str
    scientific_name: str
    threat: str
    info_page: str
    description: str
    pictures: Optional[List[Picture]] = None

    class Config:
        orm_mode = True
        arbitrary_types_allowed = True
