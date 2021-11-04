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
    # change or add fields. Those fields are examples
    crop: List[int]
    probabilities: Optional[Any] = None


class Response(BaseModel):
    is_animal_there: int
    is_it_a_dog: int
    is_the_owner_there: int
    color: int
    tail: int

    additional: Additional


@app.get("/")
async def categorize():
    return {"msg": "hello"}


@app.post("/api/categorize")
async def categorize(req: Request):
    return get_classes(req.image)


def get_classes(file: str) -> Response:
    response = Response(is_animal_there=0, is_it_a_dog=0, is_the_owner_there=0,
                        color=0, tail=0, additional=Additional(crop=[0, 0, 0, 0]))

    response = detect.run_analytics(file, response)
    return response


if __name__ == "__main__":
    uvicorn.run("server:app", reload=True, host="0.0.0.0", port=6002)
