FROM golang:alpine

WORKDIR /backend
COPY . /backend

ENV DB "bricked-up_prod.db"
EXPOSE 3100

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite
RUN sqlite3 $DB ".read sql/init.sql"

RUN go get .
RUN go build .

CMD ["./backend"]
