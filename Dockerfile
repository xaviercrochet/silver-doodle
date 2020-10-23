FROM golang:latest

WORKDIR /go/src/app
COPY . .


ENV API_END_POINT="https://storage.googleapis.com/coding-session-rest-api"

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]