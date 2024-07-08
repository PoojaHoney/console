import os
from pydantic_settings import BaseSettings
# from pymongo import MongoClient


class GCP_Settings(BaseSettings):
    CLOUD_PROVIDER: str = "gcp"
    IAM_TYPE: str = os.environ.get("GCP_IAM_TYPE", "service_account")
    DEFAULT_REGION: str = os.environ.get("GCP_DEFAULT_REGION", "asia-south1")
    PROJECT_ID: str = os.environ.get("GCP_PROJECT_ID", "console-424205")
    SRV_ACC_PRIVATE_KEY_ID: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY_ID", "84c6cd36101a188c42648801b5cd476794d51613")
    SRV_ACC_PRIVATE_KEY: str = os.environ.get("GCP_SRV_ACC_PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCu0CXKguMu7zQ9\nWy22f6Dn7Xni8ABd+Tw1A6iUN0eTZFYdKqyG1uU36LUb5VVcAkqZ9XOgTb/gATSC\ndAIEQzmaT3Mcs/RRQgx7Qr8OSjDvy3nk8t1wbovWk9W9dLtL1pYMG/uvtySMK2yj\ne85NtTlaKDQ8KIrT0wCsPQT/IlhSY8gpvoX0axwIdX4032PCOxJk05Oan+ywqFu2\nYQi/DfUjo5JrkzHZkkTkOKBoe++vAc1zad4V9OUo8S1UcgC/HjCjsXiALt1nERal\nEtsuKHU9r2jHkD5C6ImG5GIJSfPWP8DLvu/nAezgXxdGytx5fth9IMFVGHZvPpHU\nywxAPPgbAgMBAAECggEAIOm5o8Y+5dBNopuvIKQo9GHGyA5t00OgrU6SB83coQiM\nGhO1xfFh2MPsjWMBTkB0MUaATZc14AP6EfGri4foqX/56VaMUlxAI1juxc78JbnG\nqm97d7aoh27wGCqJNP8i4wbFvVCbBfULEEC4F/Nd5/aX1xe/A6C/iK44aZzayjAF\njyAcf9mDWLi38ZAFY60OuHR5BeK+Gmna0tJHX5OC3Xn8TONryWcBHMO6T/Hwjt4P\n5TbhC/H27GxyTdHBxfF6DoQTAiG19wAN90kQzF9c3AZ5EM0Bqsl1j7qFtuPntr3H\nfpb3WTkwjTigUaXwFqULwOK8Ua1HaekkQsjnkhXzgQKBgQDVX09mu0m21qVbGAdL\nUw5Ou2eRJkXn1TUjOulxB9r3ImL15X26ZFl8sDREtbJH6iH4Ei374eS4iPJMnXCG\n4UPkcssPFIcZnlgE9kMAFKM6N7ve5CEeMPeC2qIRhX4GFUbuQP27Th+7B6CkDcs4\nNrCnrRuGEQuO2lAtQRsV00KEQQKBgQDRvMS7mIajPqIrEuLdP6zZdN18q0DvkTz2\nM54Nz9QbIl4460DLJ+a3UdrXZm7FpifU0FxY/CDg6Qn2Y+Qgjlgb2aJu2h6Ptuvh\nNdWPQm484egQQzOnbeP9VS5jUvHc5MV9KNbIqoyOCe6Si5aFWzONqGdh43r8zV92\nMVziIpa1WwKBgQDISL/Scaj+7Dqo6EQoi1zZwIW1tNRXikVnvWvUwofiBLAZmwJK\njNmG9isJ9/fSE08xeXcagW7dXQb7rakropXFAY/jpuQtwvliJvq7P6P8CMbbsQ35\nNgOPG8SXZ9Pkx+Id0HmbcPjN6Zn7hEsmiAeITTCRxrwvqnMqF3viH7zHQQKBgEzc\ngbhB2C1g2asxW6Q7Hov+cOCsHbrtncqX9fcXz6I2AROY2wDYWYua2rkPV87k6jnq\nNb2CgEYMANOjncl0gMOMCSPYmVSPc3fv85hxftae7x1ukzJrah6/paVOk0YhLGKG\noded8K5o414e+4VZe9YQ/fn4lGyqkq5Mvr48cMB7AoGAE2NSJaqQTilQWQA4psLb\naEoLHdo8CBrb6whXpTTgYu8LayRBuUNRqiOxJq0eFODkLnFTG6MxR86+Pqgs27RE\ndjEmLoGtHNT78Xqnud4OI0V8CbuWAvKaHsWOERnJR7E7IAXxiSDPPzgMbydVEAup\n+8AFcRV00F1eu1q3kt+m03w=\n-----END PRIVATE KEY-----\n")
    SRV_ACC_CLIENT_EMAIL: str = os.environ.get("GCP_SRV_ACC_CLIENT_EMAIL", "console-automation@console-424205.iam.gserviceaccount.com")
    SRV_ACC_CLIENT_ID: str = os.environ.get("GCP_SRV_ACC_CLIENT_ID", "107868089407702649699")
    AUTH_URL: str = os.environ.get("GCP_AUTH_URL", "https://accounts.google.com/o/oauth2/auth")
    TOKEN_URL: str = os.environ.get("GCP_TOKEN_URL", "https://oauth2.googleapis.com/token")
    CLIENT_X509_CERT_URL: str = os.environ.get("GCP_CLIENT_X509_CERT_URL", "https://www.googleapis.com/robot/v1/metadata/x509/console-automation%40console-424205.iam.gserviceaccount.com")
    BUCKET_NAME: str = os.environ.get("GCP_BUCKET_NAME", "inno-console")
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


# class Databases():
#     mongo_client = None

#     @ classmethod
#     def get_mongo_connection(cls):
#         if cls.mongo_client is None:
#             if settings.MONGO_USERNAME != "" and settings.MONGO_PASSWORD != "":
#                 cls.mongo_client = MongoClient(
#                     f"mongodb://{settings.MONGO_USERNAME}:{settings.MONGO_PASSWORD}@{settings.MONGO_HOST}:{settings.MONGO_PORT}/")
#             else:
#                 cls.mongo_client = MongoClient(
#                     f"mongodb://{settings.MONGO_HOST}:{settings.MONGO_PORT}/")
#         return cls.mongo_client

#     @ classmethod
#     def get_mongo_database(cls, db_name: str):
#         return cls.get_mongo_connection()[db_name]

#     @classmethod
#     def get_mongo_collection(cls, db_name, collection_name: str):
#         return cls.get_mongo_database(db_name)[collection_name]


# databases = Databases()
