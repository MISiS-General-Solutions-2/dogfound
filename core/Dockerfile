FROM vlasovanton/gocv:latest as build

WORKDIR /app/

COPY . ./

RUN cd gocr-0.52/src; make lib.a;
RUN cd /app/; mkdir libs;mv /app/gocr-0.52/src/libPgm2asc.a ./libs/

RUN go build -o core .

FROM vlasovanton/gocv:latest
COPY --from=build /app/core /opt/pet-track/

RUN mkdir -p /opt/pet-track/data/img
COPY ./data /opt/pet-track/data

# TODO add sqlite3 

CMD ["/opt/pet-track/core"]