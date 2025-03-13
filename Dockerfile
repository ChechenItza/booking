FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /grpc_booking -trimpath /app/cmd/grpc

FROM alpine

WORKDIR /app

COPY --from=builder /grpc_booking /grpc_booking

CMD ["/grpc_booking"]