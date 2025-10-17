from sqlalchemy import Column, Integer, String, Table, ForeignKey
from sqlalchemy.orm import relationship, backref

from app.db.base_class import Base


association_table = Table('galleries_pictures', Base.metadata,
                          Column('picture_id', Integer, ForeignKey('pictures.id')),
                          Column('gallery_id', Integer, ForeignKey('galleries.id'))
                          )


class Gallery(Base):
    __tablename__ = "galleries"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True)
    description = Column(String, nullable=True)

    pictures = relationship("Picture",
                            secondary=association_table,
                            cascade="all,delete",
                            backref=backref('galleries')
                            )
