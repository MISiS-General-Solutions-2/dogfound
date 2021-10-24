FROM vlasovanton/gocv:latest as build

WORKDIR /app/

COPY . ./

RUN mkdir libs;cd include/gocr-0.52/src; make lib.a; mv libPgm2asc.a /app/libs/

RUN go build -o core .

FROM vlasovanton/gocv:latest
COPY --from=build /app/core /opt/dogfound/

RUN mkdir -p /opt/dogfound/data/img
COPY ./data /opt/dogfound/data

COPY ./include/sqlite3 /usr/bin/sqlite3

CMD ["/opt/dogfound/core"]