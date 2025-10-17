from typing import Optional, List

from pydantic import BaseModel


class FiltersList(BaseModel):
    aperture: Optional[List[str]]
    exposure: Optional[List[str]]
    focal: Optional[List[str]]
    iso: Optional[List[str]]
    lens: Optional[List[str]]
