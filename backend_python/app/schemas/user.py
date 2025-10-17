from typing import Optional

from pydantic import BaseModel


# Shared properties
class UserBase(BaseModel):
    is_superuser: Optional[bool] = False
    username: Optional[str] = None


class UserBaseInDB(UserBase):
    id: int = None

    class Config:
        orm_mode = True


# Properties to receive via API on creation
class UserCreate(UserBaseInDB):
    username: str
    password: str


# Properties to receive via API on update
class UserUpdate(UserBaseInDB):
    pass


# Additional properties to return via API
class User(UserBaseInDB):
    pass


# Additional properties stored in DB
class UserInDB(UserBaseInDB):
    hashed_password: str
