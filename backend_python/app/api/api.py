from fastapi import APIRouter

from app.api.endpoints import login, pictures, species, gallery

api_router = APIRouter()
api_router.include_router(login.router)
api_router.include_router(
    pictures.router, prefix="/pictures", tags=['pictures'])
api_router.include_router(species.router, prefix="/species", tags=['species'])
api_router.include_router(
    gallery.router, prefix='/albums', tags=['galleries'])
