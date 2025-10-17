from typing import List

from sqlalchemy.orm import Session

from app.models.picture import Picture
from app.models.species import Species
from app.schemas.picture import PictureCreate, PictureUpdate
from app.crud.base import CRUDBase


class CRUDItem(CRUDBase[Picture, PictureCreate, PictureUpdate]):
    def get_by_species(
            self, db_session: Session, species: Species
    ) -> List[Picture]:
        return db_session.query(Picture).filter(Picture.species == species)

    def toggle_star(self, db_session: Session, picture: Picture) -> Picture:
        picture.stared = not picture.stared
        db_session.commit()
        return picture

    def toggle_blur(self, db_session: Session, picture: Picture) -> Picture:
        picture.blured = not picture.blured
        db_session.commit()
        return picture

    def get_stared(self, db_session: Session) -> List[Picture]:
        return db_session.query(Picture).filter(Picture.stared).all()


    def remove(self, db_session: Session, picture: Picture):
        db_session.delete(picture)
        db_session.commit()
        return picture


pictures = CRUDItem(Picture)
