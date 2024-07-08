from google.oauth2 import service_account as GCPSrvAcc
import base64


def get_gcp_crds(crds: dict):

    print("crds private key: ", crds["private_key"])
    crds = {
        "type": "service_account",
        "project_id": "console-424205",
        "private_key_id": "84c6cd36101a188c42648801b5cd476794d51613",
        "private_key": crds["private_key"],
        "client_email": "console-automation@console-424205.iam.gserviceaccount.com",
        "client_id": "107868089407702649699",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/console-automation%40console-424205.iam.gserviceaccount.com",
        "universe_domain": "googleapis.com"
    }
    print("crds new: ", crds)
    result = GCPSrvAcc.Credentials.from_service_account_info(
        crds,
        scopes=["https://www.googleapis.com/auth/cloud-platform"]
    )
    print(result)
    return result
