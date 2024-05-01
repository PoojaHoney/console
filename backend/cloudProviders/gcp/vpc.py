from google.oauth2 import service_account as GCPSrvAcc
from googleapiclient.discovery import build
from typing import Any
import time
from cloudProviders.gcp.credentials import get_gcp_crds


def get_vpc_client(gcp_crds):
    return build(serviceName="compute", version="v1", credentials=gcp_crds)

def list_vpcs(gcp_crds: GCPSrvAcc.Credentials, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    vpc = vpc_client.networks().list(project=project_id,).execute()
    return vpc

def get_vpc(gcp_crds: GCPSrvAcc.Credentials, vpc_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    try:
        vpc = vpc_client.networks().get(project=project_id, network=vpc_name).execute()
    except:
        vpc = None
    return vpc

def check_vpc_exists(gcp_crds: GCPSrvAcc.Credentials, vpc_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    if get_vpc(gcp_crds, vpc_name, vpc_client) is None:
        return False
    return True

def create_vpc(gcp_crds: GCPSrvAcc.Credentials, details: dict, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    vpc = vpc_client.networks().insert(project=project_id, body=details).execute()
    time.sleep(30)
    return vpc

def create_subnetwork(gcp_crds: GCPSrvAcc.Credentials, details: dict, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)   
    response = vpc_client.subnetworks().insert(project=gcp_crds.project_id, region=details["region"], body=details).execute()
    time.sleep(10)
    return response

def get_subnetworks(gcp_crds: GCPSrvAcc.Credentials,subnetwork_name: str ,vpc_client: Any= None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    try:
        subnetwork = vpc_client.subnetworks().get(project=project_id, subnetwork=subnetwork_name).execute()
    except: 
        subnetwork = None
    return subnetwork

def list_subnetworks(gcp_crds: GCPSrvAcc.Credentials, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    subnetworks = vpc_client.subnetworks().list(project=project_id).execute()
    return subnetworks

def list_vpc_subnetworks(gcp_crds: GCPSrvAcc.Credentials, vpc_target_link: str, region: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    subnetworks = vpc_client.subnetworks().list(project=project_id, filter=f'network="{vpc_target_link}"', region=region).execute()
    return subnetworks

def check_subnetwork_exists(gcp_crds: GCPSrvAcc.Credentials, subnetwork_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    if get_subnetworks(gcp_crds, subnetwork_name, vpc_client) is None:
        return False
    return True

def list_vpc_firewalls(gcp_crds: GCPSrvAcc.Credentials, vpc_target_link: str,vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    firewalls = vpc_client.firewalls().list(project=project_id, filter=f'network="{vpc_target_link}"').execute()
    return firewalls

def list_firewalls(gcp_crds: GCPSrvAcc.Credentials, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    firewalls = vpc_client.firewalls().list(project=project_id).execute()
    return firewalls

def get_firewall(gcp_crds: GCPSrvAcc.Credentials, firewall_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    try:
        firewall = vpc_client.firewalls().get(project=project_id, firewall=firewall_name).execute()
    except:
        firewall = None
    return firewall

def create_firewall(gcp_crds: GCPSrvAcc.Credentials, details: dict, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    project_id = gcp_crds.project_id
    firewall = vpc_client.firewalls().insert(project=project_id, body=details).execute()
    return firewall

def check_firewall_exists(gcp_crds: GCPSrvAcc.Credentials, firewall_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    if get_firewall(gcp_crds, firewall_name, vpc_client) is None:
        return False
    return True

def delete_vpc(gcp_crds: GCPSrvAcc.Credentials, vpc_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    response = vpc_client.networks().delete(project=gcp_crds.project_id, network=f"{vpc_name}").execute()
    return response

def delete_firewall(gcp_crds: GCPSrvAcc.Credentials, firewall_name: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    response = vpc_client.firewalls().delete(project=gcp_crds.project_id, firewall=firewall_name).execute()
    return response

def delete_subnetwork(gcp_crds: GCPSrvAcc.Credentials, subnetwork_name: str, region: str, vpc_client: Any = None):
    if vpc_client is None:
        vpc_client = get_vpc_client(gcp_crds)
    response = vpc_client.subnetworks().delete(project=gcp_crds.project_id, region=region, subnetwork=subnetwork_name).execute()
    return response