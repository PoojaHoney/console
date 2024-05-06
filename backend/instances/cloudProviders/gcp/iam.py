from google.oauth2 import service_account as GCPSrvAcc
from googleapiclient.discovery import build
from typing import Any


def get_iam_client(gcp_crds):
    return build(serviceName="iam", version="v1", credentials=gcp_crds)

def get_resource_manager(gcp_crds: GCPSrvAcc.Credentials):
    return build("cloudresourcemanager", "v1", credentials=gcp_crds)

def create_service_account(gcp_crds: GCPSrvAcc.Credentials, details: dict, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    project_id = gcp_crds.project_id
    service_account = iam_client.projects().serviceAccounts().create(
        # name="projects/{}/serviceAccounts/".format(gcp_crds.project_id), 
        name=f"projects/{project_id}",
        body={"accountId": details["name"], "serviceAccount": {
                        "displayName": details["name"],
                        "description": details["description"]}}
    ).execute()
    return service_account

def rename_service_account(gcp_crds: GCPSrvAcc.Credentials, details: dict, iam: Any, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    iam["displayName"] = details["name"]
    iam["description"] = details["description"]
    service_account = iam_client.projects().serviceAccounts().update(
        name=f"projects/-/serviceAccounts/{iam['email']}",
        body=iam
    ).execute()
    return service_account

def disable_service_account(gcp_crds: GCPSrvAcc.Credentials, service_account: str, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    service_account = iam_client.projects().serviceAccounts().disable(
        name=f"projects/-/serviceAccounts/{service_account}"
    ).execute()

def enable_service_account(gcp_crds: GCPSrvAcc.Credentials, service_account: str, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    service_account = iam_client.projects().serviceAccounts().enable(
        name=f"projects/-/serviceAccounts/{service_account}"
    ).execute()

def delete_service_account(gcp_crds: GCPSrvAcc.Credentials, service_account: str, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    iam_client.projects().serviceAccounts().delete(
        name=f"projects/-/serviceAccounts/{service_account}"
    ).execute()

def list_service_accounts(gcp_crds: GCPSrvAcc.Credentials, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    project_id = gcp_crds.project_id
    service_accounts = iam_client.projects().serviceAccounts().list(
        name=f"projects/{project_id}"
    ).execute()
    return service_accounts

def get_service_account(gcp_crds: GCPSrvAcc.Credentials, details: dict, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds)
    email = ""
    if not details.get("email") or details["email"] == "":
        email = f"{details['name']}@{gcp_crds.project_id}.iam.gserviceaccount.com"
    else:
        email = details["email"]
        
    try:
        account = iam_client.projects().serviceAccounts().get(
            name=f"projects/{gcp_crds.project_id}/serviceAccounts/{email}"
        ).execute()
        return account
    except:
        return None

def check_service_account_exists(gcp_crds: GCPSrvAcc.Credentials, details: dict, iam_client: Any = None):
    service_account = get_service_account(gcp_crds=gcp_crds, details=details, iam_client =iam_client)
    if service_account is None:
        return False
    return True

def get_service_account_policy(gcp_crds: GCPSrvAcc.Credentials, service_account: str, version: int = 1):
    resource_manager = get_resource_manager(gcp_crds=gcp_crds)
    policy = list_policies_project(gcp_crds=gcp_crds, resource_manager=resource_manager, version=version)
    polices = []
    for role in policy["bindings"]:
        for member in role["members"]:
            if member == f"serviceAccount:{service_account}":
                polices.append(role["role"])

    return polices

def list_predefined_roles(gcp_crds: GCPSrvAcc.Credentials):
    iam = get_iam_client(gcp_crds=gcp_crds)
    roles = []
    nextToken = ""
    while True:
        if nextToken == "":
            response = iam.roles().list().execute()
        else:
            response = iam.roles().list(pageToken=nextToken).execute()
        nextToken = response.get("nextPageToken")
        roles.extend(response.get("roles", []))
        if nextToken == "" or nextToken is None:
            break
    return roles


def list_policies_project(gcp_crds: GCPSrvAcc.Credentials, version: int = 3, resource_manager: Any = None):
    if resource_manager is None:
        resource_manager = get_resource_manager(gcp_crds=gcp_crds)
    policy = resource_manager.projects().getIamPolicy(
        resource = gcp_crds.project_id,
        body = {"options": {"requestedPolicyVersion": version}}
    ).execute()
    return policy

def set_service_account_policy(gcp_crds: GCPSrvAcc.Credentials, service_account: str, roles: list, resource_manager: Any = None):
    if not len(roles):
        return
    if resource_manager is None:
        resource_manager = get_resource_manager(gcp_crds=gcp_crds)
    policies = list_policies_project(gcp_crds=gcp_crds)
    if policies.get("bindings") is None:
        policies["bindings"] = []
    for role in roles:
        policies["bindings"].append({
                        'role': role,
                        'members': [f'serviceAccount:{service_account}']
                    })
    response = resource_manager.projects().setIamPolicy(
        resource=gcp_crds.project_id,
        body={
            'policy': policies
        }
    ).execute()
    return response

def delete_service_account_policy(gcp_crds: GCPSrvAcc.Credentials, service_account: str, roles: list, resource_manager: Any = None):
    if not len(roles):
        return
    if resource_manager is None:
        resource_manager = get_resource_manager(gcp_crds=gcp_crds)
    policies = list_policies_project(gcp_crds=gcp_crds)
    for role in roles:
        for binding in policies["bindings"]:
            if binding["role"] == role:
                for member in binding["members"]:
                    if member == f"serviceAccount:{service_account}":
                        binding["members"].remove(member)
    response = resource_manager.projects().setIamPolicy(
        resource=gcp_crds.project_id,
        body={
            'policy': policies
        }
    ).execute()
    return response

def update_service_account_policy(gcp_crds: GCPSrvAcc.Credentials, service_account: str, deleted_roles: list[str], added_roles: list[str], resource_manager: Any = None):
    if not len(deleted_roles) and not len(added_roles):
        return
    if resource_manager is None:
        resource_manager = get_resource_manager(gcp_crds=gcp_crds)
    policies = list_policies_project(gcp_crds=gcp_crds)
    if policies.get("bindings") is None:
        policies["bindings"] = []
    if len(added_roles):
        count = len(added_roles)
        for role in added_roles:
            for binding in policies["bindings"]:
                if binding["role"] == role:
                    binding["members"].append(f"serviceAccount:{service_account}")
                    count -= 1
                if count == 0:
                    break
    if len(deleted_roles):
        count = len(deleted_roles)
        for role in deleted_roles:
            for binding in policies["bindings"]:
                if binding["role"] == role:
                    for member in binding["members"]:
                        if member == f"serviceAccount:{service_account}":
                            binding["members"].remove(member)
                            count -= 1
                if count == 0:
                    break
    response = resource_manager.projects().setIamPolicy(
        resource=gcp_crds.project_id,
        body={
            'policy': policies
        }
    ).execute()
    return response

def create_service_account_key(gcp_crds: GCPSrvAcc.Credentials, service_account: str, iam_client: Any = None):
    if iam_client is None:
        iam_client = get_iam_client(gcp_crds=gcp_crds)
    key = iam_client.projects().serviceAccounts().keys().create(
        name=f"projects/-/serviceAccounts/{service_account}"
        ).execute()
    return key