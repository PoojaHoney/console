from schemas import IAM_Role, Response as API_Response
from fastapi import APIRouter, Query
from api_v1.handlers import iam as IAM

router = APIRouter()

@router.post("/iam_role/{cloud_provider}")
def create_iam_role(cloud_provider: str, details: IAM_Role):
    try:
        return IAM.create_iam_role(cloud_provider=cloud_provider, details=details)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.put("/iam_role/{cloud_provider}/{action}")
def update_iam_role(cloud_provider: str, details: IAM_Role, action: str):
    try:
        return IAM.update_iam_role(cloud_provider=cloud_provider, details=details, action=action)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.delete("/iam_role/{cloud_provider}")
def delete_iam_role(cloud_provider: str, service_account: str):
    try:
        return IAM.delete_service_account(cloud_provider=cloud_provider, service_account=service_account)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()

@router.post("/iam_role_keys")
def create_iam_role_keys(service_account: str):
    try:
        return IAM.create_iam_role_keys(service_account=service_account)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()

@router.get("/iam_role_permissions/{cloud_provider}")
def get_iam_role_permissions(cloud_provider: str, role_name: str):
    try:
        return IAM.get_iam_role_permissions(cloud_provider=cloud_provider, role_name=role_name)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()

@router.post("/iam_predefined_roles/{cloud_provider}")
def create_iam_predefined_roles(cloud_provider: str) -> API_Response:
    try:
        response = IAM.create_iam_predefined_roles(cloud_provider=cloud_provider)
        return response
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()
    
@router.get("/iam_predefined_roles/{cloud_provider}")
def get_iam_predefined_roles(cloud_provider: str,
                filter: str = Query(default=None, max_length=200)):
    try:
        return IAM.get_iam_predefined_roles(cloud_provider=cloud_provider, filter=filter)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()