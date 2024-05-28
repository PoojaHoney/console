from schemas import Response as API_Response, ComputeEngine as ComputeEngine_Schema
from fastapi import APIRouter, Query
from api_v1.handlers import compute as Compute
from config import settings

router = APIRouter()

@router.get("/compute/{cloud_provider}")
def get_compute(cloud_provider: str, compute_engine_name: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return Compute.list_compute_engines(compute_engine_name=compute_engine_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
    
@router.post("/compute/{cloud_provider}")
def create_compute(cloud_provider: str, details: ComputeEngine_Schema):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return Compute.create_compute_engine(details=details)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()