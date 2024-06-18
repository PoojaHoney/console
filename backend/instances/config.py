import os
from pydantic_settings import BaseSettings
from pymongo import MongoClient


class GCP_Settings(BaseSettings):
    CLOUD_PROVIDER: str = "gcp"
    IAM_TYPE: str = os.environ.get("GCP_IAM_TYPE", "service_account")
    DEFAULT_REGION: str = os.environ.get("GCP_DEFAULT_REGION", "asia-south1")
    PROJECT_ID: str = os.environ.get("GCP_PROJECT_ID", "")
    SRV_ACC_PRIVATE_KEY_ID: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY_ID", "")
    SRV_ACC_PRIVATE_KEY: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY", "")
    SRV_ACC_CLIENT_EMAIL: str = os.environ.get("GCP_SRV_ACC_CLIENT_EMAIL", "")
    SRV_ACC_CLIENT_ID: str = os.environ.get("GCP_SRV_ACC_CLIENT_ID", "")
    AUTH_URL: str = os.environ.get("GCP_AUTH_URL", "")
    TOKEN_URL: str = os.environ.get("GCP_TOKEN_URL", "")
    CLIENT_X509_CERT_URL: str = os.environ.get("GCP_CLIENT_X509_CERT_URL", "")
    BUCKET_NAME: str = os.environ.get("GCP_BUCKET_NAME", "alph-a_console")
    DEFAULT_ZONE: str = os.environ.get("GCP_DEFAULT_ZONE", "asia-south1-c")


class Settings(BaseSettings):
    SERVICE_PORT: int = os.environ.get("SERVICE_PORT", 6001)
    SERVICE_NAME: str = os.environ.get("SERVICE_NAME", "Console Instances")
    SERVICE_DOMAIN: str = os.environ.get("SERVICE_DOMAIN", "localhost")
    SERVICE_BASEPATH: str = os.environ.get(
        "SERVICE_BASEPATH", "/api/instances")
    SERVICE_VERSION: str = os.environ.get("SERVICE_VERSION", "v1")
    AWS_CLOUD_PROVIDER: str = "aws"
    GCP_Config: GCP_Settings = GCP_Settings()  # type: ignore
    MONGO_DATABASE: str = os.environ.get("MONGO_DATABASE", "instances")
    USER_COLLECTION: str = os.environ.get("USER_COLLECTION", "users")
    TOKEN_COLLECTION: str = os.environ.get("TOKEN_COLLECTION", "tokens")
    MONGO_HOST: str = os.environ.get("MONGO_HOST", "localhost")
    MONGO_USERNAME: str = os.environ.get("MONGO_USERNAME", "")
    MONGO_PASSWORD: str = os.environ.get("MONGO_PASSWORD", "")
    MONGO_PORT: int = os.environ.get("MONGO_PORT", 27017)
    INSTANCE_COLLECTION: str = os.environ.get(
        "INSTANCE_COLLECTION", "instances")
    INSTANCE_CONFIGURATION_COLLECTION: str = os.environ.get(
        "INSTANCE_CONFIGURATION_COLLECTION", "instance_configurations")
    USER_OTP_COLLECTION: str = os.environ.get(
        "USER_OTP_COLLECTION", "userOTP")
    USER_DATABASE: str = os.environ.get("USER_DATABASE", "users")
    PRODUCT_DATABASE: str = os.environ.get("PRODUCT_DATABASE", "products")
    PRODUCT_MICROSERVICE_COLLECTION: str = os.environ.get(
        "PRODUCT_MICROSERVICE_COLLECTION", "microservices")
    PRODUCT_RESOURCE_COLLECTION: str = os.environ.get(
        "PRODUCT_RESOURCE_COLLECTION", "resources")
    PRODUCT_PLANS_COLLECTION: str = os.environ.get(
        "PRODUCT_PLANS_COLLECTION", "plans")
    PRODUCT_VERSIONS_COLLECTION: str = os.environ.get(
        "PRODUCT_VERSIONS_COLLECTION", "versions")
    PRODUCT_CONFIG_COLLECTION: str = os.environ.get(
        "PRODUCT_CONFIG_COLLECTION", "configurations")
    PRODUCT_COLLECTION: str = os.environ.get(
        "PRODUCT_COLLECTION", "products")
    PRODUCT_FULL_DETAILS: str = os.environ.get(
        "PRODUCT_FULL_DETAILS", "productFullDetails")


settings = Settings()


class Constants():
    TERRAFORM_FRAMEWORK: str = "terraform"
    SDK_FRAMEWORK: str = "sdk"


constants = Constants()


class Databases():
    mongo_client = None

    @ classmethod
    def get_mongo_connection(cls):
        if cls.mongo_client is None:
            if settings.MONGO_USERNAME != "" and settings.MONGO_PASSWORD != "":
                cls.mongo_client = MongoClient(
                    f"mongodb://{settings.MONGO_USERNAME}:{settings.MONGO_PASSWORD}@{settings.MONGO_HOST}:{settings.MONGO_PORT}/")
            else:
                cls.mongo_client = MongoClient(
                    f"mongodb://{settings.MONGO_HOST}:{settings.MONGO_PORT}/")
        return cls.mongo_client

    @ classmethod
    def get_mongo_database(cls, db_name: str):
        return cls.get_mongo_connection()[db_name]

    @classmethod
    def get_mongo_collection(cls, db_name, collection_name: str):
        return cls.get_mongo_database(db_name)[collection_name]


databases = Databases()
