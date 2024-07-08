from config import settings
import cloudProviders.gcp.credentials as GCP_Crds
from fastapi.encoders import jsonable_encoder
from cloudProviders.gcp import vpc as VPC
from google.oauth2 import service_account as GCPSrvAcc
from schemas import (
    Response as API_Response,
    VPC as VPC_Schema,
    SubNetwork as SubNetwork_Schema,
    FireWall as Firewall_Schema
)
import time
from typing import Any


def list_vpcs(vpc_name: str, gcp_client: GCP_Crds.get_gcp_crds = None):
    if gcp_client == None:
        print("GCP Creds: ",  settings.GCP_Config.PROJECT_ID, settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
              settings.GCP_Config.SRV_ACC_PRIVATE_KEY, settings.GCP_Config.SRV_ACC_CLIENT_EMAIL, settings.GCP_Config.TOKEN_URL)
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    if vpc_name == "" or vpc_name == None:
        vpc = VPC.list_vpcs(gcp_crds=gcp_client)
        return API_Response(message="VPCs fetched successfully", data=vpc, statusCode=200).model_dump()
    vpc = VPC.get_vpc(gcp_crds=gcp_client, vpc_name=vpc_name)
    return API_Response(message="VPCs fetched successfully", data=vpc, statusCode=200).model_dump()


def create_vpc_subnetwork_firewall(details: VPC_Schema, gcp_client: GCP_Crds.get_gcp_crds = None):
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    vpc_client = VPC.get_vpc_client(gcp_crds=gcp_client)
    if not VPC.check_vpc_exists(gcp_crds=gcp_client, vpc_name=details.name, vpc_client=vpc_client):
        vpc = create_vpc(gcp_client=gcp_client,
                         details=details, vpc_client=vpc_client)
        if type(vpc) == API_Response:
            return vpc.model_dump()
        if vpc is None:
            return API_Response(error="VPC creation failed", statusCode=400).model_dump()
    else:
        vpc = VPC.get_vpc(gcp_crds=gcp_client,
                          vpc_name=details.name, vpc_client=vpc_client)
        vpc["targetLink"] = vpc["selfLink"]
    if details.autoCreateSubnetworks != True:
        subnetwork_details = SubNetwork_Schema(
            region=settings.GCP_Config.DEFAULT_REGION,
            name=f"{details.name}-subnet",
            ipCidrRange="10.0.0.0/24",
            vpc_name=details.name,
            network=vpc["targetLink"]
        )
        subnetwork = create_subnetwork(gcp_client=gcp_client, details=jsonable_encoder(
            subnetwork_details), vpc_client=vpc_client)
        if type(subnetwork) == API_Response:
            return subnetwork.model_dump()
        if subnetwork is None:
            return API_Response(error="VPC Created butSubnetwork creation failed", statusCode=400).model_dump()
        vpc["subnetworks"] = [subnetwork]

    firewall_rules = get_sample_firewalls(
        vpc_name=details.name, vpc_target_link=vpc["targetLink"], subnetwork_ip=subnetwork_details.ipCidrRange, ports=details.ports)
    firewalls = create_firewall(
        gcp_client=gcp_client, firewall_rules=firewall_rules, vpc_client=vpc_client)
    # if type(firewalls) == API_Response:
    #     return firewalls.model_dump()
    if firewalls is None or len(firewalls) == 0:
        return API_Response(error="VPC Created but Firewall creation failed", statusCode=400).model_dump()
    vpc["firewall"] = firewalls
    return API_Response(message="VPC created successfully", data=vpc, statusCode=200).model_dump()


def create_firewall(firewall_rules: list[Firewall_Schema], gcp_client: GCPSrvAcc.Credentials = None, vpc_client: Any = None):
    result = []
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    for firewall in firewall_rules:
        if VPC.check_firewall_exists(gcp_crds=gcp_client, firewall_name=firewall["name"], vpc_client=vpc_client):
            continue
        firewall = VPC.create_firewall(
            gcp_crds=gcp_client, details=firewall, vpc_client=vpc_client)
        result.append(firewall)
    return result


def create_subnetwork(details: SubNetwork_Schema, gcp_client: GCPSrvAcc.Credentials = None, vpc_client: Any = None):
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    details_json = jsonable_encoder(details)
    if VPC.check_subnetwork_exists(gcp_crds=gcp_client, subnetwork_name=details_json["name"], vpc_client=vpc_client):
        return API_Response(message="Subnetwork already exists", statusCode=400)
    subnetwork = VPC.create_subnetwork(
        gcp_crds=gcp_client, details=details_json, vpc_client=vpc_client)
    return subnetwork


def create_vpc(details: VPC_Schema, gcp_client: GCPSrvAcc.Credentials = None, vpc_client: Any = None):
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    details_json = {
        "name": details.name,
        "auto_create_subnetworks": details.autoCreateSubnetworks,
        "description": details.description,
        "mtu": details.mtu,
        "routingConfig": jsonable_encoder(details.routingConfig),
    }
    if VPC.check_vpc_exists(gcp_crds=gcp_client, vpc_name=details_json["name"], vpc_client=vpc_client):
        return API_Response(message="VPC already exists", statusCode=400)
    vpc = VPC.create_vpc(gcp_crds=gcp_client,
                         details=details_json, vpc_client=vpc_client)
    return vpc


