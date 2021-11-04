# syntax=docker/dockerfile:1.3
FROM vlasovanton/gocv:latest as environment

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
RUN apt-get install -y -qq tesseract-ocr-rus
RUN apt-get install -y sqlite3 libsqlite3-dev

FROM environment as build

WORKDIR /app/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN --mount=type=cache,mode=0777,target=~/.cache\
    go build -o core .

FROM environment as exec
COPY --from=build /app/core /opt/dogfound/

RUN mkdir -p /opt/dogfound/data/img && \
    mkdir -p /opt/dogfound/data/new_images
COPY ./data /opt/dogfound/data

CMD ["/opt/dogfound/core", "neural_network:6002","5","4"]