FROM python:3.9.7-buster

WORKDIR /app/
COPY . ./
RUN python3 -m pip install -r requirements.txt

CMD [ "uvicorn", "server:app","--host", "0.0.0.0", "--port", "80" ]