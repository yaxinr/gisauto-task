FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app
WORKDIR /app

# ENV GO111MODULE=on

COPY . .

WORKDIR /app/server
RUN go get -d -v
# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gisauto
RUN OOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./gisauto

# Run container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/server/. .
# COPY --from=builder /app/server/migrations .

ENTRYPOINT ["./gisauto"]