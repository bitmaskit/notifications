FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o router-service .

FROM alpine:latest AS runner
COPY --from=builder /app/router-service /app/router-service

CMD ["/app/router-service"]