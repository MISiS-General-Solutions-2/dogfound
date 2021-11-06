from fastapi import FastAPI
from pydantic import BaseModel
from typing import Any, List, Optional
from pathlib import Path
import uvicorn

import detect

app = FastAPI()


class Request(BaseModel):
    image: str


class Additional(BaseModel):
    crop: List[int]


class Response(BaseModel):
    is_animal_there: int
    is_it_a_dog: int
    is_the_owner_there: int
    color: int
    tail: int
    breed: str

    additional: Additional


@app.get("/")
async def categorize():
    return {"msg": "hello"}


@app.post("/api/categorize")
async def categorize(req: Request):
    return get_classes(req.image)


def get_classes(file: str) -> Response:
    response = Response(is_animal_there=0, is_it_a_dog=0, is_the_owner_there=0,
                        color=0, tail=0, breed="", additional=Additional(crop=[0, 0, 0, 0]))

    response = detect.run_analytics(file, response)
    return response


@app.post("/api/cam_id")
async def cam_id(req: str):
    return cam_id(req.image)


def cam_id(req: str) -> str:
    return "PVN_hd_TSAO_5300_3"


if __name__ == "__main__":
    uvicorn.run("server:app", reload=True, host="0.0.0.0", port=6002)
