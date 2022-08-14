FROM golang:1.18.2-alpine3.14 as builder

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download 

COPY . .
RUN CGO_ENALBE=0 GOOS=linux go build -v -o /usr/local/bin/scout cmd/scout/main.go

FROM alpine:3.14

RUN mkdir /app
COPY --from=builder /usr/local/bin/scout /app/scout

USER 1000

ENTRYPOINT ["/app/scout"]

