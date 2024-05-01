from google.oauth2 import service_account as GCPSrvAcc
from googleapiclient.discovery import build
from typing import Any
import time
from cloudProviders.gcp.credentials import get_gcp_crds


def get_artifact_registry_client(gcp_crds):
    return build(serviceName="artifactregistry", version="v1", credentials=gcp_crds)


def list_artifact_repositories(gcp_crds: GCPSrvAcc.Credentials,region: str, artifact_registry_client: Any = None):
    if artifact_registry_client is None:
        artifact_registry_client = get_artifact_registry_client(gcp_crds)
    project_id = gcp_crds.project_id
    artifact_registry = artifact_registry_client.projects().locations().repositories().list(parent=f"projects/{project_id}/locations/{region}").execute()
    return artifact_registry

def get_artifact_registry_permissions(gcp_crds: GCPSrvAcc.Credentials, region: str, repository_id: str, artifact_registry_client: Any = None):
    if artifact_registry_client is None:
        artifact_registry_client = get_artifact_registry_client(gcp_crds)
    project_id = gcp_crds.project_id
    artifact_registry = artifact_registry_client.projects().locations().repositories().getIamPolicy(
        resource=f"projects/{project_id}/locations/{region}/repositories/{repository_id}").execute()
    return artifact_registry

def add_service_account_to_policy(gcp_crds: GCPSrvAcc.Credentials, region: str, repository_id: str, service_account_email: str, roles: list[str]):
    artifact_registry_client = get_artifact_registry_client(gcp_crds)
    project_id = gcp_crds.project_id
    
    # Get the current IAM policy
    current_policy = get_artifact_registry_permissions(gcp_crds, region, repository_id, artifact_registry_client)
    
    # Add the new service account to the IAM policy
    for role in roles:
        if 'bindings' not in current_policy:
            current_policy['bindings'] = []
        current_policy['bindings'].append({'role': role, 'members': [f'serviceAccount:{service_account_email}']})
    
    # Update the IAM policy with the new changes
    updated_policy = artifact_registry_client.projects().locations().repositories().setIamPolicy(
        resource=f"projects/{project_id}/locations/{region}/repositories/{repository_id}", 
        body={'policy': current_policy}).execute()
    
    return updated_policy