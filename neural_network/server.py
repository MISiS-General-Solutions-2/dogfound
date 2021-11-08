from typing import List

import pandas as pd
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel

import detect

app = FastAPI()
test_results_df = pd.read_csv('notebooks/csv/results.csv', index_col=0, sep=',')


class Request(BaseModel):
    image: str


class Additional(BaseModel):
    crop: List[int]


class AddrCamInfo(BaseModel):
    cam_id: str
    address: str


class Response(BaseModel):
    is_animal_there: int
    is_it_a_dog: int
    is_the_owner_there: int
    color: int
    tail: int
    breed: str

    additional: Additional


@app.post("/api/categorize")
async def categorize(req: Request):
    return inference(req.image)


def inference(file: str) -> Response:
    response = Response(is_animal_there=0, is_it_a_dog=0, is_the_owner_there=0,
                        color=0, tail=0, breed="", additional=Additional(crop=[0, 0, 0, 0]))
    response = detect.run_analytics(file, response)
    return response


@app.post("/api/cam-id")
async def cam_id(req: Request):
    return extra_info(req.image)


def extra_info(file: str):
    info = AddrCamInfo(cam_id="", address="")
    file_name = file.split('/')[-1]
    found = test_results_df[test_results_df['filename'] == file_name]
    if len(found) > 0:
        info.cam_id = found.iloc[0]['cam_id']
        info.address = found.iloc[0]['address']
    print(info.cam_id)
    return info.cam_id


if __name__ == "__main__":
    uvicorn.run("server:app", reload=True, host="0.0.0.0", port=6002)
