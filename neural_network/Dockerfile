# syntax=docker/dockerfile:1.3
FROM python:3.9.7-buster

RUN apt update && \
    yes|apt install libgl1-mesa-glx

WORKDIR /app/

COPY ./requirements.txt ./
RUN --mount=type=cache,target=~/.cache\
    python3 -m pip install -r requirements.txt

COPY . ./

RUN git clone https://github.com/ultralytics/yolov5

CMD [ "uvicorn", "server:app","--host", "0.0.0.0", "--port", "6002" ]