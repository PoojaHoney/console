from schemas import IAM_Role, Response as API_Response
from fastapi import APIRouter, Query
from api_v1.handlers.sdk import iam as SDK_IAM
from config import settings, constants

router = APIRouter()

# --------------------------------------SERVICE ACCOUNT/IAM ROLES APIs----------------------------------------


@router.post("/iam_role/{cloud_provider}/{framework}")
def create_iam_role(cloud_provider: str, details: IAM_Role, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.create_service_account(details=details)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


@router.put("/iam_role/{cloud_provider}/{framework}/{action}")
def update_iam_role(cloud_provider: str, details: IAM_Role, action: str, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.update_service_account(details=details, action=action)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


@router.delete("/iam_role/{cloud_provider}/{framework}")
def delete_iam_role(cloud_provider: str, service_account: str, framework: str,
                    product: str = Query(default=""),
                    instanceId: str = Query(default="")):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.delete_service_account(service_account=service_account, instanceId=instanceId, product=product)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


@router.get("/iam_role/{cloud_provider}/{framework}")
def get_iam_roles(cloud_provider: str, framework: str, service_account: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.list_service_accounts(service_account=service_account)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


@router.get("/iam_role_key/{cloud_provider}/{framework}")
def get_iam_role_key(cloud_provider: str, service_account: str, framework: str,
                    product: str = Query(default=""),
                    instanceId: str = Query(default="")):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.get_service_account_key(service_account=service_account, product=product, instanceId=instanceId)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


# --------------------------------------SERVICE ACCOUNT PERMISSIONS/POLICIES/ROLE APIs----------------------------------------

@router.get("/iam_role_permissions/{cloud_provider}/{framework}")
def get_iam_role_permissions(cloud_provider: str, role_name: str, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.get_service_account_permissions(role_name=role_name)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()

# --------------------------------------SERVICE ACCOUNT KEYS APIs----------------------------------------


@router.post("/iam_role_keys")
def create_iam_role_key(cloud_provider: str, service_account: str, framework: str):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.create_service_account_key(service_account=service_account)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()

# --------------------------------------IAM ROLES PERMISSIONS PREDEFINED APIs----------------------------------------


@router.post("/iam_predefined_roles/{cloud_provider}/{framework}")
def create_iam_predefined_roles(cloud_provider: str, framework: str) -> API_Response:
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    response = SDK_IAM.create_iam_predefined_roles()
                    return response
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()


@router.get("/iam_predefined_roles/{cloud_provider}/{framework}")
def get_iam_predefined_roles(cloud_provider: str,
                             framework: str,
                             filter: str = Query(default=None, max_length=200)):
    try:
        if cloud_provider and framework:
            if cloud_provider == settings.GCP_Config.CLOUD_PROVIDER:
                if framework == constants.SDK_FRAMEWORK:
                    return SDK_IAM.get_iam_predefined_roles(filter=filter)
        else:
            return API_Response(error="cloud_provider and framework are required", statusCode=400).model_dump()
    except Exception as exp:
        if len(exp.error_details):
            return API_Response(error=exp.error_details, statusCode=400).model_dump()
        return API_Response(error=exp.args[0], statusCode=400).model_dump()
