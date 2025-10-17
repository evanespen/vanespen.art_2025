from app import crud
from app.core import config
from app.schemas.user import UserCreate

# make sure all SQL Alchemy models are imported before initializing DB
# otherwise, SQL Alchemy might fail to initialize properly relationships
# for more details: https://github.com/tiangolo/full-stack-fastapi-postgresql/issues/28
from app.db.base_class import Base
from app.db.session import engine


def init_db(db_session):
    # Tables should be created with Alembic migrations
    # But if you don't want to use migrations, create
    # the tables un-commenting the next line
    Base.metadata.create_all(bind=engine)

    user = crud.users.get_by_username(db_session, username=config.SUPERUSER)
    if not user:
        user_in = UserCreate(
            username=config.SUPERUSER,
            password=config.SUPERUSER_PASSWORD,
            is_superuser=True,
        )
        user = crud.users.create(db_session, obj_in=user_in)
