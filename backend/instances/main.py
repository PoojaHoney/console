import uvicorn
from fastapi import FastAPI, Request
from api_v1.api import api_router
from config import settings
from fastapi.responses import JSONResponse
import logging

app = FastAPI(
    title=settings.SERVICE_NAME, openapi_url=f"{settings.SERVICE_BASEPATH}/{settings.SERVICE_VERSION}/openapi.json", docs_url="/docs"
)

@app.middleware("http")
async def log_request(request: Request, call_next):
    try:
        response = await call_next(request)
        return response
    except Exception as e:
        logging.exception("Exception during request processing: %s", e)
        return JSONResponse(status_code=500, content={"detail": "Internal Server Error"})

app.include_router(
    api_router, prefix=f"{settings.SERVICE_BASEPATH}/{settings.SERVICE_VERSION}")

if __name__ == "__main__":
    uvicorn.run(app, host=settings.SERVICE_HOST, port=settings.SERVICE_PORT)
