from schemas import Response as API_Response, VPC as VPC_Schema
from fastapi import APIRouter, Query
from api_v1.handlers.sdk import vpc as SDK_VPC
from config import settings, constants

router = APIRouter()


@router.get("/vpc/{cloud_provider}/{framework}")
async def get_vpc(cloud_provider: str, framework: str, vpc_name: str = Query(default=None, max_length=200)):
    print("VPC: ", vpc_name, framework, cloud_provider)
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_VPC.list_vpcs(vpc_name=vpc_name)
    except Exception as exp:
        print("VPC GET Exception: ", exp)
        return API_Response(error=exp, statusCode=400).model_dump()


@router.get("/version")
async def get_version():
    print("Version API - Instance Microservice 1234")
    return {
        "version": "1.0.0"
    }


@router.post("/vpc/{cloud_provider}/{framework}")
def create_vpc(cloud_provider: str, framework: str,  details: VPC_Schema):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_VPC.create_vpc_subnetwork_firewall(details=details)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()


@router.delete("/vpc/{cloud_provider}/{framework}")
def delete_vpc(cloud_provider: str, framework: str, vpc_name: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_VPC.delete_vpc(vpc_name=vpc_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
