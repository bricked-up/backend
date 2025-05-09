# Base Image
FROM golang:alpine

# Directory in container
WORKDIR /backend
COPY . /backend

# Setting up environment
ENV DB "/backend/bricked-up_prod.db"
ENV LOGS "/backend/backend.log"
ENV HOST "clabsql.clamv.constructor.university"
ENV PORT ":3100"
EXPOSE 3100

# Setting up database
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

# Building program
RUN go get .
RUN go build .

# Running the server
CMD ["./brickedup"]
