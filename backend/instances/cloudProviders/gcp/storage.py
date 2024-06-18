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
    blob.upload_from_filename(
        source_file_name, if_generation_match=generation_match_precondition)


def upload_blob_from_memory(gcp_crds: GCPSrvAcc.Credentials, bucket: str, content: Any, destination_file_name: str):
    body = content
    if type(content) == dict:
        body = json.dumps(content).encode('utf-8')
    bucket = get_storage_client(gcp_crds).bucket(bucket)
    blob = bucket.blob(destination_file_name)
    blob.upload_from_string(body)


def get_blob(gcp_crds: GCPSrvAcc.Credentials, bucket: str, blob_name: str):
    bucket = get_storage_client(gcp_crds).bucket(bucket)
    blob = bucket.blob(blob_name)
    return blob


def create_bucket(gcp_crds: GCPSrvAcc.Credentials, input: dict):
    client = get_storage_client(gcp_crds)
    bucket = client.bucket(input.get("bucket_name"))
    bucket.iam_configuration.uniform_bucket_level_access_enabled = input.get(
        "uniform_bucket_level_access")
    bucket.storage_class = input.get("storage_class")
    bucket.versioning_enabled = input.get("versioning")
    bucket = client.create_bucket(bucket)
    lifecycle_rules = input.get("lifecycle_rules", [])
    for rule in lifecycle_rules:
        action = rule["action"]
        condition = rule["condition"]
        if action == "Delete":
            bucket.add_lifecycle_delete_rule(age=condition)
    bucket.patch()
    return bucket

def delete_blob(gcp_crds: GCPSrvAcc.Credentials, bucket: str, path: str):
    bucket = get_storage_client(gcp_crds).bucket(bucket)
    blob = bucket.blob(path)
    blob.delete()
    return blob 