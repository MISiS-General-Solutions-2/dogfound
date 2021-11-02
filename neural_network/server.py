from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
from pathlib import Path
import uvicorn

import detect

app = FastAPI()


class Request(BaseModel):
    image: str


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
    return get_classes(req.image)


def get_classes(file: str) -> Response:
    # implement this function
    print(file)
    response = Response(is_animal_there=0, is_it_a_dog=0, is_the_owner_there=0,
                        color=0, tail=0, vis=Visualization(crop=[0, 0, 5, 5], probabilities="in progress"))

    response = detect.run_analytics(file, response)
    print(response)
    return response


if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=6002)
