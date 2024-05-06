from google.oauth2 import service_account as GCPSrvAcc
from config import settings
from typing import Any
import cloudProviders.gcp.storage as Storage

def upload_blob_from_memory(gcp_crds:  GCPSrvAcc.Credentials, content: Any, path: str):
    Storage.upload_blob_from_memory(gcp_crds, settings.GCP_Config.BUCKET_NAME, content, path)

def test():
    pass