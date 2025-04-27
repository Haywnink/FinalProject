FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

ARG TODO_PORT=7540
ENV TODO_PORT=${TODO_PORT}

ARG TODO_DBFILE=/data/scheduler.db
ENV TODO_DBFILE=${TODO_DBFILE}

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

EXPOSE ${TODO_PORT}

ENTRYPOINT ["./scheduler"]