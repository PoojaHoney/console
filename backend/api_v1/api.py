from fastapi import APIRouter
from api_v1.endpoints import iam, vpc, compute, instances, artifactRegistry

api_router = APIRouter()
api_router.include_router(iam.router, tags=["iam"])
api_router.include_router(vpc.router, tags=["vpc"])
api_router.include_router(compute.router, tags=["compute"])
api_router.include_router(artifactRegistry.router, tags=["artifactRegistry"])
api_router.include_router(instances.router, tags=["instances"])