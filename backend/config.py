import os
from pydantic_settings import BaseSettings
from pymongo import MongoClient
from fastapi import Depends

class GCP_Settings(BaseSettings):
    CLOUD_PROVIDER: str = "gcp"
    IAM_TYPE: str = os.environ.get("GCP_IAM_TYPE", "service_account")
    DEFAULT_REGION: str = os.environ.get("GCP_DEFAULT_REGION", "asia-south1")
    PROJECT_ID: str = os.environ.get("GCP_PROJECT_ID")
    SRV_ACC_PRIVATE_KEY_ID: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY_ID")
    SRV_ACC_PRIVATE_KEY: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY")
    SRV_ACC_CLIENT_EMAIL: str = os.environ.get("GCP_SRV_ACC_CLIENT_EMAIL")
    CLIENT_ID: str = os.environ.get("GCP_CLIENT_ID")
    AUTH_URL: str = os.environ.get("GCP_AUTH_URL")
    TOKEN_URL: str = os.environ.get("GCP_TOKEN_URL")
    CLIENT_X509_CERT_URL: str = os.environ.get("GCP_CLIENT_X509_CERT_URL")
    GCP_BUCKET_NAME: str = os.environ.get("GCP_BUCKET_NAME", "alph-a_console")

class Settings(BaseSettings):
    SERVICE_PORT: int = os.environ.get("SERVICE_PORT", 6001)
    SERVICE_NAME: str = os.environ.get("SERVICE_NAME", "Console Instances")
    SERVICE_DOMAIN: str = os.environ.get("SERVICE_DOMAIN", "localhost")
    SERVICE_BASEPATH: str = os.environ.get("SERVICE_BASEPATH", "/api/instances")
    SERVICE_VERSION: str = os.environ.get("SERVICE_VERSION", "v1")
    AWS_CLOUD_PROVIDER: str = "aws"
    GCP_Config: GCP_Settings = GCP_Settings() # type: ignore

class Databases():
    mongo_client = None
    @classmethod
    def get_mongo_connection(cls):
        if cls.mongo_client is None:
            cls.mongo_client = MongoClient("mongodb://localhost:27017/")
        return cls.mongo_client

    @classmethod
    def get_mongo_database(cls, db_name: str):
        return cls.get_mongo_connection()[db_name]

settings = Settings()
databases = Databases()