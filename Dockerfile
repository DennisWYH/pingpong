FROM golang:1.19.3-bullseye as builder
WORKDIR /app

COPY . ./

RUN go mod download

COPY *.go ./

// Run go build and output binary under helloworld directory
RUN go build -o /helloworld

EXPOSE 8081

// Run the app binary when we run the container
ENTRYPOINT ["/helloworld"]