FROM golang:1.22.2-alpine AS builder

RUN apk update && apk add alpine-sdk && rm -rf /var/cache/apk/*

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

###

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /build/main .
COPY static static
COPY templates templates

EXPOSE 8080

CMD ["/app/main"]
