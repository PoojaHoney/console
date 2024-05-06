from schemas import IAM_Role, Response as API_Response
from fastapi import Query
from fastapi.encoders import jsonable_encoder
from config import settings, databases
import json, base64
import cloudProviders.gcp.credentials as GCP_Crds
from typing import Any
import cloudProviders.gcp.iam as IAM
import api_v1.handlers.storage as Storage

#--------------------------------------SERVICE ACCOUNT/IAM ROLES----------------------------------------
def create_service_account(details: IAM_Role, gcp_client: GCP_Crds.get_gcp_crds = None) -> API_Response:
    if details.name == "" or details.name == None:
        return API_Response(error="Please provide valid role name", status_code=400).model_dump()
    if len(details.name) < 6 and len(details.name) > 30:
        return API_Response(error="Role name should be between 6 and 30 characters length", status_code=400).model_dump()
    
    if gcp_client is None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                                "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                "token_uri": settings.GCP_Config.TOKEN_URL,
                            })
    
    details_json = jsonable_encoder(details)
    iam_client = IAM.get_iam_client(gcp_client)
    if IAM.check_service_account_exists(gcp_client, details_json, iam_client):
        return API_Response(error="Service account already exist", data=None, status_code=400).model_dump()

    iam = IAM.create_service_account(gcp_client, details_json, iam_client)
    if not iam.get("name"):
        return API_Response(error="IAM Role creation failed", status_code=400).model_dump()
    
    key_response = create_service_account_key(service_account=iam["email"], iam_client=iam_client)
    if key_response.get("error") != "":
        IAM.delete_service_account(gcp_crds=gcp_client, service_account=iam["email"], iam_client=iam_client)
        return key_response
    Storage.upload_blob_from_memory(gcp_crds=gcp_client, content=key_response["data"], path=f"{details.product}/{details.instanceId}/serviceAccounts/keys/{details.name}.json")
    # if len(details_json["policies"]) == 0 or len(details_json["policies"]) == None:
    #     details_json["policies"] = ["roles/storage.admin", "roles/compute.admin"]
    if len(details_json["policies"]):
        IAM.set_service_account_policy(gcp_crds=gcp_client, roles=details_json["policies"], service_account=iam["email"])
    
    return API_Response(message="IAM Role created successfully", data=iam, status_code=200).model_dump()

def update_service_account(details: IAM_Role, action: str) -> API_Response:
    if action =="" or action == None:
        return API_Response(error=ValueError("Please provide valid action"), status_code=400).model_dump()
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                    })
    iam_client = IAM.get_iam_client(gcp_client)
    if action == "update":
        if details.name == "" or details.email == None or details.email == "":
            return API_Response(error=ValueError("Please provide valid role"), status_code=400).model_dump()
        if len(details.name) < 6 and len(details.name) > 30:
            return API_Response(error=ValueError("Role name should be between 6 and 30 characters length"), status_code=400).model_dump()
        existing_iam = {"email": details.email, "name": details.email.split("@")[0]}
        existing_iam = IAM.get_service_account(gcp_crds=gcp_client, details=jsonable_encoder(existing_iam), iam_client=iam_client)
        if not existing_iam.get("name"):
            return API_Response(error="IAM Role not found", status_code=400).model_dump()
        if existing_iam.get("displayName") != details.name or existing_iam.get("description") != details.description:
            IAM.rename_service_account(gcp_crds=gcp_client, details=jsonable_encoder(details), iam_client=iam_client, iam=existing_iam)
        if len(details.policies):
            roles = IAM.get_service_account_policy(gcp_crds=gcp_client, service_account=details.email)
            new_roles = []
            deleted_roles = []
            for role in details.policies:
                if role not in roles:
                    new_roles.append(role)
            for role in roles:
                if role not in details.policies:
                    deleted_roles.append(role)
            IAM.update_service_account_policy(gcp_crds=gcp_client, service_account=details.email, added_roles=new_roles, deleted_roles=deleted_roles)
    elif action == "enable":
        IAM.enable_service_account(gcp_crds=gcp_client, service_account=details.email, iam_client=iam_client)
    elif action == "disable":
        IAM.disable_service_account(gcp_crds=gcp_client, service_account=details.email, iam_client=iam_client)
    return API_Response(message=f"IAM Role {action} successfully", status_code=200).model_dump()

def list_service_accounts(service_account: str):
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                    })
    iam_client = IAM.get_iam_client(gcp_client)
    if service_account == "" or service_account == None:
        return IAM.list_service_accounts(gcp_crds=gcp_client, iam_client=iam_client)
    else:
        return IAM.get_service_account(gcp_crds=gcp_client, iam_client=iam_client, details={"email":service_account})

def delete_service_account(service_account: str):
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                    })
    iam_client = IAM.get_iam_client(gcp_client)
    if service_account == "" or service_account == None:
        return API_Response(error=ValueError("Please provide service account name"), status_code=400).model_dump()
    IAM.delete_service_account(gcp_crds=gcp_client, service_account=service_account, iam_client=iam_client)
    return API_Response(message="IAM Role deleted successfully", status_code=200).model_dump()

#--------------------------------------SERVICE ACCOUNT KEYS----------------------------------------
def create_service_account_key( service_account: str, iam_client: Any = None):
    if service_account == "" or service_account == None:
        return API_Response(error=ValueError("Please provide service account name"), status_code=400).model_dump()
    gcp_crds = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                        })
    key = IAM.create_service_account_key(gcp_crds=gcp_crds, service_account=service_account, iam_client=iam_client)
    json_key_data = base64.b64decode(key['privateKeyData']).decode('utf-8')
    json_key_file = json.loads(json_key_data)
    #save file to cloud store
    return API_Response(message="Service account key created successfully", data=json_key_file, status_code=200).model_dump()

#--------------------------------------SERVICE ACCOUNT PERMISSIONS/POLICIES/ROLE----------------------------------------
def get_service_account_permissions(role_name: str):
    if role_name == "" or role_name == None:
        return API_Response(error=ValueError("Please provide role name"), status_code=400)
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                        })
    roles = IAM.get_service_account_policy(gcp_client, role_name)
    return API_Response(data=roles, status_code=200, message=f"{role_name} role permissions list").model_dump()

#--------------------------------------IAM ROLES PERMISSIONS PREDEFINED APIs----------------------------------------
def create_iam_predefined_roles():
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                        })
    roles = IAM.list_predefined_roles(gcp_client)
    if len(roles):
        db = databases.get_mongo_database(db_name="gcpMasterData")
        db["iam_predefined_roles"].insert_many(roles)
    return API_Response(data=roles, status_code=200, message="list of predefined roles has been saved to database").model_dump()

def get_iam_predefined_roles(filter: str):
    try:
        roles = []
        db = databases.get_mongo_database(db_name="gcpMasterData")

        if filter == None or filter == "":
            roles = db["iam_predefined_roles"].find()
        else:
            roles = db["iam_predefined_roles"].find({"$or":[{"name": {"$regex": f".*{filter}.*"}},{"title": {"$regex": f".*{filter}.*"}}]})

        return API_Response(data=roles, status_code=200, message="list of predefined roles from JSON file").model_dump()
    except Exception as exp:
        return API_Response(error=exp, status_code=400)
    