FROM golang:alpine AS builder

RUN apk add --no-cache gcc g++ make git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory_service .

FROM alpine:latest

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /root/

COPY --from=builder /app/inventory_service .

EXPOSE 9080

ENTRYPOINT ["./inventory_service"]
