from schemas import (
    Instance, 
    Response as API_Response,
    IAM_Role,
    VPC as VPC_Schema,
    ComputeEngine as ComputeEngine_Schema
)
import random
from config import settings
import cloudProviders.gcp.credentials as GCP_Crds
import string
from api_v1.handlers import iam as IAM, vpc as VPC, compute as Compute, artifactRegistry as ArtifaceRegistry

def create_instance(details: Instance):
    instance_id = random.choice(string.ascii_lowercase) + ''.join(random.choice(string.ascii_lowercase + string.digits) for _ in range(5))
    gcp_client = GCP_Crds.get_gcp_crds({"project_id":settings.GCP_Config.PROJECT_ID,
                                "private_key_id":settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                "private_key":settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                "client_email":settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                "token_uri": settings.GCP_Config.TOKEN_URL,
                            })
    
    iam_response = IAM.create_service_account(IAM_Role(
        name=f"{instance_id}{details.name}",
        instanceId=instance_id,
        product=details.product,
        description=f"{instance_id}{details.name} service account for {details.name} of {details.product} product",
        policies=["roles/artifactregistry.writer", "roles/iam.serviceAccountTokenCreator", "roles/compute.admin", "roles/storage.admin"]
    ), gcp_client=gcp_client)
    if iam_response.get("error") != "":
        return API_Response(error=iam_response.get("error"), status_code=400).model_dump()
    iam = iam_response["data"]

    vpc_response = VPC.create_vpc_subnetwork_firewall(details=VPC_Schema(
        name=f"{instance_id}{details.name}",
        description=f"{instance_id}{details.name} VPC for {details.name} of {details.product} product",
        
    ), gcp_client=gcp_client)
    if vpc_response.get("error") != "":
        return API_Response(error=vpc_response.get("error"), status_code=400).model_dump()
    vpc = vpc_response["data"]
  
    artifact_registry_bindings = ArtifaceRegistry.set_service_account_artifact_registry(
        gcp_client=gcp_client,
        repository_id=details.product,
        service_account=iam["email"],
        roles=["roles/artifactregistry.admin"]
    )

    startup_script = None
    with open("cloudProviders/gcp/masterData/startup_script.sh", "r") as f:
        startup_script = f.read()

    startup_script = startup_script.replace("{{PROJECT_ID}}", gcp_client.project_id)
    startup_script = startup_script.replace("{{CLOUD_STORE_BUCKET}}", settings.GCP_Config.BUCKET_NAME)
    startup_script = startup_script.replace("{{REGION}}", settings.GCP_Config.DEFAULT_REGION)
    startup_script = startup_script.replace("{{SERVICE_ACCOUNT}}", iam["email"])
    startup_script = startup_script.replace("{{INSTANCE_ID}}", instance_id)
    startup_script = startup_script.replace("{{INSTANCE_NAME}}", details.name)

    compute_response = Compute.create_compute_engine(gcp_client=gcp_client, details= ComputeEngine_Schema(
        name=f"{instance_id}{details.name}",
        zone=settings.GCP_Config.DEFAULT_ZONE,
        machineType=details.machineType,
        region=settings.GCP_Config.DEFAULT_REGION,
        networkInterfaces=[ComputeEngine_Schema.NetworkInterfaces(
            name=vpc["targetLink"],
            subnetwork=vpc["subnetworks"][0]["targetLink"],
        )],
        disks=[ComputeEngine_Schema.Disks(
            initializeParams=ComputeEngine_Schema.InitializeParams(
                diskSizeGb=details.diskSizeGb,
                diskType=details.diskType,
                sourceImage=details.sourceImage
            )
        )],
        serviceAccounts=[ComputeEngine_Schema.Scopes(
            email=iam["email"]
        )],
        metadata=ComputeEngine_Schema.Metadata(
            items=[ComputeEngine_Schema.Items(
                key="startup-script",
                value=startup_script
            )]
        ),
        tags=ComputeEngine_Schema.Tags(
            items=["http-server", "https-server", details.product]
        )
    ))
    return API_Response(message="Instance created successfully", data=details, status_code=200).model_dump()