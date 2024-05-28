from schemas import Response as API_Response, VPC as VPC_Schema
from fastapi import APIRouter, Query
from api_v1.handlers import vpc as VPC
from config import settings

router = APIRouter()

@router.get("/vpc/{cloud_provider}")
def get_vpc(cloud_provider: str, vpc_name: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return VPC.list_vpcs(vpc_name=vpc_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()

@router.post("/vpc/{cloud_provider}")
def create_vpc(cloud_provider: str, details: VPC_Schema):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return VPC.create_vpc_subnetwork_firewall(details=details)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
    
@router.delete("/vpc/{cloud_provider}")
def delete_vpc(cloud_provider: str, vpc_name: str):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return VPC.delete_vpc(vpc_name=vpc_name)
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()