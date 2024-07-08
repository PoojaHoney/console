from google.oauth2 import service_account as GCPSrvAcc

def get_gcp_crds(crds: dict):
    print("crdsssss: ", crds)
    result =  GCPSrvAcc.Credentials.from_service_account_info(
    # filename="svc-acc-gcp-cloud-automation-cred.json", 
    crds,
    scopes=["https://www.googleapis.com/auth/cloud-platform"])
    print(result)
    return result