from typing import List, Optional, Dict
from pydantic import BaseModel
import datetime
from typing import Any, Optional, List


class IAM_Role(BaseModel):
    name: str = ""
    description: str = ""
    policies: list = []
    email: Optional[str] = ""
    instanceId: Optional[str] = ""
    product: Optional[str] = ""


class Response(BaseModel):
    message: Optional[str] = ""
    data: Optional[Any] = None
    statusCode: Optional[int] = 0
    error: Optional[Any] = None


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
    priority: Optional[int] = 1000
    sourceRanges: list[str]
    sourceTags: Optional[list[str]] = None
    targetTags: Optional[list[str]] = None
    allowed: list[Any] = None
    denied: Optional[list] = None
    description: Optional[str] = ""


class VPC(BaseModel):
    name: str
    description: Optional[str] = ""
    cidr: Optional[str] = ""
    network: Optional[str] = ""
    autoCreateSubnetworks: Optional[bool] = False
    enableFlowLogs: Optional[bool] = False
    mtu: Optional[int] = 1460
    routingConfig: Optional[RoutingConfig] = RoutingConfig()
    ports: Optional[List[str]] = []


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
        accessConfigs: List['ComputeEngine.AccessConfig'] = [
            {"type": "ONE_TO_ONE_NAT", "name": "External NAT"}]

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


class SSHKey(BaseModel):
    userName: Optional[str] = ""
    password: Optional[str] = ""
    publicKey: Optional[str] = ""
    keyFilePath: Optional[str] = ""


class InstanceDatabase(BaseModel):
    type: Optional[str] = ""
    database: Optional[str] = ""
    name: Optional[str] = ""
    portNumber: Optional[int] = 0
    userName: Optional[str] = ""
    password: Optional[str] = ""


class InstanceConfiguration(BaseModel):
    instanceId: str = ""
    productId: str = ""
    ssh: Optional[SSHKey] = SSHKey()
    databases: List[InstanceDatabase] = []


class Instance(BaseModel):
    name: str = ""
    provider: str = ""
    productID: str = ""
    memory: float = 10
    diskType: str = "pd-standard"
    machineType: str = ""
    planId: str = ""
    status: Optional[str] = ""
    createdOn: Optional[datetime.date] = datetime.date.today()
    changedOn: Optional[datetime.date] = datetime.date.today()
    region: str = ""
    prefix: Optional[str] = ""
    description: Optional[str] = ""
    instanceId: Optional[str] = ""
    startTimeInDay: Optional[datetime.time] = datetime.time()
    endTimeInDay: Optional[datetime.time] = datetime.time()
    reinitializeOnSameDay: Optional[bool] = False
    weekOffDays: Optional[bool] = False
    alwaysRun: bool = False
    createdBy: Optional[str] = ""
    runningStatus: Optional[str] = ""
    url: Optional[str] = ""
    version: str = ""
    stage: Optional[str] = ""
    deployedOn: str = ""
    autoScaling: bool = False
    vmSourceImage: str = ""


class VPC_RoutingConfig_Response(BaseModel):
    routingMode: str = ""


class VPC_Response(BaseModel):
    name: str = ""
    selfLink: str = ""
    autoCreateSubnetwork: bool = False
    creationTimestamp: str = ""
    description: str = ""
    id: int = ""
    kind: str = ""
    networkFirewallPolicyEnforcementOrder: str = ""
    routingConfig: VPC_RoutingConfig_Response = {}
    selfLinkWithId: str = ""
    subnetworks: list = []

