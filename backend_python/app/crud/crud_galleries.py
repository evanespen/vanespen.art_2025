from typing import List

from sqlalchemy.orm import Session

from app.models.gallery import Gallery
from app.schemas.gallery import GalleryCreate, GalleryUpdate
from app.schemas.gallery import Gallery as GallerySchema
from app.schemas.picture import Picture
from app.crud.base import CRUDBase


class CRUDItem(CRUDBase[Gallery, GalleryCreate, GalleryUpdate]):
    def get_all(
            self, db_session: Session
    ) -> List[Gallery]:
        return db_session.query(Gallery).all()

    def get_by_name(
            self, db_session: Session, name: str
    ) -> List[Gallery]:
        return db_session.query(Gallery).filter(Gallery.name == name).first()

    def add_picture(self, db_session: Session, gallery: GallerySchema, picture: Picture):
        gallery.pictures.append(picture)
        db_session.commit()

    def remove_picture(self, db_session: Session, gallery: GallerySchema, picture: Picture) -> GallerySchema:
        gallery.pictures = list(filter(
            lambda p: p.id != picture.id, gallery.pictures
        ))

        db_session.commit()
        return gallery


galleries = CRUDItem(Gallery)
