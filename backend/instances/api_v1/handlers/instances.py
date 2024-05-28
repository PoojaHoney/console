from schemas import (
    Instance,
    Response as API_Response,
    IAM_Role,
    VPC as VPC_Schema,
    ComputeEngine as ComputeEngine_Schema
)
import time
import random
from config import settings
import cloudProviders.gcp.credentials as GCP_Crds
import string
from api_v1.handlers import iam as IAM, vpc as VPC, compute as Compute, artifactRegistry as ArtifaceRegistry, storage as CloudStorage


def create_instance_compute_engine(details: Instance, product: dict):
    instance_id = random.choice(string.ascii_lowercase) + ''.join(
        random.choice(string.ascii_lowercase + string.digits) for _ in range(5))
    gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                        "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                        "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                        "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                        "token_uri": settings.GCP_Config.TOKEN_URL,
                                        })
    iam = None
    vpc = None
    compute_response = None
    try:
        ports = []
        for microservice in product["microservices"]:
            ports.append(microservice["portNumber"])
        vpc_response = VPC.create_vpc_subnetwork_firewall(details=VPC_Schema(
            name=f"{instance_id}{details.name}",
            instanceId=instance_id,
            ports=ports,
            description=f"{instance_id}{details.name} VPC for {details.name} of {details.productID} product",
        ), gcp_client=gcp_client)
        if vpc_response.get("error") != "":
            return API_Response(error=vpc_response.get("error"), statusCode=400).model_dump()
        vpc = vpc_response["data"]

        policies = None
        for permissions in product["configuration"]["providerPermissions"]:
            if permissions["provider"] == settings.GCP_Config.CLOUD_PROVIDER:
                policies = permissions["permissions"]
                break

        iam_response = IAM.create_service_account(IAM_Role(
            name=f"{instance_id}{details.name}",
            instanceId=instance_id,
            product=details.productID,
            description=f"{instance_id}{details.name} service account for {details.name} of {details.productID} product",
            policies=policies
        ), gcp_client=gcp_client)
        if iam_response.get("error") != "":
            return API_Response(error=iam_response.get("error"), statusCode=400).model_dump()
        iam = iam_response["data"]

        artifact_registry_bindings = ArtifaceRegistry.set_service_account_artifact_registry(
            gcp_client=gcp_client,
            repository_id=details.productID,
            service_account=iam["email"],
            roles=["roles/artifactregistry.admin"]
        )

        if product["configuration"]["startupScriptFilePath"] != "":
            startup_script = CloudStorage.get_blob_text(
                gcp_crds=gcp_client, bucket=settings.GCP_Config.BUCKET_NAME, path=product["configuration"]["startupScriptFilePath"])

            startup_script = startup_script.replace(
                "{{PROJECT_ID}}", gcp_client.project_id)
            startup_script = startup_script.replace(
                "{{CLOUD_STORE_BUCKET}}", settings.GCP_Config.BUCKET_NAME)
            startup_script = startup_script.replace(
                "{{REGION}}", settings.GCP_Config.DEFAULT_REGION)
            startup_script = startup_script.replace(
                "{{SERVICE_ACCOUNT}}", iam["email"])
            startup_script = startup_script.replace(
                "{{INSTANCE_ID}}", instance_id)
            startup_script = startup_script.replace(
                "{{INSTANCE_NAME}}", details.name)

            for microservice in product["microservices"]:
                for resource in product["resources"]:
                    if microservice["name"] == resource["name"] and details.version == resource["productVersion"]:
                        port = "{{"+microservice["name"].upper()+"_PORT}}"
                        startup_script = startup_script.replace(
                            port, str(microservice["portNumber"]))
                        latest_version = next(
                            (resource_version for resource_version in resource["versions"] if resource_version["latest"]), None)
                        image_version = "{{"+microservice["name"].upper(
                        ) + "_IMAGE_VERSION}}"
                        startup_script = startup_script.replace(
                            image_version, latest_version[f"{details.stage}Tag"])
                        service_name = "{{"+microservice["name"].upper(
                        ) + "_NAME}}"
                        startup_script = startup_script.replace(
                            service_name, latest_version["resourceName"])

        compute_response = Compute.create_compute_engine(gcp_client=gcp_client, details=ComputeEngine_Schema(
            name=f"{instance_id}{details.name}",
            zone=settings.GCP_Config.DEFAULT_ZONE,
            machineType=details.machineType,
            region=settings.GCP_Config.DEFAULT_REGION,
            networkInterfaces=[ComputeEngine_Schema.NetworkInterfaces(
                name=vpc["targetLink"],
                subnetwork=vpc["subnetworks"][0]["targetLink"],
                accessConfigs=[ComputeEngine_Schema.AccessConfig(
                    type="ONE_TO_ONE_NAT",
                    name="External NAT"
                )]
            )],
            disks=[ComputeEngine_Schema.Disks(
                initializeParams=ComputeEngine_Schema.InitializeParams(
                    diskSizeGb=details.memory,
                    diskType=details.diskType,
                    sourceImage=details.vmSourceImage
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
                items=["http-server", "https-server", details.productID]
            )
        ))
        print(compute_response)
        return API_Response(message="Instance created successfully", data=details, statusCode=200).model_dump()
    except Exception as exp:
        if "not ready" in exp.get("reason") and "subnetwors" in exp.get("reason"):
            return API_Response(error=str(exp), statusCode=400, message="VPC limit reached and unable to create subnetwork").model_dump()
        if "was not found" in exp.get("reason"):
            return API_Response(error=str(exp), statusCode=400, message="VPC limit crossed").model_dump()
        delete_instance_compute_engine(
            instance_id=instance_id, name=details.name, gcp_client=gcp_client)
        return API_Response(error=str(exp), statusCode=400).model_dump()


def delete_instance_compute_engine(instance_id: str, name: str, gcp_client: GCP_Crds.get_gcp_crds = None):
    if gcp_client == None:
        gcp_client = GCP_Crds.get_gcp_crds({"project_id": settings.GCP_Config.PROJECT_ID,
                                            "private_key_id": settings.GCP_Config.SRV_ACC_PRIVATE_KEY_ID,
                                            "private_key": settings.GCP_Config.SRV_ACC_PRIVATE_KEY,
                                            "client_email": settings.GCP_Config.SRV_ACC_CLIENT_EMAIL,
                                            "token_uri": settings.GCP_Config.TOKEN_URL,
                                            })
    try:
        Compute.delete_compute_engine(
            compute_engine_name=f"{instance_id}{name}", gcp_client=gcp_client)
        IAM.delete_service_account(
            service_account=f"{instance_id}{name}@{settings.GCP_Config.PROJECT_ID}.iam.gserviceaccount.com", gcp_client=gcp_client)
        time.sleep(20)
        VPC.delete_vpc(vpc_name=f"{instance_id}{name}", gcp_client=gcp_client)
        return API_Response(message="Instance deleted successfully", statusCode=200).model_dump()
    except Exception as exp:
        return API_Response(error=str(exp), statusCode=400).model_dump()
