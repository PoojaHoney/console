from google.oauth2 import service_account as GCPSrvAcc
import base64


def get_gcp_crds(crds: dict):

    print("crds new: ", crds)
    result = GCPSrvAcc.Credentials.from_service_account_info(
        info=crds,
        scopes=["https://www.googleapis.com/auth/cloud-platform"]
    )
    print(result)
    return result
