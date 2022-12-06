FROM golang:1.19.3-bullseye as builder
WORKDIR /app

COPY . ./

RUN go mod download

COPY *.go ./

RUN go build -o /helloworld

EXPOSE 8081

ENTRYPOINT ["/helloworld"]