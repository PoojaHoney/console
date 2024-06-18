from google.oauth2 import service_account as GCPSrvAcc
from config import settings
from typing import Any
import cloudProviders.gcp.storage as Storage

def upload_blob_from_memory(gcp_crds:  GCPSrvAcc.Credentials, content: Any, path: str):
    Storage.upload_blob_from_memory(gcp_crds, settings.GCP_Config.BUCKET_NAME, content, path)

def get_blob_text(gcp_crds:  GCPSrvAcc.Credentials, path: str, bucket: str):
    return Storage.get_blob(gcp_crds, bucket, path).download_as_text()

def create_bucket(gcp_crds:  GCPSrvAcc.Credentials, bucket_name: str):
    return 

def delete_blob(gcp_crds:  GCPSrvAcc.Credentials, bucket: str, path: str):
    return Storage.delete_blob(gcp_crds, bucket, path)