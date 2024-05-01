from schemas import IAM_Role, Response as API_Response
from fastapi import APIRouter, Query
from api_v1.handlers import iam as IAM
from config import settings

router = APIRouter()

#--------------------------------------SERVICE ACCOUNT/IAM ROLES APIs----------------------------------------
@router.post("/iam_role/{cloud_provider}")
def create_iam_role(cloud_provider: str, details: IAM_Role):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.create_service_account(details=details)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.put("/iam_role/{cloud_provider}/{action}")
def update_iam_role(cloud_provider: str, details: IAM_Role, action: str):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.update_service_account( details=details, action=action)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.delete("/iam_role/{cloud_provider}")
def delete_iam_role(cloud_provider: str, service_account: str):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.delete_service_account(service_account=service_account)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.get("/iam_role/{cloud_provider}")
def get_iam_roles(cloud_provider: str, service_account: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.list_service_accounts(service_account=service_account)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()


#--------------------------------------SERVICE ACCOUNT PERMISSIONS/POLICIES/ROLE APIs----------------------------------------

@router.get("/iam_role_permissions/{cloud_provider}")
def get_iam_role_permissions(cloud_provider: str, role_name: str):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.get_service_account_permissions(role_name=role_name)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()

#--------------------------------------SERVICE ACCOUNT KEYS APIs----------------------------------------
@router.post("/iam_role_keys")
def create_iam_role_key(cloud_provider: str, service_account: str):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.create_service_account_key(service_account=service_account)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()

#--------------------------------------IAM ROLES PERMISSIONS PREDEFINED APIs----------------------------------------
@router.post("/iam_predefined_roles/{cloud_provider}")
def create_iam_predefined_roles(cloud_provider: str) -> API_Response:
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            response = IAM.create_iam_predefined_roles()
            return response
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.get("/iam_predefined_roles/{cloud_provider}")
def get_iam_predefined_roles(cloud_provider: str,
                filter: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
            return IAM.get_iam_predefined_roles(filter=filter)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
