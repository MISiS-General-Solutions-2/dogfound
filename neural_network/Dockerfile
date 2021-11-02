FROM python:3.9.7-buster

WORKDIR /app/

COPY ./requirements.txt ./
RUN python3 -m pip install -r requirements.txt

RUN apt update && \
    yes|apt install libgl1-mesa-glx

COPY . ./

CMD [ "uvicorn", "server:app","--host", "0.0.0.0", "--port", "80" ]