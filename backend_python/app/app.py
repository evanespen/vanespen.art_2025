import jwt
from fastapi import FastAPI, Depends, HTTPException, status
from starlette.middleware.cors import CORSMiddleware
from starlette.requests import Request

from app.api.api import api_router
from app.core import config
from app.db.session import Session
import app.db.init_db as init_db

app = FastAPI(
    title='Khazad-Dum backend',
    description='Khazad-Dum backend',
    version="0.1",
    # openapi_url=None,
    # docs_url=None,
    # redoc_url=None
)

# CORS
origins = []

init_db.init_db(Session())

origins = ['*']

app.add_middleware(
    CORSMiddleware,
    allow_origins=['http://localhost:5173'],#origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/alive", include_in_schema=False)
def api_is_up():
    return {"alive": "yes"}


@app.middleware("http")
async def db_session_middleware(request: Request, call_next):
    request.state.db = Session()
    response = await call_next(request)
    request.state.db.close()
    return response

app.include_router(api_router, prefix='/api')
