FROM golang:alpine

WORKDIR /backend
COPY . /backend

ENV DB "/backend/bricked-up_prod.db"
EXPOSE 3100

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

RUN go get .
RUN go build .

CMD ["./backend"]
