from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, Boolean
from sqlalchemy.orm import relationship

from app.db.base_class import Base


class Picture(Base):
    __tablename__ = "pictures"
    id = Column(Integer, primary_key=True, index=True)
    aperture = Column(String)
    cam_model = Column(String)
    exposure = Column(String)
    flash = Column(String)
    focal = Column(String)
    focal_equiv = Column(String)
    iso = Column(String)
    lens = Column(String)
    mode = Column(String)
    timestamp = Column(DateTime)
    path = Column(String, unique=True)
    landscape = Column(Boolean, default=True)

    stared = Column(Boolean, default=False)
    blured = Column(Boolean, default=False)

    species_id = Column(Integer, ForeignKey("species.id"))
    species = relationship("Species", back_populates="pictures")

    # galleries = relationship("Gallery", back_populates="pictures")
