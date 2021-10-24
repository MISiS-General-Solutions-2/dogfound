from fastapi import FastAPI
from pydantic import BaseModel
from typing import List

app = FastAPI()

class Request(BaseModel):
    dir: str
    imgs:List[str]

class Timestamp(BaseModel):
    timestamp: int

@app.post("/api/recognize")
async def recognize(req:Request):
    result = []
    for img in req.imgs:
        result.append(GetTextData(req.dir+img))
    return result

# implement ocr
def GetTextData(file: str) ->Timestamp:
    return Timestamp(timestamp=0)