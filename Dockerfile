FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jetapi ./cmd/jetapi

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/jetapi .

EXPOSE 8922

ENV PORT=8922
ENV HOST=0.0.0.0

CMD ["./jetapi"]
