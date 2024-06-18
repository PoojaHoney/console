from schemas import Response as API_Response, ComputeEngine as ComputeEngine_Schema
from fastapi import APIRouter, Query
from api_v1.handlers.sdk import compute as SDK_Compute
from config import settings, constants

router = APIRouter()


@router.get("/compute/{cloud_provider}/{framework}")
def get_compute(cloud_provider: str, framework: str, compute_engine_name: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_Compute.list_compute_engines(compute_engine_name=compute_engine_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()


@router.post("/compute/{cloud_provider}/{framework}")
def create_compute(cloud_provider: str, details: ComputeEngine_Schema, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_Compute.create_compute_engine(details=details)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()


@router.delete("/compute/{cloud_provider}/{framework}")
def delete_compute(cloud_provider: str, compute_engine: str, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_Compute.delete_compute_engine(compute_engine_name=compute_engine)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
