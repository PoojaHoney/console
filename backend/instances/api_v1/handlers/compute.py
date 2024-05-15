from config import settings
from cloudProviders.gcp import compute as Compute
import cloudProviders.gcp.credentials as GCP_Crds
from fastapi.encoders import jsonable_encoder
from schemas import ComputeEngine as ComputeEngine_Schema, Response as API_Response
from typing import Any

def list_compute_engines(compute_engine_name: str):
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                            "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                            "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                            "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                            "token_uri": settings.GCP_Config.TOKEN_URL,
                        })
    if compute_engine_name == "" or compute_engine_name == None:
        compute_engines = Compute.list_compute_engines(gcp_crds=gcp_client, zone=settings.GCP_Config.DEFAULT_ZONE)
        return API_Response(message="Compute Engines fetched successfully", data=compute_engines, status_code=200).model_dump()
    compute_engine = Compute.get_compute_engine(gcp_crds=gcp_client, compute_engine_name=compute_engine_name, zone=settings.GCP_Config.DEFAULT_ZONE)
    return API_Response(message="Compute Engines fetched successfully", data=compute_engine, status_code=200).model_dump()

def create_compute_engine(details: ComputeEngine_Schema, gcp_client: GCP_Crds.get_gcp_crds = None):
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                                "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                "token_uri": settings.GCP_Config.TOKEN_URL,
                            })
    compute_client = Compute.get_compute_client(gcp_crds=gcp_client)
    if details.disks[0].initializeParams.sourceImage == "" or details.disks[0].initializeParams.sourceImage == None:
        details.disks[0].initializeParams.sourceImage = Compute.get_latest_debian_compute_image(gcp_crds=gcp_client, compute_client=compute_client)
    # details.tags.items.append("http-server")
    # details.tags.items.append("https-server")
    # details.tags.items.append(details.networkInterfaces[0].name)
    details.machineType = f"zones/{details.zone}/machineTypes/{details.machineType}"
    # details.networkInterfaces[0].subnetwork = f"regions/{details.region}/subnetworks/{details.networkInterfaces[0].subnetwork}"
    # details.networkInterfaces[0].name = f"projects/{gcp_client.project_id}/regions/{details.region}/networks/{details.networkInterfaces[0].name}"
    details.disks[0].initializeParams.diskType = f"zones/{details.zone}/diskTypes/{details.disks[0].initializeParams.diskType}" 
    # details.metadata.items.append({"key": "startup-script", "value": startup_script})

    details_json = jsonable_encoder(details)
    compute_engine = Compute.create_compute_engine(gcp_crds=gcp_client, details=details_json)
    return API_Response(message="Compute Engine created successfully", data=compute_engine, status_code=200).model_dump()

