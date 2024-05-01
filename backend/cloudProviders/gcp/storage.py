from google.cloud import storage
from google.oauth2 import service_account as GCPSrvAcc
from typing import Any
import json

def get_storage_client(gcp_crds):
    return storage.Client(credentials=gcp_crds)

def upload_blob(gcp_crds: GCPSrvAcc.Credentials, bucket: str, source_file_name: str, destination_file_name: str):
    bucket = get_storage_client(gcp_crds).bucket(bucket)
    blob = bucket.blob(destination_file_name)
    generation_match_precondition = 0
    blob.upload_from_filename(source_file_name, if_generation_match=generation_match_precondition)

def upload_blob_from_memory(gcp_crds: GCPSrvAcc.Credentials, bucket: str, content: Any, destination_file_name: str):
    body = content
    if type(content) == dict:
        body = json.dumps(content).encode('utf-8')
    bucket = get_storage_client(gcp_crds).bucket(bucket)
    blob = bucket.blob(destination_file_name)
    blob.upload_from_string(body)