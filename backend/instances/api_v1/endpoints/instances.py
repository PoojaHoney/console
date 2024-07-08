import json
from schemas import Instance, Response as API_Response
from config import settings
# , databases
from fastapi import APIRouter
from api_v1.handlers import instances as Instances
from api_v1.handlers.instances import delete_instance_compute_engine

router = APIRouter()


@router.post("/create", response_model=API_Response)
def create_instance(details: Instance):
    try:
        if details.version == "":
            return API_Response(error="product version cannot be empty in instance creation", statusCode=400).model_dump()
        # product = databases.get_mongo_collection(
        #     db_name=settings.PRODUCT_DATABASE, collection_name=settings.PRODUCT_FULL_DETAILS).find_one({"productID": details.productID})
        with open("input.json", "r") as f:
            product = json.load(f)
        if product is None or product.get("productID") != details.productID or product.get("productID") == "":
            return API_Response(error="Product not found or does not exists", statusCode=404).model_dump()
        if details.provider not in product["providers"]:
            return API_Response(error="Provider not supported for the product", statusCode=400).model_dump()
        if details.provider == settings.GCP_Config.CLOUD_PROVIDER:
            for environmentSupport in product["configuration"]["environmentsSupport"]:
                if environmentSupport["environment"] == details.deployedOn and details.deployedOn == "compute_engine":
                    result = Instances.create_instance_compute_engine(
                        details=details, product=product)
                    # if result["statusCode"] != 200:
                    return result
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()


@router.delete("/delete", response_model=API_Response)
def delete_instance(instance_id: str, product: str, name: str):
    try:
        result = delete_instance_compute_engine(
            name=name, instanceId=instance_id, product=product)
        return result
    except Exception as exp:
        return API_Response(error=exp, statusCode=400).model_dump()
