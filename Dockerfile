# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /expensetracking

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /expensetracking /expensetracking

ENV PORT=:2545
ENV DATABASE_URL=localhost

USER nonroot:nonroot

ENTRYPOINT ["/expensetracking"]
