from google.oauth2 import service_account as GCPSrvAcc
from googleapiclient.discovery import build
from typing import Any
import time
from cloudProviders.gcp.credentials import get_gcp_crds


def get_compute_client(gcp_crds):
    return build(serviceName="compute", version="v1", credentials=gcp_crds)

def list_compute_engines(gcp_crds: GCPSrvAcc.Credentials, zone: str, compute_client: Any = None):
    if compute_client is None:
        compute_client = get_compute_client(gcp_crds)
    project_id = gcp_crds.project_id
    compute = compute_client.instances().list(project=project_id, zone=zone ).execute()
    return compute

def get_compute_engine(gcp_crds: GCPSrvAcc.Credentials, compute_engine_name: str, zone: str, compute_client: Any = None):
    if compute_client is None:
        compute_client = get_compute_client(gcp_crds)
    project_id = gcp_crds.project_id
    try:
        compute = compute_client.instances().get(project=project_id, instance=compute_engine_name, zone=zone).execute()
    except:
        compute = None
    return compute

def create_compute_engine(gcp_crds: GCPSrvAcc.Credentials, details: Any, compute_client: Any = None):
    if compute_client is None:
        compute_client = get_compute_client(gcp_crds)
    project_id = gcp_crds.project_id
    compute = compute_client.instances().insert(project=project_id, zone=details["zone"], body=details).execute()
    return compute

def get_latest_debian_compute_image(gcp_crds: GCPSrvAcc.Credentials, compute_client: Any = None):
    if compute_client is None:
        compute_client = get_compute_client(gcp_crds)
    image = compute_client.images().getFromFamily(project="debian-cloud", family="debian-11").execute()
    return image.get("selfLink")