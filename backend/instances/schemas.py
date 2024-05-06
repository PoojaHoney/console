from pydantic import BaseModel
from typing import Any, Optional, List
class IAM_Role(BaseModel):
    name: str
    description: str
    policies: list
    email: Optional[str] = ""
    product: Optional[str] = ""
    instanceId: Optional[str] = ""


class Response(BaseModel):
    message: Optional[str] = ""
    data: Optional[Any] = {}
    statusCode: Optional[int] = 0
    error: Optional[Any] = ""

class RoutingConfig(BaseModel):
    routingMode: Optional[str] = "REGIONAL"

class SubNetwork(BaseModel):
    vpc_name: str
    name: str
    ipCidrRange: Optional[str] = ""
    region: Optional[str] = ""
    network: str

class FireWall(BaseModel):
    name: str = None
    network: str = None
    direction: Optional[str] = "INGRESS"
    priority: Optional[int]=1000
    sourceRanges: list[str]
    sourceTags: Optional[list[str]] = None
    targetTags: Optional[list[str]] =None
    allowed: list[Any] = None
    denied: Optional[list] = None
    description: Optional[str]=""

class VPC(BaseModel):
    name: str
    product: Optional[str] = ""
    instanceId: Optional[str] = ""
    description: Optional[str] = ""
    cidr: Optional[str] = ""
    network: Optional[str] = ""
    autoCreateSubnetworks: Optional[bool] = False
    enableFlowLogs: Optional[bool] = False
    mtu: Optional[int] = 1460
    routingConfig: Optional[RoutingConfig] = RoutingConfig()
    
from pydantic import BaseModel
from typing import List, Optional

class ComputeEngine(BaseModel):
    class InitializeParams(BaseModel):
        sourceImage: Optional[str] = "projects/debian-cloud/global/images/family/debian-11"
        diskSizeGb: Optional[int] = 10
        diskType: Optional[str] = "pd-standard"

    class Disks(BaseModel):
        boot: Optional[bool] = True
        autoDelete: Optional[bool] = True
        type: Optional[str] = "PERSISTENT"
        initializeParams: 'ComputeEngine.InitializeParams'

    class AccessConfig(BaseModel):
        type: str = "ONE_TO_ONE_NAT"
        name: str = "External NAT"

    class NetworkInterfaces(BaseModel):
        subnetwork: str
        name: str
        accessConfigs: List['ComputeEngine.AccessConfig'] = [{"type": "ONE_TO_ONE_NAT", "name": "External NAT"}]

    class Scopes(BaseModel):
        email: str = "default"
        scopes: List[str] = [
            "https://www.googleapis.com/auth/devstorage.read_write",
            "https://www.googleapis.com/auth/logging.write"
        ]

    class Items(BaseModel):
        key: str
        value: str

    class Metadata(BaseModel):
        items: List['ComputeEngine.Items']

    class Tags(BaseModel):
        items: List[str]

    class ShieldedInstanceConfig(BaseModel):
        enableIntegrityMonitoring: Optional[bool] = True
        enableSecureBoot: Optional[bool] = False
        enableVtpm: Optional[bool] = True

    name: str
    machineType: str
    disks: List[Disks]
    region: str = ""
    networkInterfaces: List[NetworkInterfaces]
    serviceAccounts: List[Scopes]
    metadata: Metadata
    tags: Tags
    zone: Optional[str] = ""
    cpuPlatform: Optional[str] = "Intel Broadwell"
    shieldedInstanceConfig: Optional[ShieldedInstanceConfig] = None

class Instance(BaseModel):
    name: str
    cloudProvider: str
    product: str
    machineType: str
    diskSizeGb: Optional[int] = 10
    diskType: Optional[str] = "pd-standard"
    region: Optional[str]
    zone: Optional[str]
    sourceImage: Optional[str] = None
    
