FROM golang:1.22 as builder
WORKDIR /go/server

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o server .

FROM ubuntu

COPY --from=builder /go/server/server /local/bin/server

RUN chmod u+x /local/bin/server

ENTRYPOINT /local/bin/server
