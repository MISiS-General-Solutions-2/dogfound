# syntax=docker/dockerfile:1.3
FROM vlasovanton/gocv:latest as environment

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
RUN apt-get install -y -qq tesseract-ocr-rus

FROM environment as build

WORKDIR /app/

COPY ./include ./include
RUN mkdir libs && \
    cd include/gocr-0.52/src && \
    make lib.a && \
    mv libPgm2asc.a /app/libs/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN --mount=type=cache,target=~/.cache\
    go build -o core .

FROM environment as exec
COPY --from=build /app/core /opt/dogfound/

RUN mkdir -p /opt/dogfound/data/img && \
    mkdir -p /opt/dogfound/data/new_images
COPY ./data /opt/dogfound/data

COPY ./include/sqlite3 /usr/bin/sqlite3

CMD ["/opt/dogfound/core"]