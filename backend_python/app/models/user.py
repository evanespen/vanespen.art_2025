from sqlalchemy import Boolean, Column, Integer, String

from app.db.base_class import Base


class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, index=True)
    hashed_password = Column(String)
    is_superuser = Column(Boolean(), default=False)