def delete_vpc(vpc_name: str, gcp_client: GCP_Crds.get_gcp_crds = None):
    try:
        if gcp_client == None:
            gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                                "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                                "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                                "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                                "token_uri": settings.GCP_Config.TOKEN_URL,
                                                })
        vpc_client = VPC.get_vpc_client(gcp_crds=gcp_client)
        vpc = VPC.get_vpc(gcp_crds=gcp_client,
                          vpc_name=vpc_name, vpc_client=vpc_client)
        if type(vpc) == API_Response:
            return vpc.model_dump()
        if not vpc:
            return API_Response(message="VPC not found", statusCode=400).model_dump()
        firewalls = VPC.list_vpc_firewalls(
            gcp_crds=gcp_client, vpc_target_link=vpc["selfLink"], vpc_client=vpc_client)
        if len(firewalls.get("items", [])):
            for firewall in firewalls.get("items", []):
                VPC.delete_firewall(
                    firewall_name=firewall["name"], gcp_crds=gcp_client, vpc_client=vpc_client)
            time.sleep(10)
        subnetworks = VPC.list_vpc_subnetworks(
            gcp_crds=gcp_client, vpc_client=vpc_client, vpc_target_link=vpc["selfLink"], region=settings.GCP_Config.DEFAULT_REGION)
        if len(subnetworks.get("items", [])):
            for subnetwork in subnetworks.get("items", []):
                VPC.delete_subnetwork(subnetwork_name=subnetwork["name"], gcp_crds=gcp_client, vpc_client=vpc_client,
                                      region=settings.GCP_Config.DEFAULT_REGION)
            time.sleep(10)
        time.sleep(30)
        VPC.delete_vpc(vpc_name=vpc_name, gcp_crds=gcp_client)
        time.sleep(30)
        return API_Response(message="VPC deleted successfully", statusCode=200).model_dump()
    except Exception as exp:
        return API_Response(error=str(exp), message="can not delete vpc or vpc does not exists", statusCode=400).model_dump()


def get_sample_firewalls(vpc_name: str, vpc_target_link: str, subnetwork_ip: str, ports: list):
    all_ports = ["22", "80", "443", "3000"] + ports
    return [
        {
            "name": f"{vpc_name}-allow-ssh",
            "direction": "INGRESS",
            "priority": 65534,
            "network": vpc_target_link,
            "sourceRanges": ["0.0.0.0/0"],
            "allowed": [
                {
                    "IPProtocol": "tcp",
                    "ports": ["22"],
                }
            ],
            "description": f"Allow SSH from all on 0.0.0.0/0",
        },
        {
            "name": f"{vpc_name}-allow-rdp",
            "direction": "INGRESS",
            "priority": 65534,
            "network": vpc_target_link,
            "sourceRanges": ["0.0.0.0/0"],
            "allowed": [
                {
                    "IPProtocol": "tcp",
                    "ports": ["3389"],
                }
            ],
            "description": f"Allow RDP from all on 0.0.0.0/0",
        },
        {
            "name": f"{vpc_name}-allow-icmp",
            "direction": "INGRESS",
            "priority": 65534,
            "network": vpc_target_link,
            "sourceRanges": ["0.0.0.0/0"],
            "allowed": [
                {
                    "IPProtocol": "icmp",
                }
            ],
            "description": f"Allow ICMP from all on 0.0.0.0/0",
        },
        {
            "name": f"{vpc_name}-allow-custom",
            "direction": "INGRESS",
            "priority": 65534,
            "network": vpc_target_link,
            "sourceRanges": [subnetwork_ip],
            "allowed": [
                {
                    "IPProtocol": "all",
                }
            ],
            "description": f"Allow ALL from all on 10.0.0.0/24",
        },
        {
            "name": f"{vpc_name}-allow-https",
            "direction": "INGRESS",
            "priority": 1000,
            "network": vpc_target_link,
            "sourceRanges": ["0.0.0.0/0", "35.235.240.0/20"],
            "allowed": [
                {
                    "IPProtocol": "tcp",
                    "ports":  all_ports,
                }
            ],
            "targetTags": ["https-server", vpc_name],
            "description": f"Allow SSH, HTTP, HTTPS, {all_ports}",
        },
        {
            "name": f"{vpc_name}-allow-http",
            "direction": "INGRESS",
            "priority": 1000,
            "network": vpc_target_link,
            "sourceRanges": ["0.0.0.0/0", "35.235.240.0/20"],
            "allowed": [
                {
                    "IPProtocol": "tcp",
                    "ports":  ["80"],
                }
            ],
            "targetTags": ["http-server", vpc_name],
            "description": "Allow HTTP 80 port",
        },
    ]
