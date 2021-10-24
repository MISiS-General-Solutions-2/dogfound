from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
from pathlib import Path

import detect

app = FastAPI()


class Request(BaseModel):
    dir: str
    imgs: List[str]


class Visualization(BaseModel):
    # change or add fields. Those fields are examples
    crop: List[int]
    probabilities: str


class Response(BaseModel):
    # don't change field names
    is_animal_there: int
    is_it_a_dog: int
    is_the_owner_there: int
    color: int
    tail: int

    vis: Visualization


@app.get("/")
async def categorize():
    return {"msg": "hello"}


@app.post("/api/categorize")
async def categorize(req: Request):
    result = []
    for img in req.imgs:
        result.append(get_classes(req.dir+img))
    return result


def get_classes(file: str) -> Response:
    # implement this function
    response = Response(is_animal_there=0, is_it_a_dog=0, is_the_owner_there=0,
                        color=0, tail=0, vis=Visualization(crop=[0, 0, 5, 5], probabilities="in progress"))

    response = detect.run_analytics(file, response)
    return response

