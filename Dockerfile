FROM golang:1.19.10-alpine3.17 as build

ENV GOPATH /go

WORKDIR $GOPATH/src/app

COPY . .

RUN go mod download

RUN go mod tidy

RUN go build -o server ./cmd/server

RUN chmod +x server

FROM alpine:3.17

COPY --from=build go/src/app/server /

EXPOSE 8000

ENTRYPOINT [ "./server" ]