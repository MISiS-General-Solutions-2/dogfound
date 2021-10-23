FROM vlasovanton/gocv:latest as build

WORKDIR /app/

COPY . ./

RUN mkdir app;cd include/gocr-0.52/src; make lib.a; mv libPgm2asc.a /app/libs/

RUN go build -o core .

FROM vlasovanton/gocv:latest
COPY --from=build /app/core /opt/pet-track/

RUN mkdir -p /opt/pet-track/data/img
COPY ./data /opt/pet-track/data

COPY ./include/sqlite3 /usr/bin/sqlite3

CMD ["/opt/pet-track/core"]