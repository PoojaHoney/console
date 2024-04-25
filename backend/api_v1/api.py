from fastapi import APIRouter
from api_v1.endpoints import iam

api_router = APIRouter()
api_router.include_router(iam.router, tags=["iam"])