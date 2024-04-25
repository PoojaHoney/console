from pydantic import BaseModel
from typing import Any, Optional
class IAM_Role(BaseModel):
    name: str
    description: str
    policies: list
    email: Optional[str] = ""


class Response(BaseModel):
    message: Optional[str] = ""
    data: Optional[Any] = {}
    status_code: Optional[int] = 0
    error: Optional[Any] = ""
