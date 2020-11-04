FROM golang AS build-env
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT [ "localsearch-api" ]
