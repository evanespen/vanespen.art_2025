from sqlalchemy import Column, Integer, String
from sqlalchemy.orm import relationship

from app.db.base_class import Base


class Species(Base):
    __tablename__ = "species"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True)
    scientific_name = Column(String, unique=True)
    threat = Column(String, nullable=True)
    info_page = Column(String, nullable=True)
    description = Column(String, nullable=True)

    pictures = relationship("Picture", back_populates="species")
