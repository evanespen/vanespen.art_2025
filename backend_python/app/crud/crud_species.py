from typing import List

from sqlalchemy.orm import Session

from app.models.species import Species
from app.schemas.species import SpeciesCreate, SpeciesUpdate
from app.crud.base import CRUDBase


class CRUDItem(CRUDBase[Species, SpeciesCreate, SpeciesUpdate]):
    def get_by_name(
            self, db_session: Session, name: str
    ) -> List[Species]:
        return db_session.query(Species).filter(Species.name == name).first()


species = CRUDItem(Species)
