from schemas import Response as API_Response
from fastapi import APIRouter, Query
from api_v1.handlers import artifactRegistry as ArtifactRegistry
from config import settings

router = APIRouter()

@router.get("/artifact_registry/{cloud_provider}")
def get_artifact_registry(cloud_provider: str, registry_name: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return ArtifactRegistry.list_artifact_registries(registry_name=registry_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
    
@router.get("/artifact_registry_permissions/{cloud_provider}")
def get_artifact_registry_permissions(cloud_provider: str, registry_name: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return ArtifactRegistry.get_artifact_registry_permissions(repository_id=registry_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
    