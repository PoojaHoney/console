from config import settings
from cloudProviders.gcp import artifactRegistry as ArtifaceRegistry
import cloudProviders.gcp.credentials as GCP_Crds
from schemas import Response as API_Response
from typing import Any

def list_artifact_registries(registry_name: str):
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                        })
    artifact_registries = ArtifaceRegistry.list_artifact_repositories(gcp_crds=gcp_client, region=settings.GCP_Config.DEFAULT_REGION)
    return API_Response(message="Artifact Registries fetched successfully", data=artifact_registries, status_code=200).model_dump()

def get_artifact_registry_permissions(repository_id: str, gcp_client = None):
    if gcp_client is None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                                "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                "token_uri": settings.GCP_Config.TOKEN_URL,
                            })
    artifact_registries = ArtifaceRegistry.get_artifact_registry_permissions(gcp_crds=gcp_client, region=settings.GCP_Config.DEFAULT_REGION, repository_id=repository_id)
    return API_Response(message="Artifact Registries fetched successfully", data=artifact_registries, status_code=200).model_dump()

def set_service_account_artifact_registry(repository_id: str, service_account: str, roles: list[str] ,gcp_client = None):
    if gcp_client is None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                                "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                "token_uri": settings.GCP_Config.TOKEN_URL,
                            })
    bindings = ArtifaceRegistry.add_service_account_to_policy(gcp_crds=gcp_client, region=settings.GCP_Config.DEFAULT_REGION, repository_id=repository_id, service_account_email=service_account, roles=roles)
    return API_Response(message="Artifact Registries fetched successfully", data=bindings, status_code=200).model_dump() 