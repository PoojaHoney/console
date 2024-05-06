from schemas import Instance, Response as API_Response
from config import settings
from fastapi import APIRouter
from api_v1.handlers import instances as Instances

router = APIRouter()

@router.post("/instances", response_model=API_Response)
def create_instance(details: Instance):
    try:
        if details.cloudProvider == settings.GCP_Config.CLOUD_PROVIDER:
            return Instances.create_instance(details=details)
    except Exception as exp:
        return API_Response(error=exp, status_code=400).model_dump()