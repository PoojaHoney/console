import uvicorn
from fastapi import FastAPI
from api_v1.api import api_router
from config import settings
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI(
    title=settings.SERVICE_NAME, openapi_url=f"{settings.SERVICE_BASEPATH}/{settings.SERVICE_VERSION}/openapi.json", docs_url="/docs"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  
    allow_credentials=False,
    expose_headers=["Content-Length"], 
    allow_methods=["*"], 
    max_age=86400,
    allow_headers=["*"],
)
app.include_router(
    api_router, prefix=f"{settings.SERVICE_BASEPATH}/{settings.SERVICE_VERSION}")

if __name__ == "__main__":
    uvicorn.run(app, host=settings.SERVICE_DOMAIN, port=settings.SERVICE_PORT)
